package middleware

import (
	"context"
	"github.com/simabdi/auth-service/contextkey"
	"github.com/simabdi/auth-service/service"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

//type contextKey struct {
//	name string
//}
//
//var (
//	contextUserID = &contextKey{"user_id"}
//	contextUUID   = &contextKey{"uuid"}
//	contextRef    = &contextKey{"ref"}
//	contextRefID  = &contextKey{"ref_id"}
//)

func Authenticate(authService service.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractBearerToken(r.Header.Get("Authorization"))
			if token == "" {
				http.Error(w, "Unauthorized: Missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}

			authInfo, err := authService.ParseToken(r.Context(), token)
			if err != nil {
				http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
				return
			}

			//ctx := context.WithValue(r.Context(), contextUserID, authInfo.UserID)
			//ctx = context.WithValue(ctx, contextUUID, authInfo.Uuid)
			//ctx = context.WithValue(ctx, contextRef, authInfo.RefType)
			//ctx = context.WithValue(ctx, contextRefID, authInfo.RefID)

			ctx := context.WithValue(r.Context(), contextkey.UserIDKey, authInfo.UserID)
			ctx = context.WithValue(ctx, contextkey.UUIDKey, authInfo.Uuid)
			ctx = context.WithValue(ctx, contextkey.RefKey, authInfo.RefType)
			ctx = context.WithValue(ctx, contextkey.RefIDKey, authInfo.RefID)

			if setter, ok := w.(interface{ SetContext(context.Context) }); ok {
				log.Debug("Injecting user context to response writer")
				setter.SetContext(ctx)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractBearerToken(header string) string {
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
}
