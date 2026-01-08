package httpserver

import (
	"context"

	commonhttp "{{MODULE_PATH}}/common/httpserver"
)

const (
	defaultServiceName    = "{{SERVICE_NAME}}"
	defaultAPITitle       = "{{SERVICE_NAME}} API"
	defaultAPIDescription = "API for {{SERVICE_NAME}}"
	defaultAPIVersion     = "1.0.0"
)

// New constructs a HTTP server with app defaults (logger, service name, docs) layered on the shared builder.
func New(routes []commonhttp.Route, opts ...func(*commonhttp.Server)) *commonhttp.Server {
	logger := From(context.Background())
	openAPI := commonhttp.NewOpenAPI(logger, defaultAPITitle, defaultAPIVersion, defaultAPIDescription)

	defaultOpts := []func(*commonhttp.Server){
		commonhttp.WithLogger(logger),
		commonhttp.WithServiceName(defaultServiceName),
		commonhttp.WithOpenAPI(openAPI),
	}

	defaultOpts = append(defaultOpts, opts...)

	return commonhttp.New(routes, defaultOpts...)
}
