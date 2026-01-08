package httpserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/swgui/v5cdn"
	"go.uber.org/zap"
)

// RouteSpec describes a single route and the OpenAPI metadata needed to register it.
type RouteSpec struct {
	Method          string
	GinPath         string
	SpecPath        string
	Summary         string
	Tags            []string
	Req             any
	Resp            any
	ReqContentType  string
	RespContentType string
	Secure          bool
}

// OpenAPI holds the OpenAPI reflector and related settings for registration.
type OpenAPI struct {
	ref    *openapi31.Reflector
	logger *zap.Logger
}

// OpenAPIOption configures the OpenAPI helper.
type OpenAPIOption func(*OpenAPI)

// NewOpenAPI constructs an OpenAPI helper with basic API info.
func NewOpenAPI(logger *zap.Logger, title, version, description string, opts ...OpenAPIOption) *OpenAPI {
	ref := openapi31.NewReflector()
	ref.Spec.Info.Title = title
	ref.Spec.Info.Version = version
	if description != "" {
		ref.Spec.Info.Description = &description
	}

	o := &OpenAPI{
		ref:    ref,
		logger: logger,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

// AddRoute registers the OpenAPI operation for the provided route.
func (o *OpenAPI) AddRoute(route RouteSpec) {
	specPath := route.SpecPath
	if specPath == "" {
		specPath = route.GinPath
	}

	oc, err := o.ref.NewOperationContext(route.Method, specPath)
	if err != nil {
		o.logger.Warn("failed to create operation context", zap.Error(err), zap.String("method", route.Method), zap.String("path", specPath))
		return
	}
	if route.Summary != "" {
		oc.SetSummary(route.Summary)
	}
	if len(route.Tags) > 0 {
		oc.SetTags(route.Tags...)
	}
	if route.Req != nil {
		oc.AddReqStructure(route.Req)
	} else if strings.TrimSpace(route.ReqContentType) != "" {
		type rawBody struct {
			Raw string `json:"raw"`
		}
		oc.AddReqStructure(new(rawBody), openapi.WithContentType(route.ReqContentType))
	}
	if route.Resp != nil {
		oc.AddRespStructure(route.Resp)
	} else if strings.TrimSpace(route.RespContentType) != "" {
		type rawResp struct {
			Raw string `json:"raw"`
		}
		oc.AddRespStructure(new(rawResp), openapi.WithContentType(route.RespContentType))
	}

	if err := o.ref.AddOperation(oc); err != nil {
		o.logger.Warn("failed to add operation", zap.Error(err), zap.String("method", route.Method), zap.String("path", specPath))
	}
}

// AttachDocs exposes the OpenAPI JSON and Swagger UI endpoints.
func (o *OpenAPI) AttachDocs(engine *gin.Engine, jsonPath, docsPath string) {
	if jsonPath == "" {
		jsonPath = "/openapi.json"
	}
	if docsPath == "" {
		docsPath = "/docs"
	}

	specBytes, err := json.MarshalIndent(o.ref.Spec, "", "  ")
	if err != nil {
		o.logger.Warn("failed to marshal openapi spec", zap.Error(err))
		return
	}

	engine.GET(jsonPath, func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", specBytes)
	})
	ui := v5cdn.NewHandler("API Docs", jsonPath, docsPath)
	engine.GET(docsPath+"/*any", gin.WrapH(ui))
	engine.GET(docsPath, gin.WrapH(ui))
}
