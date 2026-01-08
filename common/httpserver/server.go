package httpserver

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

// Server represents a HTTP server that can handle requests for this microservice.
type Server struct {
	routes        []Route
	authenticator Authenticator
	logger        *zap.Logger
	serviceName   string
	metrics       http.Handler
	ginMode       string
	openAPI       *OpenAPI
}

// New will instantiate a new instance of Server.
func New(routes []Route, opts ...func(*Server)) *Server {
	s := &Server{
		routes:      routes,
		logger:      zap.NewNop(),
		serviceName: "",
		ginMode:     gin.ReleaseMode,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithGinDebug toggles gin debug mode (otherwise release mode).
func WithGinDebug(enabled bool) func(*Server) {
	return func(s *Server) {
		if enabled {
			s.ginMode = gin.DebugMode
		} else {
			s.ginMode = gin.ReleaseMode
		}
	}
}

// WithLogger configures the structured logger used by the router.
func WithLogger(l *zap.Logger) func(*Server) {
	return func(s *Server) {
		if l != nil {
			s.logger = l
		}
	}
}

// WithServiceName sets the service name used for telemetry.
func WithServiceName(name string) func(*Server) {
	return func(s *Server) {
		if name != "" {
			s.serviceName = name
		}
	}
}

// WithOpenAPI sets the OpenAPI helper used to document routes.
func WithOpenAPI(openAPI *OpenAPI) func(*Server) {
	return func(s *Server) {
		if openAPI != nil {
			s.openAPI = openAPI
		}
	}
}

// WithMetricsHandler wires a metrics endpoint when provided.
func WithMetricsHandler(h http.Handler) func(*Server) {
	return func(s *Server) {
		s.metrics = h
	}
}

// WithAuthenticator configures request authentication (e.g., client cert).
func WithAuthenticator(a Authenticator) func(*Server) {
	return func(s *Server) {
		s.authenticator = a
	}
}

// Router constructs the gin engine with middleware and routes.
func (s *Server) Router() *gin.Engine {
	gin.SetMode(s.ginMode)
	engine := gin.New()

	// Observability middleware stack: request ID, tracing, structured logs, recovery.
	engine.Use(RequestIDMiddleware())
	engine.Use(ginzap.Ginzap(s.logger, time.RFC3339Nano, true))
	engine.Use(ginzap.RecoveryWithZap(s.logger, true))
	engine.Use(otelgin.Middleware(s.serviceName))

	var secureGroup *gin.RouterGroup
	if s.authenticator != nil {
		secureGroup = engine.Group("")
		secureGroup.Use(AuthMiddleware(s.authenticator))
	}

	if s.metrics != nil {
		engine.GET("/metrics", gin.WrapH(s.metrics))
	}

	for _, rt := range s.routes {
		if rt.Handler == nil {
			s.logger.Warn("missing handler for route", zap.String("path", rt.GinPath))
			continue
		}
		var target gin.IRoutes = engine
		if rt.Secure && secureGroup != nil {
			target = secureGroup
		}
		target.Handle(rt.Method, rt.GinPath, rt.Handler)
		if s.openAPI != nil {
			s.openAPI.AddRoute(rt.RouteSpec)
		}
	}

	if s.openAPI != nil {
		s.openAPI.AttachDocs(engine, "{{API_SPEC_FILE}}", "{{API_DOC_ROOT}}")
	}

	return engine
}
