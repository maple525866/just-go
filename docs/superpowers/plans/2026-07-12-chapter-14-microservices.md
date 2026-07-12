# Chapter 14 Microservices Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (- [ ]) syntax for tracking.

**Goal:** Build a runnable, self-contained Chapter 14 that teaches real protobuf/gRPC communication, service discovery, dynamic configuration, and an HTTP API Gateway through a product-and-inventory example.

**Architecture:** Product and inventory services communicate over real gRPC transports using committed generated contracts. Narrow discovery and configuration interfaces have deterministic in-memory adapters; a standard-library HTTP Gateway applies edge policy, resolves gRPC endpoints, and aggregates both services under one deadline.

**Tech Stack:** Go 1.24, grpc-go v1.75.1, protobuf-go v1.36.11, protoc-gen-go-grpc v1.5.1, Buf CLI v1.65.0, net/http, Go standard testing, bufconn.

## Global Constraints

- Keep the repository module declaration at Go 1.24; do not select dependencies that require Go 1.25.
- Default execution and all Chapter 14 tests require no Docker, Consul, etcd, Nacos, Kubernetes, TLS endpoint, or message broker.
- Hand-written .proto files are the contract source; generated .pb.go files are committed, are never hand-edited, and must be reproducible with pinned commands.
- Service discovery stores network addresses and metadata, never Go service objects.
- Unknown internal errors are hidden; Gateway required-downstream failure returns no partial aggregate.
- Every shared-state component owns its lock and lifecycle; streaming and subscription loops observe context cancellation.
- Use test-first RED/GREEN cycles, gofmt all Go files, and preserve unrelated untracked .agents/ and docs/superpowers/specs/github-issue-add-github-actions-ci.md.
- Learner-facing prose is simplified Chinese; identifiers, protocols, and commands retain standard English names.

---

## Planned File Structure

- stage-3-architecture/14-microservices/buf.yaml: protobuf module definition.
- stage-3-architecture/14-microservices/buf.gen.yaml: pinned local Go plugin commands.
- stage-3-architecture/14-microservices/api/product/v1/product.proto: ProductService contract.
- stage-3-architecture/14-microservices/api/inventory/v1/inventory.proto: InventoryService contract and all three RPC forms.
- stage-3-architecture/14-microservices/api/**/v1/*.pb.go: generated protobuf and gRPC code.
- stage-3-architecture/14-microservices/internal/product/catalog.go: concurrency-safe immutable product lookup.
- stage-3-architecture/14-microservices/internal/product/service.go: ProductService gRPC adapter.
- stage-3-architecture/14-microservices/internal/inventory/store.go: concurrency-safe stock state and subscriptions.
- stage-3-architecture/14-microservices/internal/inventory/service.go: InventoryService unary and streaming adapter.
- stage-3-architecture/14-microservices/internal/discovery/discovery.go: discovery value types and interfaces.
- stage-3-architecture/14-microservices/internal/discovery/memory.go: deterministic in-memory registry.
- stage-3-architecture/14-microservices/internal/configcenter/config.go: Gateway configuration, validation, and rollout.
- stage-3-architecture/14-microservices/internal/configcenter/memory.go: versioned in-memory configuration store.
- stage-3-architecture/14-microservices/internal/gateway/auth.go: Bearer authentication.
- stage-3-architecture/14-microservices/internal/gateway/limiter.go: single-process fixed-window limiter.
- stage-3-architecture/14-microservices/internal/gateway/connections.go: discovery-backed cached gRPC connections.
- stage-3-architecture/14-microservices/internal/gateway/handler.go: route policy and concurrent aggregation.
- stage-3-architecture/14-microservices/internal/gateway/errors.go: stable gRPC-to-HTTP error mapping.
- stage-3-architecture/14-microservices/app.go: composition lifecycle and demonstration request.
- stage-3-architecture/14-microservices/main.go: executable entry point.
- Matching _test.go files: focused unit, transport, HTTP, cancellation, and lifecycle tests.
- stage-3-architecture/14-microservices/README.md and EXERCISES.md: learning materials.
- ROADMAP.md: Chapter 14 output and completion state.
- openspec/changes/chapter-14-microservices/: approved proposal, design, specifications, and tracked tasks.

### Task 1: Protocol Contracts and Reproducible Generation

**Files:**
- Create: stage-3-architecture/14-microservices/buf.yaml
- Create: stage-3-architecture/14-microservices/buf.gen.yaml
- Create: stage-3-architecture/14-microservices/api/product/v1/product.proto
- Create: stage-3-architecture/14-microservices/api/inventory/v1/inventory.proto
- Generate: stage-3-architecture/14-microservices/api/product/v1/product.pb.go
- Generate: stage-3-architecture/14-microservices/api/product/v1/product_grpc.pb.go
- Generate: stage-3-architecture/14-microservices/api/inventory/v1/inventory.pb.go
- Generate: stage-3-architecture/14-microservices/api/inventory/v1/inventory_grpc.pb.go
- Modify: go.mod
- Modify: go.sum
- Test: stage-3-architecture/14-microservices/api/contracts_test.go

**Interfaces:**
- Produces: productv1.ProductService with GetProduct(context.Context, *GetProductRequest) and inventoryv1.InventoryService with GetStock, WatchStock, and SyncStock generated interfaces.
- Produces: GetProductRequest{Sku}, GetProductResponse{Sku, Name, PriceCents}, GetStockRequest{Sku}, GetStockResponse{Sku, Quantity, Version}, WatchStockRequest{Sku}, WatchStockResponse{Sku, Quantity, Version}, SyncStockRequest{Sku, Delta}, and SyncStockResponse{Sku, Quantity, Version}.

- [ ] **Step 1: Write the failing contract reflection test**

Create api/contracts_test.go. The test asserts exact fully-qualified services, method names, and streaming flags so a renamed or incorrectly shaped protocol fails visibly.

~~~go
package api_test

import (
    "testing"

    inventoryv1 "just-go/stage-3-architecture/14-microservices/api/inventory/v1"
    productv1 "just-go/stage-3-architecture/14-microservices/api/product/v1"
)

func TestGeneratedServiceDescriptors(t *testing.T) {
    if got := productv1.ProductService_ServiceDesc.ServiceName; got != "product.v1.ProductService" {
        t.Fatalf("product service name = %q", got)
    }
    if got := inventoryv1.InventoryService_ServiceDesc.ServiceName; got != "inventory.v1.InventoryService" {
        t.Fatalf("inventory service name = %q", got)
    }
    methods := inventoryv1.InventoryService_ServiceDesc.Methods
    if len(methods) != 1 || methods[0].MethodName != "GetStock" {
        t.Fatalf("unary methods = %#v", methods)
    }
    streams := inventoryv1.InventoryService_ServiceDesc.Streams
    if len(streams) != 2 {
        t.Fatalf("stream count = %d", len(streams))
    }
    if streams[0].StreamName != "WatchStock" || !streams[0].ServerStreams || streams[0].ClientStreams {
        t.Fatalf("WatchStock descriptor = %#v", streams[0])
    }
    if streams[1].StreamName != "SyncStock" || !streams[1].ServerStreams || !streams[1].ClientStreams {
        t.Fatalf("SyncStock descriptor = %#v", streams[1])
    }
}
~~~

- [ ] **Step 2: Run the contract test and verify RED**

Run:

~~~powershell
go test ./stage-3-architecture/14-microservices/api/... -run TestGeneratedServiceDescriptors -count=1
~~~

Expected: FAIL because productv1 and inventoryv1 generated packages do not exist.

- [ ] **Step 3: Add protocol sources and pinned generator configuration**

Create buf.yaml:

~~~yaml
version: v2
modules:
  - path: api
~~~

Create buf.gen.yaml:

~~~yaml
version: v2
plugins:
  - local: ["go", "run", "google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11"]
    out: api
    opt: paths=source_relative
  - local: ["go", "run", "google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1"]
    out: api
    opt: paths=source_relative
~~~

Create product.proto:

~~~proto
syntax = "proto3";

package product.v1;

option go_package = "just-go/stage-3-architecture/14-microservices/api/product/v1;productv1";

service ProductService {
  rpc GetProduct(GetProductRequest) returns (GetProductResponse);
}

message GetProductRequest {
  string sku = 1;
}

message GetProductResponse {
  string sku = 1;
  string name = 2;
  int64 price_cents = 3;
}
~~~

Create inventory.proto:

~~~proto
syntax = "proto3";

package inventory.v1;

option go_package = "just-go/stage-3-architecture/14-microservices/api/inventory/v1;inventoryv1";

service InventoryService {
  rpc GetStock(GetStockRequest) returns (GetStockResponse);
  rpc WatchStock(WatchStockRequest) returns (stream WatchStockResponse);
  rpc SyncStock(stream SyncStockRequest) returns (stream SyncStockResponse);
}

message GetStockRequest {
  string sku = 1;
}

message WatchStockRequest {
  string sku = 1;
}

message SyncStockRequest {
  string sku = 1;
  int64 delta = 2;
}

message GetStockResponse {
  string sku = 1;
  int64 quantity = 2;
  uint64 version = 3;
}

message WatchStockResponse {
  string sku = 1;
  int64 quantity = 2;
  uint64 version = 3;
}

message SyncStockResponse {
  string sku = 1;
  int64 quantity = 2;
  uint64 version = 3;
}
~~~

- [ ] **Step 4: Generate code and add compatible runtime dependencies**

Run:

~~~powershell
Push-Location stage-3-architecture/14-microservices
go run github.com/bufbuild/buf/cmd/buf@v1.65.0 generate
Pop-Location
go get google.golang.org/grpc@v1.75.1 google.golang.org/protobuf@v1.36.11
go mod tidy
~~~

Expected: four generated .pb.go files exist; go.mod remains go 1.24 and contains grpc v1.75.1 plus protobuf v1.36.11.

- [ ] **Step 5: Run the contract test and protocol lint**

Run:

~~~powershell
go test ./stage-3-architecture/14-microservices/api/... -count=1
Push-Location stage-3-architecture/14-microservices
go run github.com/bufbuild/buf/cmd/buf@v1.65.0 lint
Pop-Location
~~~

Expected: contract tests PASS and Buf reports no lint errors. Generated-code drift is checked after the generated files have a committed baseline in Task 8.

- [ ] **Step 6: Commit protocol work**

~~~powershell
git add go.mod go.sum stage-3-architecture/14-microservices/buf.yaml stage-3-architecture/14-microservices/buf.gen.yaml stage-3-architecture/14-microservices/api
git commit -m "feat: add chapter 14 grpc contracts"
~~~

### Task 2: Product Catalog and Unary gRPC Services

**Files:**
- Create: stage-3-architecture/14-microservices/internal/product/catalog.go
- Create: stage-3-architecture/14-microservices/internal/product/catalog_test.go
- Create: stage-3-architecture/14-microservices/internal/product/service.go
- Create: stage-3-architecture/14-microservices/internal/product/service_test.go
- Create: stage-3-architecture/14-microservices/internal/inventory/store.go
- Create: stage-3-architecture/14-microservices/internal/inventory/store_test.go
- Create: stage-3-architecture/14-microservices/internal/inventory/service.go
- Create: stage-3-architecture/14-microservices/internal/inventory/service_test.go

**Interfaces:**
- Produces: product.Product{SKU string, Name string, PriceCents int64}, product.NewCatalog([]Product) (*Catalog, error), (*Catalog).Get(string) (Product, error), product.NewService(*Catalog) *Service.
- Produces: inventory.Stock{SKU string, Quantity int64, Version uint64}, inventory.NewStore(map[string]int64) (*Store, error), (*Store).Get(string) (Stock, error), (*Store).Adjust(string, int64) (Stock, error), inventory.NewService(*Store) *Service.
- Stable package errors: ErrInvalidProduct, ErrProductNotFound, ErrInvalidStock, ErrStockNotFound.

- [ ] **Step 1: Write failing table tests for domain stores**

Tests must cover blank SKU/name, non-positive price, duplicate SKU, immutable product value reads, unknown product, blank inventory SKU, negative initial quantity, unknown stock, adjustment below zero, version increment, and concurrent reads/adjustments. The first RED test is:

~~~go
func TestCatalogGet(t *testing.T) {
    catalog, err := NewCatalog([]Product{{SKU: "book-1", Name: "Go Book", PriceCents: 9900}})
    if err != nil {
        t.Fatal(err)
    }
    got, err := catalog.Get("book-1")
    if err != nil {
        t.Fatal(err)
    }
    if got != (Product{SKU: "book-1", Name: "Go Book", PriceCents: 9900}) {
        t.Fatalf("product = %#v", got)
    }
}
~~~

Inventory's first RED test is:

~~~go
func TestStoreAdjustIncrementsVersion(t *testing.T) {
    store, err := NewStore(map[string]int64{"book-1": 10})
    if err != nil {
        t.Fatal(err)
    }
    got, err := store.Adjust("book-1", -2)
    if err != nil {
        t.Fatal(err)
    }
    if got.Quantity != 8 || got.Version != 2 {
        t.Fatalf("stock = %#v", got)
    }
}
~~~

- [ ] **Step 2: Run domain tests and verify RED**

Run:

~~~powershell
go test ./stage-3-architecture/14-microservices/internal/product ./stage-3-architecture/14-microservices/internal/inventory -count=1
~~~

Expected: FAIL because Product, Catalog, Stock, and Store are undefined.

- [ ] **Step 3: Implement minimal concurrency-safe stores**

Catalog owns map[string]Product behind sync.RWMutex; construction trims and validates all strings before copying inputs. Store owns map[string]Stock behind sync.RWMutex; initial version is 1, Adjust rejects blank SKU, zero delta, unknown SKU, and negative resulting quantity, then increments version exactly once. Return value structs, never internal pointers.

Required declarations:

~~~go
var (
    ErrInvalidProduct  = errors.New("invalid product")
    ErrProductNotFound = errors.New("product not found")
)

type Product struct {
    SKU        string
    Name       string
    PriceCents int64
}

type Catalog struct {
    mu       sync.RWMutex
    products map[string]Product
}

var (
    ErrInvalidStock  = errors.New("invalid stock")
    ErrStockNotFound = errors.New("stock not found")
)

type Stock struct {
    SKU      string
    Quantity int64
    Version  uint64
}

type Store struct {
    mu    sync.RWMutex
    stock map[string]Stock
}
~~~

- [ ] **Step 4: Run domain tests and verify GREEN**

Run:

~~~powershell
go test ./stage-3-architecture/14-microservices/internal/product ./stage-3-architecture/14-microservices/internal/inventory -count=1
~~~

Expected: PASS.

- [ ] **Step 5: Write failing unary service tests**

Use direct service calls for status mapping and cover valid response, blank SKU => codes.InvalidArgument, unknown SKU => codes.NotFound, and canceled context => codes.Canceled. Example:

~~~go
func TestServiceGetProductMapsMissingProduct(t *testing.T) {
    catalog, err := NewCatalog([]Product{{SKU: "book-1", Name: "Go Book", PriceCents: 9900}})
    if err != nil {
        t.Fatal(err)
    }
    _, err = NewService(catalog).GetProduct(context.Background(), &productv1.GetProductRequest{Sku: "missing"})
    if status.Code(err) != codes.NotFound {
        t.Fatalf("code = %v, err = %v", status.Code(err), err)
    }
}
~~~

- [ ] **Step 6: Run unary service tests and verify RED**

Run:

~~~powershell
go test ./stage-3-architecture/14-microservices/internal/product ./stage-3-architecture/14-microservices/internal/inventory -run 'TestService(GetProduct|GetStock)' -count=1
~~~

Expected: FAIL because Service constructors and generated server implementations are missing.

- [ ] **Step 7: Implement unary service adapters**

Both services embed the generated Unimplemented server, check ctx.Err before work, validate request and SKU, delegate to the store, and map only known package errors.

~~~go
type Service struct {
    productv1.UnimplementedProductServiceServer
    catalog *Catalog
}

func NewService(catalog *Catalog) *Service {
    return &Service{catalog: catalog}
}

func (s *Service) GetProduct(ctx context.Context, req *productv1.GetProductRequest) (*productv1.GetProductResponse, error)
~~~

Inventory uses the parallel signatures:

~~~go
type Service struct {
    inventoryv1.UnimplementedInventoryServiceServer
    store *Store
}

func NewService(store *Store) *Service
func (s *Service) GetStock(context.Context, *inventoryv1.GetStockRequest) (*inventoryv1.GetStockResponse, error)
~~~

- [ ] **Step 8: Run focused and race tests**

Run:

~~~powershell
gofmt -w stage-3-architecture/14-microservices/internal/product stage-3-architecture/14-microservices/internal/inventory
go test ./stage-3-architecture/14-microservices/internal/product ./stage-3-architecture/14-microservices/internal/inventory -count=1
go test -race -count=1 ./stage-3-architecture/14-microservices/internal/product ./stage-3-architecture/14-microservices/internal/inventory
~~~

Expected: PASS with no race reports.

- [ ] **Step 9: Commit unary services**

~~~powershell
git add stage-3-architecture/14-microservices/internal/product stage-3-architecture/14-microservices/internal/inventory
git commit -m "feat: add chapter 14 unary grpc services"
~~~

### Task 3: Inventory Streaming and Transport Integration

**Files:**
- Modify: stage-3-architecture/14-microservices/internal/inventory/store.go
- Modify: stage-3-architecture/14-microservices/internal/inventory/store_test.go
- Modify: stage-3-architecture/14-microservices/internal/inventory/service.go
- Create: stage-3-architecture/14-microservices/internal/inventory/stream_test.go
- Create: stage-3-architecture/14-microservices/internal/inventory/transport_test.go

**Interfaces:**
- Produces: (*Store).Watch(context.Context, string) (<-chan Stock, error).
- Implements generated InventoryService_WatchStockServer and InventoryService_SyncStockServer handlers.
- Watch sends the current snapshot first, then coalesces unread changes so a slow subscriber receives the latest accepted state; cancellation closes the subscriber channel.

- [ ] **Step 1: Write failing Store.Watch tests**

Cover initial snapshot, ordered changes, unknown SKU, canceled context, slow subscriber not blocking Adjust, and watcher cleanup. Use a buffered channel expectation:

~~~go
func TestStoreWatchReceivesCurrentAndUpdatedStock(t *testing.T) {
    store, err := NewStore(map[string]int64{"book-1": 10})
    if err != nil {
        t.Fatal(err)
    }
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    updates, err := store.Watch(ctx, "book-1")
    if err != nil {
        t.Fatal(err)
    }
    if got := <-updates; got.Quantity != 10 || got.Version != 1 {
        t.Fatalf("initial = %#v", got)
    }
    if _, err := store.Adjust("book-1", -2); err != nil {
        t.Fatal(err)
    }
    if got := <-updates; got.Quantity != 8 || got.Version != 2 {
        t.Fatalf("updated = %#v", got)
    }
}
~~~

- [ ] **Step 2: Run watch tests and verify RED**

Run:

~~~powershell
go test ./stage-3-architecture/14-microservices/internal/inventory -run TestStoreWatch -count=1
~~~

Expected: FAIL because Store.Watch is undefined.

- [ ] **Step 3: Implement Store subscriptions**

Add subscriber IDs and per-SKU buffered channels under the existing mutex. Watch validates before registration, enqueues the current snapshot, and starts one goroutine that removes and closes its subscription on ctx.Done. Adjust performs non-blocking latest-snapshot delivery while holding the same mutex used by cancellation and Close; replace a stale buffered value before sending the newest snapshot. Keeping send and close under one lock prevents send-on-closed-channel races without letting a slow subscriber block a producer.

Required fields:

~~~go
type subscription struct {
    id uint64
    ch chan Stock
}

type Store struct {
    mu          sync.RWMutex
    stock       map[string]Stock
    watchers    map[string]map[uint64]chan Stock
    nextWatcher uint64
}
~~~

- [ ] **Step 4: Run watch and race tests**

Run:

~~~powershell
go test ./stage-3-architecture/14-microservices/internal/inventory -run TestStoreWatch -count=1
go test -race -count=1 ./stage-3-architecture/14-microservices/internal/inventory
~~~

Expected: PASS with no race reports.

- [ ] **Step 5: Write failing real-transport streaming tests**

Start grpc.NewServer on bufconn.Listen, register Service, and dial with grpc.NewClient plus grpc.WithContextDialer and insecure.NewCredentials. Cover:

- WatchStock receives current and adjusted snapshots, then Recv returns codes.Canceled after client cancellation.
- SyncStock Send/Recv returns one result per delta in order.
- SyncStock invalid request returns codes.InvalidArgument.
- SyncStock CloseSend yields io.EOF after all responses.

The test helper has this exact signature:

~~~go
func newInventoryClient(t *testing.T, service inventoryv1.InventoryServiceServer) inventoryv1.InventoryServiceClient
~~~

- [ ] **Step 6: Run streaming transport tests and verify RED**

Run:

~~~powershell
go test ./stage-3-architecture/14-microservices/internal/inventory -run 'Test(WatchStock|SyncStock)' -count=1
~~~

Expected: FAIL because WatchStock and SyncStock retain generated unimplemented behavior.

- [ ] **Step 7: Implement server-streaming and bidirectional-streaming**

WatchStock calls store.Watch(stream.Context, req.Sku), maps known errors, then loops over updates and stream.Send until the channel closes or context is canceled.

SyncStock loops on stream.Recv, returns nil on io.EOF, maps canceled contexts, validates each request through Store.Adjust, and calls stream.Send once per accepted adjustment. A failed request terminates the stream with the stable mapped status.

Required signatures:

~~~go
func (s *Service) WatchStock(req *inventoryv1.WatchStockRequest, stream inventoryv1.InventoryService_WatchStockServer) error
func (s *Service) SyncStock(stream inventoryv1.InventoryService_SyncStockServer) error
~~~

- [ ] **Step 8: Run all inventory and transport tests**

Run:

~~~powershell
gofmt -w stage-3-architecture/14-microservices/internal/inventory
go test ./stage-3-architecture/14-microservices/internal/inventory -count=1
go test -race -count=1 ./stage-3-architecture/14-microservices/internal/inventory
~~~

Expected: PASS; cancellation tests finish within their explicit one-second test deadline.

- [ ] **Step 9: Commit streaming work**

~~~powershell
git add stage-3-architecture/14-microservices/internal/inventory
git commit -m "feat: demonstrate grpc inventory streams"
~~~

### Task 4: Service Discovery

**Files:**
- Create: stage-3-architecture/14-microservices/internal/discovery/discovery.go
- Create: stage-3-architecture/14-microservices/internal/discovery/memory.go
- Create: stage-3-architecture/14-microservices/internal/discovery/memory_test.go

**Interfaces:**
- Produces Instance{Service, ID, Address string}.
- Produces Registry.Register(Instance) (func() error, error), Resolve(string) (Instance, error), Watch(context.Context, string) (<-chan []Instance, error), and Close() error.
- Stable errors: ErrInvalidInstance, ErrDuplicateInstance, ErrUnavailable, ErrClosed.

- [ ] **Step 1: Write failing discovery contract tests**

Table tests cover invalid fields, duplicate service-plus-ID, idempotent deregistration, deterministic lexical Resolve order, immutable sorted watcher snapshots, immediate initial snapshot, cancellation, slow watcher, close, and concurrent registration/resolution.

~~~go
func TestMemoryRegistryResolveIsDeterministic(t *testing.T) {
    registry := NewMemoryRegistry()
    defer registry.Close()
    for _, instance := range []Instance{
        {Service: "product", ID: "b", Address: "127.0.0.1:2"},
        {Service: "product", ID: "a", Address: "127.0.0.1:1"},
    } {
        if _, err := registry.Register(instance); err != nil {
            t.Fatal(err)
        }
    }
    got, err := registry.Resolve("product")
    if err != nil {
        t.Fatal(err)
    }
    if got.ID != "a" {
        t.Fatalf("resolved ID = %q", got.ID)
    }
}
~~~

- [ ] **Step 2: Run discovery tests and verify RED**

Run:

~~~powershell
go test ./stage-3-architecture/14-microservices/internal/discovery -count=1
~~~

Expected: FAIL because the discovery package does not exist.

- [ ] **Step 3: Define the exact discovery contract**

~~~go
type Instance struct {
    Service string
    ID      string
    Address string
}

type Registry interface {
    Register(Instance) (deregister func() error, err error)
    Resolve(service string) (Instance, error)
    Watch(ctx context.Context, service string) (<-chan []Instance, error)
    Close() error
}
~~~

Use errors.Is-compatible sentinel errors. Validate trimmed non-empty values and net.SplitHostPort-compatible addresses.

- [ ] **Step 4: Implement the in-memory registry**

MemoryRegistry owns service maps and watchers under sync.RWMutex. Snapshot creates a fresh sorted []Instance. Register atomically checks duplicates, stores the instance, and publishes a non-blocking latest snapshot while holding the mutex. Deregistration uses sync.Once. Watch installs a one-element buffered channel, sends the current snapshot immediately, and removes then closes the watcher under the mutex on cancellation. Close marks closed and closes every watcher under the same mutex. Shared locking makes publishing and channel closure race-free; one-element buffers keep publishing non-blocking.

- [ ] **Step 5: Run focused and race tests**

Run:

~~~powershell
gofmt -w stage-3-architecture/14-microservices/internal/discovery
go test ./stage-3-architecture/14-microservices/internal/discovery -count=1
go test -race -count=1 ./stage-3-architecture/14-microservices/internal/discovery
~~~

Expected: PASS with no race reports.

- [ ] **Step 6: Commit discovery**

~~~powershell
git add stage-3-architecture/14-microservices/internal/discovery
git commit -m "feat: add in-memory service discovery"
~~~

### Task 5: Dynamic Configuration and Deterministic Rollout

**Files:**
- Create: stage-3-architecture/14-microservices/internal/configcenter/config.go
- Create: stage-3-architecture/14-microservices/internal/configcenter/config_test.go
- Create: stage-3-architecture/14-microservices/internal/configcenter/memory.go
- Create: stage-3-architecture/14-microservices/internal/configcenter/memory_test.go

**Interfaces:**
- Produces GatewayConfig{RouteEnabled bool, RequestTimeout time.Duration, RateLimit int, RateWindow time.Duration, RolloutPercent uint32, BearerToken string}.
- Produces Snapshot{Version uint64, Config GatewayConfig}.
- Produces Store.Current() (Snapshot, error), Update(GatewayConfig) (Snapshot, error), Watch(context.Context) (<-chan Snapshot, error), Close() error.
- Produces InRollout(key string, percent uint32) bool.
- Stable errors: ErrInvalidConfig, ErrClosed.

- [ ] **Step 1: Write failing validation and rollout tests**

Use a table for zero/negative timeout, non-positive rate limit/window, rollout above 100, and blank token. Assert zero percent always false, 100 always true, and 1–99 produces the same result for a repeated key.

~~~go
func TestInRolloutIsStable(t *testing.T) {
    first := InRollout("learner-42", 25)
    for i := 0; i < 100; i++ {
        if got := InRollout("learner-42", 25); got != first {
            t.Fatalf("decision changed: first=%v got=%v", first, got)
        }
    }
    if InRollout("any", 0) {
        t.Fatal("zero-percent rollout enabled")
    }
    if !InRollout("any", 100) {
        t.Fatal("full rollout disabled")
    }
}
~~~

- [ ] **Step 2: Run config tests and verify RED**

Run:

~~~powershell
go test ./stage-3-architecture/14-microservices/internal/configcenter -count=1
~~~

Expected: FAIL because GatewayConfig and InRollout are undefined.

- [ ] **Step 3: Implement configuration values and stable rollout**

Validate returns fmt.Errorf("%w: field", ErrInvalidConfig). InRollout hashes the key with fnv.New32a and compares hash%100 with percent; special-case 0 and 100 before hashing.

~~~go
type GatewayConfig struct {
    RouteEnabled  bool
    RequestTimeout time.Duration
    RateLimit      int
    RateWindow     time.Duration
    RolloutPercent uint32
    BearerToken    string
}

type Snapshot struct {
    Version uint64
    Config  GatewayConfig
}
~~~

- [ ] **Step 4: Write failing memory-store tests**

Cover initial version 1, Current immutable value, valid update increments exactly once, invalid update preserves version, immediate Watch snapshot, latest-only updates for a slow watcher, canceled watcher, Close, and concurrent Current/Update.

~~~go
func TestMemoryStoreRejectsInvalidUpdateWithoutAdvancingVersion(t *testing.T) {
    store, err := NewMemoryStore(validConfig())
    if err != nil {
        t.Fatal(err)
    }
    defer store.Close()
    before, _ := store.Current()
    invalid := before.Config
    invalid.RolloutPercent = 101
    if _, err := store.Update(invalid); !errors.Is(err, ErrInvalidConfig) {
        t.Fatalf("error = %v", err)
    }
    after, _ := store.Current()
    if after.Version != before.Version {
        t.Fatalf("version advanced: before=%d after=%d", before.Version, after.Version)
    }
}
~~~

- [ ] **Step 5: Run store tests and verify RED**

Run:

~~~powershell
go test ./stage-3-architecture/14-microservices/internal/configcenter -run TestMemoryStore -count=1
~~~

Expected: FAIL because NewMemoryStore and Store methods are undefined.

- [ ] **Step 6: Implement the memory configuration store**

Use the same one-element latest-snapshot watcher pattern as discovery, but do not share implementation across packages. NewMemoryStore validates and creates version 1. Update validates before locking, checks closed, increments version, and publishes non-blocking while holding the mutex. Cancellation and idempotent Close remove or close watcher channels under that same mutex so a publisher can never send to a closed channel.

- [ ] **Step 7: Run focused and race tests**

Run:

~~~powershell
gofmt -w stage-3-architecture/14-microservices/internal/configcenter
go test ./stage-3-architecture/14-microservices/internal/configcenter -count=1
go test -race -count=1 ./stage-3-architecture/14-microservices/internal/configcenter
~~~

Expected: PASS with no race reports.

- [ ] **Step 8: Commit configuration**

~~~powershell
git add stage-3-architecture/14-microservices/internal/configcenter
git commit -m "feat: add dynamic gateway configuration"
~~~

### Task 6: HTTP API Gateway

**Files:**
- Create: stage-3-architecture/14-microservices/internal/gateway/auth.go
- Create: stage-3-architecture/14-microservices/internal/gateway/auth_test.go
- Create: stage-3-architecture/14-microservices/internal/gateway/limiter.go
- Create: stage-3-architecture/14-microservices/internal/gateway/limiter_test.go
- Create: stage-3-architecture/14-microservices/internal/gateway/connections.go
- Create: stage-3-architecture/14-microservices/internal/gateway/connections_test.go
- Create: stage-3-architecture/14-microservices/internal/gateway/errors.go
- Create: stage-3-architecture/14-microservices/internal/gateway/handler.go
- Create: stage-3-architecture/14-microservices/internal/gateway/handler_test.go
- Create: stage-3-architecture/14-microservices/internal/gateway/integration_test.go

**Interfaces:**
- Consumes: discovery Registry Resolve, configcenter Store Current, generated product and inventory clients.
- Produces: NewLimiter(clock func() time.Time) *Limiter and (*Limiter).Allow(key string, limit int, window time.Duration) bool.
- Produces: Connections.Product(context.Context, string) (productv1.ProductServiceClient, error), Inventory(context.Context, string) (inventoryv1.InventoryServiceClient, error), Close() error.
- Produces: NewHandler(ConfigReader, Resolver, ClientProvider, *Limiter) http.Handler.
- HTTP route: GET /api/v1/products/{sku}; stable request key is X-Request-Key, falling back to RemoteAddr.

- [ ] **Step 1: Write failing authentication and limiter tests**

Authentication accepts exactly the configured token in the Authorization: Bearer header. Limiter uses a fixed window per request key; changing limit/window applies on the next Allow call. Cover first allowed, limit exhausted, separate keys, and window reset with an injected clock.

~~~go
func TestLimiterResetsAfterWindow(t *testing.T) {
    now := time.Unix(100, 0)
    limiter := NewLimiter(func() time.Time { return now })
    if !limiter.Allow("client-a", 1, time.Minute) {
        t.Fatal("first request rejected")
    }
    if limiter.Allow("client-a", 1, time.Minute) {
        t.Fatal("second request allowed")
    }
    now = now.Add(time.Minute)
    if !limiter.Allow("client-a", 1, time.Minute) {
        t.Fatal("request after reset rejected")
    }
}
~~~

- [ ] **Step 2: Run policy tests and verify RED**

Run:

~~~powershell
go test ./stage-3-architecture/14-microservices/internal/gateway -run 'Test(Auth|Limiter)' -count=1
~~~

Expected: FAIL because the gateway package does not exist.

- [ ] **Step 3: Implement authentication and local limiter**

Use subtle.ConstantTimeCompare for token comparison after exact Bearer prefix parsing. Limiter owns map[string]bucket under sync.Mutex; bucket stores window start and count. It rejects invalid limit/window defensively. Document in Go comments that state is process-local.

- [ ] **Step 4: Write failing connection ownership tests**

Use a fake dial function that returns a real bufconn-backed *grpc.ClientConn and increments a count. Assert repeated address lookup returns clients backed by one connection, distinct addresses dial separately, and Close is idempotent and causes later calls to return ErrConnectionsClosed.

Required constructor:

~~~go
type DialFunc func(context.Context, string) (*grpc.ClientConn, error)

func NewConnections(dial DialFunc) *Connections
~~~

Production dial uses grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials())).

- [ ] **Step 5: Run connection tests and verify RED**

Run:

~~~powershell
go test ./stage-3-architecture/14-microservices/internal/gateway -run TestConnections -count=1
~~~

Expected: FAIL because Connections is undefined.

- [ ] **Step 6: Implement cached connections**

Connections owns map[string]*grpc.ClientConn under sync.Mutex. Product and Inventory call a private conn(ctx,address), validate address, reuse existing connection, and avoid storing a newly dialed connection if Close raced with dialing. Close marks closed, copies and clears the map, unlocks, and joins close errors with errors.Join.

- [ ] **Step 7: Write failing handler tests with fakes**

Define small test fakes for ConfigReader, Resolver, and ClientProvider. Cover:

- missing/wrong bearer token => 401 and zero downstream calls;
- limiter rejection => 429 and zero downstream calls;
- route disabled or rollout rejected => 404 and zero downstream calls;
- missing service => 503;
- success => 200 JSON with sku, name, price_cents, quantity, stock_version;
- product and inventory calls both start before either fake is released, proving concurrency;
- InvalidArgument/NotFound/Unavailable/DeadlineExceeded => 400/404/503/504;
- unknown error => 500 with generic body;
- one downstream failure => no product fields in the error response.

Test-facing interfaces:

~~~go
type ConfigReader interface {
    Current() (configcenter.Snapshot, error)
}

type Resolver interface {
    Resolve(service string) (discovery.Instance, error)
}

type ClientProvider interface {
    Product(context.Context, string) (productv1.ProductServiceClient, error)
    Inventory(context.Context, string) (inventoryv1.InventoryServiceClient, error)
}

func NewHandler(config ConfigReader, resolver Resolver, clients ClientProvider, limiter *Limiter) http.Handler
~~~

- [ ] **Step 8: Run handler tests and verify RED**

Run:

~~~powershell
go test ./stage-3-architecture/14-microservices/internal/gateway -run TestHandler -count=1
~~~

Expected: FAIL because NewHandler and error mapping are undefined.

- [ ] **Step 9: Implement Gateway routing and aggregation**

Register GET /api/v1/products/{sku} on http.NewServeMux. Per request:

1. Read current config; closed/unavailable config => 503.
2. Authenticate, then rate-limit by request key.
3. Reject disabled route or failed InRollout decision with 404.
4. Validate non-empty PathValue("sku") or return 400.
5. Resolve service names "product" and "inventory".
6. Create context.WithTimeout using RequestTimeout.
7. Start two goroutines writing typed results to two buffered channels.
8. Read both results, cancel on failure, and map the first required error.
9. Encode exactly one application/json success object; unknown error bodies use {"error":"request failed"}.

Use a response struct with explicit JSON field tags:

~~~go
type productDetails struct {
    SKU          string `json:"sku"`
    Name         string `json:"name"`
    PriceCents   int64  `json:"price_cents"`
    Quantity     int64  `json:"quantity"`
    StockVersion uint64 `json:"stock_version"`
}
~~~

- [ ] **Step 10: Add full HTTP-over-real-gRPC integration test**

Start product and inventory gRPC servers on separate bufconn listeners, register their addresses in discovery, create the real config store, Connections, and Handler, then call httptest.NewServer. Assert a real HTTP request receives the aggregate response and a canceled inventory service maps to 504. This test must use generated gRPC clients, not fake service interfaces.

- [ ] **Step 11: Run gateway tests and race detector**

Run:

~~~powershell
gofmt -w stage-3-architecture/14-microservices/internal/gateway
go test ./stage-3-architecture/14-microservices/internal/gateway -count=1
go test -race -count=1 ./stage-3-architecture/14-microservices/internal/gateway
~~~

Expected: PASS with no races, leaks, or flaky timing assertions.

- [ ] **Step 12: Commit Gateway**

~~~powershell
git add stage-3-architecture/14-microservices/internal/gateway
git commit -m "feat: add microservices api gateway"
~~~

### Task 7: Composition Root, Documentation, and Curriculum

**Files:**
- Create: stage-3-architecture/14-microservices/app.go
- Create: stage-3-architecture/14-microservices/app_test.go
- Create: stage-3-architecture/14-microservices/main.go
- Replace: stage-3-architecture/14-microservices/README.md
- Create: stage-3-architecture/14-microservices/EXERCISES.md
- Delete: stage-3-architecture/14-microservices/.gitkeep
- Modify: ROADMAP.md
- Modify: openspec/changes/chapter-14-microservices/tasks.md as each item is completed

**Interfaces:**
- Produces: run(context.Context, io.Writer) error for deterministic lifecycle testing.
- Runtime output contains an aggregate product result and no nondeterministic port numbers.
- Owns and closes configuration store, discovery registry, deregistration handles, Gateway connections, HTTP server, gRPC servers, and listeners.

- [ ] **Step 1: Write the failing end-to-end lifecycle test**

~~~go
func TestRunCompletesProductDetailsFlow(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    var out bytes.Buffer
    if err := run(ctx, &out); err != nil {
        t.Fatal(err)
    }
    got := out.String()
    for _, want := range []string{"book-1", "Go Microservices", "quantity"} {
        if !strings.Contains(got, want) {
            t.Fatalf("output %q missing %q", got, want)
        }
    }
}
~~~

Also add TestRunHonorsCanceledContext with a pre-canceled context and require errors.Is(err, context.Canceled).

- [ ] **Step 2: Run lifecycle tests and verify RED**

Run:

~~~powershell
go test ./stage-3-architecture/14-microservices -run TestRun -count=1
~~~

Expected: FAIL because run is undefined.

- [ ] **Step 3: Implement app lifecycle and main**

app.go creates fixed teaching data, starts separate product and inventory grpc.Server values on 127.0.0.1:0 listeners, registers addresses, creates version-1 config, starts an http.Server on another random listener, performs an authorized GET with X-Request-Key, decodes and prints stable JSON, then shuts down in reverse ownership order. Every start goroutine sends its Serve error to a buffered channel; http.ErrServerClosed and grpc.ErrServerStopped are normal. Use defer immediately after each successful allocation so every early return releases prior resources.

main.go:

~~~go
package main

import (
    "context"
    "fmt"
    "os"
)

func main() {
    if err := run(context.Background(), os.Stdout); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
~~~

- [ ] **Step 4: Run executable and lifecycle tests**

Run:

~~~powershell
go test ./stage-3-architecture/14-microservices -count=1
go run ./stage-3-architecture/14-microservices
~~~

Expected: tests PASS; the executable prints stable JSON containing book-1, Go Microservices, price, quantity, and stock version, then exits.

- [ ] **Step 5: Replace README with the complete Chinese tutorial**

README must contain:

- learning goals and concrete outputs;
- exact package tree;
- protobuf compatibility rules and pinned generation commands;
- unary/server-streaming/bidirectional-streaming comparison;
- discovery registration/resolution/watch lifecycle;
- versioned dynamic config and deterministic rollout;
- Gateway request sequence, auth, local rate limit, aggregation, status mapping;
- synchronous gRPC versus asynchronous MQ decision table;
- Mermaid or text architecture and sequence diagrams;
- run, focused test, full test, race, vet, build, lint, and OpenSpec commands;
- explicit production limitations and self-check checklist.

Every referenced command must be runnable from the repository root unless the text explicitly changes directory.

- [ ] **Step 6: Add exercises with measurable acceptance criteria**

EXERCISES.md contains seven exercises: add client streaming, implement health-aware discovery, create Consul/etcd adapter, add TLS/mTLS, persist/version config, design distributed limiting, and replace one synchronous flow with MQ. Each exercise lists goal, constraints, verification command or observable result, and edge cases.

- [ ] **Step 7: Update ROADMAP and OpenSpec checkboxes**

Replace the Chapter 14 placeholder output with concrete product/inventory gRPC, discovery/configuration, and Gateway outputs. Mark only 14-microservices complete. Check every completed task in openspec/changes/chapter-14-microservices/tasks.md; do not check verification/review/issue/PR tasks until their evidence exists.

- [ ] **Step 8: Run documentation and chapter checks**

Run:

~~~powershell
rg -n "待 OpenSpec change 填充|落地后填充" stage-3-architecture/14-microservices/README.md
go test ./stage-3-architecture/14-microservices/... -count=1
go run ./stage-3-architecture/14-microservices
openspec validate chapter-14-microservices --strict
git diff --check
~~~

Expected: rg returns no Chapter 14 placeholder match; tests, executable, OpenSpec, and diff check succeed.

- [ ] **Step 9: Commit composition and materials**

~~~powershell
git add ROADMAP.md stage-3-architecture/14-microservices openspec/changes/chapter-14-microservices/tasks.md
git commit -m "docs: complete chapter 14 microservices tutorial"
~~~

### Task 8: Repository Verification, Requested Review, Issue, and PR

**Files:**
- Modify only files required to fix verified failures or Critical/Important review findings.
- Modify: openspec/changes/chapter-14-microservices/tasks.md
- Read: docs/superpowers/specs/2026-07-12-chapter-14-microservices-design.md
- Read: openspec/changes/chapter-14-microservices/specs/**/*.md

**Interfaces:**
- Consumes the full origin/main...HEAD change.
- Produces fresh verification evidence, one subagent code review, a linked GitHub issue, a pushed branch, and a PR against main.

- [ ] **Step 1: Run all fresh quality gates**

Run separately and record exit codes:

~~~powershell
gofmt -w stage-3-architecture/14-microservices
go test ./stage-3-architecture/14-microservices/... -count=1
go test ./... -count=1
go vet ./...
go test -race -count=1 ./...
go build ./...
golangci-lint run
openspec validate chapter-14-microservices --strict
Push-Location stage-3-architecture/14-microservices
go run github.com/bufbuild/buf/cmd/buf@v1.65.0 generate
Pop-Location
git diff --check
~~~

Expected: every available command exits 0, race reports none, OpenSpec is valid, and regeneration adds no generated-code diff. If golangci-lint is unavailable, record the exact command-not-found evidence rather than claiming it passed.

- [ ] **Step 2: Audit every requirement against evidence**

Create a temporary in-memory checklist mapping each scenario in microservices-foundations-tutorial/spec.md and learning-curriculum/spec.md to a test, executable output, generated descriptor, or file diff. Add missing evidence before proceeding. Re-read the approved design and verify every goal, non-goal, lifecycle rule, status mapping, and generated-file rule.

- [ ] **Step 3: Commit any verification fixes**

After reproducing each failure with a focused test, fix it test-first, rerun the focused and full relevant gates, update task checkboxes, then:

~~~powershell
git add go.mod go.sum ROADMAP.md stage-3-architecture/14-microservices openspec/changes/chapter-14-microservices
git commit -m "fix: address chapter 14 verification findings"
~~~

Omit this commit when no files changed.

- [ ] **Step 4: Dispatch the requested code-review subagent**

Calculate:

~~~powershell
git rev-parse origin/main
git rev-parse HEAD
~~~

Dispatch one reviewer with fork_turns="none" and the code-reviewer template. Provide the approved design path, OpenSpec paths, exact BASE_SHA and HEAD_SHA, verification output, and instruct it to inspect origin/main...HEAD for correctness, security, concurrency, lifecycle, generated-source consistency, teaching clarity, and requirement coverage. The reviewer must report Critical, Important, and Minor findings with file and line evidence.

- [ ] **Step 5: Resolve review findings and re-verify**

Fix every valid Critical or Important finding with a failing regression test first. Technically rebut incorrect findings with code/test evidence. Minor findings may be fixed when scoped or listed in the PR. Rerun the complete Step 1 gate after all fixes and update OpenSpec tasks.

- [ ] **Step 6: Create the GitHub issue**

Use gh issue create with a Chinese title and body containing motivation, approved scope, acceptance checklist, verification evidence, and the branch name. Capture the issue number and URL. Do not duplicate an existing Chapter 14 issue; search first with:

~~~powershell
gh issue list --state all --limit 100 --search "chapter 14 microservices in:title" --json number,title,state,url
~~~

When no matching issue exists, run:

~~~powershell
$issueURL = gh issue create --title "实现第 14 章：微服务基础设施" --body "## 目标`n完成真实 gRPC、服务发现、动态配置和 API Gateway 教学章节。`n`n## 验收`n- [ ] 三种 RPC 形态有测试`n- [ ] Gateway 覆盖鉴权、限流、聚合和错误映射`n- [ ] 全仓质量门通过`n- [ ] 文档与 OpenSpec 同步`n`n实现分支：codex/chapter-14-microservices"
$issueNumber = [int]($issueURL -replace '.*/','')
~~~

- [ ] **Step 7: Final commit, push, and PR**

Commit the completed OpenSpec task checkboxes and any issue reference, push without force, then create a non-draft PR:

~~~powershell
git push -u origin codex/chapter-14-microservices
$prBody = "Closes #$issueNumber`n`n## Summary`n- add product and inventory protobuf/gRPC services`n- add discovery, dynamic configuration, and HTTP Gateway`n- add runnable tutorial, tests, exercises, and OpenSpec artifacts`n`n## Verification`n- go test ./... -count=1`n- go vet ./...`n- go test -race -count=1 ./...`n- go build ./...`n- golangci-lint run`n- openspec validate chapter-14-microservices --strict"
gh pr create --base main --head codex/chapter-14-microservices --title "feat: complete chapter 14 microservices" --body $prBody
~~~

The PR body links the issue with Closes #N, summarizes architecture and learning outputs, lists every actually executed verification command and limitation, and notes the subagent review outcome. Verify the returned PR with gh pr view --json number,title,state,url,baseRefName,headRefName.
