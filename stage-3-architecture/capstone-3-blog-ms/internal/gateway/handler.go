package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	blogv1 "just-go/stage-3-architecture/capstone-3-blog-ms/api/blog/v1"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/resilience"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/tracekit"
)

type Options struct {
	Users          blogv1.UserServiceClient
	Posts          blogv1.PostServiceClient
	Comments       blogv1.CommentServiceClient
	Exporter       tracekit.Exporter
	Limiter        *resilience.Limiter
	CommentBreaker *resilience.Breaker[*blogv1.ListCommentsResponse]
	CommentTimeout time.Duration
	RetryAttempts  int
}

type Handler struct {
	options Options
}

func New(options Options) (http.Handler, error) {
	if options.Users == nil || options.Posts == nil || options.Comments == nil {
		return nil, errors.New("all backend clients are required")
	}
	if options.CommentTimeout <= 0 {
		options.CommentTimeout = 200 * time.Millisecond
	}
	if options.RetryAttempts <= 0 {
		options.RetryAttempts = 2
	}
	if options.CommentBreaker == nil {
		options.CommentBreaker = resilience.NewBreaker[*blogv1.ListCommentsResponse]("comments", 3, time.Second)
	}
	return &Handler{options: options}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := tracekit.WithRemote(r.Context(), r.Header.Get("traceparent"))
	ctx, end := tracekit.Start(ctx, h.options.Exporter, "gateway", r.Method+" "+r.URL.Path)
	var requestErr error
	defer func() { end(requestErr) }()
	r = r.WithContext(ctx)
	w.Header().Set("X-Trace-ID", traceID(ctx))
	switch {
	case r.URL.Path == "/healthz":
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	case r.URL.Path == "/api/register" && r.Method == http.MethodPost:
		requestErr = h.register(w, r)
	case r.URL.Path == "/api/login" && r.Method == http.MethodPost:
		requestErr = h.login(w, r)
	case r.URL.Path == "/api/posts" && r.Method == http.MethodPost:
		requestErr = h.createPost(w, r)
	case r.URL.Path == "/api/posts" && r.Method == http.MethodGet:
		requestErr = h.listPosts(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/posts/"):
		requestErr = h.postSubroutes(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) error {
	var req blogv1.RegisterRequest
	if err := decode(r, &req); err != nil {
		return writeError(w, err)
	}
	out, err := h.options.Users.Register(r.Context(), &req)
	if err != nil {
		return writeError(w, err)
	}
	writeJSON(w, http.StatusCreated, out)
	return nil
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) error {
	var req blogv1.LoginRequest
	if err := decode(r, &req); err != nil {
		return writeError(w, err)
	}
	out, err := h.options.Users.Login(r.Context(), &req)
	if err != nil {
		return writeError(w, err)
	}
	writeJSON(w, http.StatusOK, out)
	return nil
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) error {
	user, err := h.authenticate(r)
	if err != nil {
		return writeError(w, err)
	}
	var req blogv1.CreatePostRequest
	if err := decode(r, &req); err != nil {
		return writeError(w, err)
	}
	req.AuthorId = user.Id
	post, err := h.options.Posts.CreatePost(r.Context(), &req)
	if err != nil {
		return writeError(w, err)
	}
	writeJSON(w, http.StatusCreated, post)
	return nil
}

func (h *Handler) listPosts(w http.ResponseWriter, r *http.Request) error {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	out, err := h.options.Posts.ListPosts(r.Context(), &blogv1.ListPostsRequest{Tag: r.URL.Query().Get("tag"), Page: int32(page), PageSize: int32(pageSize)})
	if err != nil {
		return writeError(w, err)
	}
	writeJSON(w, http.StatusOK, out)
	return nil
}

func (h *Handler) postSubroutes(w http.ResponseWriter, r *http.Request) error {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/posts/"), "/")
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || id <= 0 {
		return writeError(w, status.Error(codes.InvalidArgument, "invalid post id"))
	}
	if len(parts) == 1 && r.Method == http.MethodGet {
		return h.getPost(w, r, id)
	}
	if len(parts) == 2 && parts[1] == "comments" && r.Method == http.MethodPost {
		return h.createComment(w, r, id)
	}
	http.NotFound(w, r)
	return nil
}

func (h *Handler) getPost(w http.ResponseWriter, r *http.Request, id int64) error {
	if !h.options.Limiter.Allow() {
		writeJSON(w, http.StatusTooManyRequests, map[string]string{"error": resilience.ErrLimited.Error()})
		return resilience.ErrLimited
	}
	post, err := h.options.Posts.GetPost(r.Context(), &blogv1.GetPostRequest{Id: id})
	if err != nil {
		return writeError(w, err)
	}
	author, err := h.options.Users.GetUser(r.Context(), &blogv1.GetUserRequest{Id: post.AuthorId})
	if err != nil {
		return writeError(w, err)
	}
	commentCtx, cancel := context.WithTimeout(r.Context(), h.options.CommentTimeout)
	defer cancel()
	comments, commentErr := h.options.CommentBreaker.Execute(func() (*blogv1.ListCommentsResponse, error) {
		return resilience.Retry(commentCtx, h.options.RetryAttempts, 5*time.Millisecond, retryable, func(ctx context.Context) (*blogv1.ListCommentsResponse, error) {
			return h.options.Comments.ListComments(ctx, &blogv1.ListCommentsRequest{PostId: id})
		})
	})
	degraded := commentErr != nil
	if comments == nil {
		comments = &blogv1.ListCommentsResponse{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"post": post, "author": author, "comments": comments.Comments, "comments_degraded": degraded})
	return nil
}

func (h *Handler) createComment(w http.ResponseWriter, r *http.Request, postID int64) error {
	user, err := h.authenticate(r)
	if err != nil {
		return writeError(w, err)
	}
	if _, err := h.options.Posts.GetPost(r.Context(), &blogv1.GetPostRequest{Id: postID}); err != nil {
		return writeError(w, err)
	}
	var req blogv1.CreateCommentRequest
	if err := decode(r, &req); err != nil {
		return writeError(w, err)
	}
	req.PostId, req.AuthorId = postID, user.Id
	comment, err := h.options.Comments.CreateComment(r.Context(), &req)
	if err != nil {
		return writeError(w, err)
	}
	writeJSON(w, http.StatusCreated, comment)
	return nil
}

func (h *Handler) authenticate(r *http.Request) (*blogv1.User, error) {
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return nil, status.Error(codes.Unauthenticated, "missing bearer token")
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "missing bearer token")
	}
	return h.options.Users.ValidateToken(r.Context(), &blogv1.ValidateTokenRequest{Token: token})
}

func retryable(err error) bool {
	code := status.Code(err)
	return code == codes.Unavailable || code == codes.ResourceExhausted || code == codes.DeadlineExceeded
}

func decode(r *http.Request, target any) error {
	defer func() { _ = r.Body.Close() }()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return status.Error(codes.InvalidArgument, "request body must contain one JSON object")
	}
	return nil
}

func writeError(w http.ResponseWriter, err error) error {
	code := http.StatusInternalServerError
	switch status.Code(err) {
	case codes.InvalidArgument:
		code = http.StatusUnprocessableEntity
	case codes.Unauthenticated:
		code = http.StatusUnauthorized
	case codes.PermissionDenied:
		code = http.StatusForbidden
	case codes.NotFound:
		code = http.StatusNotFound
	case codes.AlreadyExists:
		code = http.StatusConflict
	case codes.ResourceExhausted:
		code = http.StatusTooManyRequests
	case codes.Unavailable:
		code = http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		code = http.StatusGatewayTimeout
	}
	writeJSON(w, code, map[string]string{"error": status.Convert(err).Message()})
	return err
}

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

func traceID(ctx context.Context) string {
	value, _ := tracekit.FromContext(ctx)
	return value.TraceID
}
