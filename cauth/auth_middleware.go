package cauth

import (
	"net/http"

	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
)

// AuthMiddleware provides a middleware that verifies the auth header using basic auth.
// The username is expected to be the email and the password should be the session token.
// On success, the user is stored in the context.
type AuthMiddleware interface {
	AllowVerified(next http.Handler) http.Handler
	AllowUnverified(next http.Handler) http.Handler
}

type authMiddleware struct {
	resp   chttp.Responder
	users  UsersSvc
	logger clogger.Logger
}

func newAuthMiddleware(resp chttp.Responder, users UsersSvc, logger clogger.Logger) AuthMiddleware {
	return &authMiddleware{
		resp:   resp,
		users:  users,
		logger: logger,
	}
}

func (m *authMiddleware) AllowVerified(next http.Handler) http.Handler {
	return m.verifyAuth(next, false)
}

func (m *authMiddleware) AllowUnverified(next http.Handler) http.Handler {
	return m.verifyAuth(next, true)
}

func (m *authMiddleware) verifyAuth(next http.Handler, allowUnverified bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, sessionToken, ok := r.BasicAuth()
		if !ok {
			m.resp.Unauthorized(w)
			return
		}

		ctx := r.Context()

		user, err := m.users.VerifySessionToken(ctx, email, sessionToken)
		if err != nil && err != ErrInvalidCredentials {
			m.logger.Error("Failed to verify user session token", err)
			m.resp.InternalErr(w)
			return
		} else if err == ErrInvalidCredentials {
			m.resp.Unauthorized(w)
			return
		}

		if !allowUnverified && user.Verified == false {
			m.resp.Unauthorized(w)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctxWithUser(ctx, user)))
	})
}
