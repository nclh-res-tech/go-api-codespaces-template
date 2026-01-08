package httpserver

import (
	"context"

	commonhttp "{{MODULE_PATH}}/common/httpserver"
)

const (
	defaultServiceName    = "{{ .ServiceName }}"
	defaultAPITitle       = "{{ .APITitle }}"
	defaultAPIDescription = "{{ .APIDescription }}"
	defaultAPIVersion     = "{{ .APIVersion }}"
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
