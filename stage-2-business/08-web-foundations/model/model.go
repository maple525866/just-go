package model

// Article is the JSON representation returned by the demo API.
type Article struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

// CreateArticleRequest is the JSON request body accepted by POST /api/articles.
type CreateArticleRequest struct {
	Title string   `json:"title" validate:"required,notblank,max=80"`
	Body  string   `json:"body" validate:"required,notblank,min=10"`
	Tags  []string `json:"tags" validate:"required,dive,required,notblank,max=20"`
}

// ArticleListResponse wraps list endpoints so future pagination fields can be
// added without changing the top-level JSON shape.
type ArticleListResponse struct {
	Items []Article `json:"items"`
	Total int       `json:"total"`
}

// FieldError describes one validation problem in a client-readable shape.
type FieldError struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

// ErrorResponse is the shared JSON error envelope for REST examples.
type ErrorResponse struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Fields  []FieldError `json:"fields,omitempty"`
}
