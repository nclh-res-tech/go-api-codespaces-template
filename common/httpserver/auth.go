package httpserver

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authenticator validates an incoming HTTP request.
type Authenticator interface {
	Authenticate(ctx context.Context, r *http.Request) error
}

// CertAuthenticator validates client certificates by fingerprint or common name.
type CertAuthenticator struct {
	allowedFingerprints map[string]struct{}
	allowedCNs          map[string]struct{}
}

// NewCertAuthenticator constructs a cert authenticator.
// Provide SHA256 fingerprints (hex) and/or CommonNames that are allowed.
func NewCertAuthenticator(fingerprints []string, commonNames []string) Authenticator {
	fpMap := make(map[string]struct{}, len(fingerprints))
	for _, fp := range fingerprints {
		fpMap[strings.ToLower(fp)] = struct{}{}
	}

	cnMap := make(map[string]struct{}, len(commonNames))
	for _, cn := range commonNames {
		cnMap[cn] = struct{}{}
	}

	return &CertAuthenticator{
		allowedFingerprints: fpMap,
		allowedCNs:          cnMap,
	}
}

// Authenticate validates the client certificate against allowed fingerprints or CNs.
func (a *CertAuthenticator) Authenticate(_ context.Context, r *http.Request) error {
	if a == nil {
		return ErrUnauthorized
	}

	if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
		return ErrUnauthorized
	}

	for _, cert := range r.TLS.PeerCertificates {
		if len(a.allowedFingerprints) > 0 {
			fp := sha256.Sum256(cert.Raw)
			if _, ok := a.allowedFingerprints[strings.ToLower(hex.EncodeToString(fp[:]))]; ok {
				return nil
			}
		}
		if len(a.allowedCNs) > 0 {
			if _, ok := a.allowedCNs[cert.Subject.CommonName]; ok {
				return nil
			}
		}
	}

	return ErrUnauthorized
}

// CompositeAuthenticator succeeds if any underlying authenticator succeeds.
type CompositeAuthenticator struct {
	authenticators []Authenticator
}

// NewCompositeAuthenticator returns an authenticator that tries each provided one.
func NewCompositeAuthenticator(auths ...Authenticator) Authenticator {
	return &CompositeAuthenticator{authenticators: auths}
}

// Authenticate checks each configured authenticator until one succeeds.
func (c *CompositeAuthenticator) Authenticate(ctx context.Context, r *http.Request) error {
	for _, a := range c.authenticators {
		if a == nil {
			continue
		}
		if err := a.Authenticate(ctx, r); err == nil {
			return nil
		}
	}
	return ErrUnauthorized
}

// AuthMiddleware returns a gin middleware wrapper that applies the authenticator when configured.
func AuthMiddleware(a Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		if a == nil {
			c.Next()
			return
		}

		if err := a.Authenticate(c.Request.Context(), c.Request); err != nil {
			HandleError(c, err)
			c.Abort()
			return
		}

		c.Next()
	}
}
