package http

import (
	"context"
	"fmt"
	"gold-gym-be/internal/entity"
	"gold-gym-be/pkg/response"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func (s *Server) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			resp := &response.Response{}
			defer resp.RenderJSON(w, r)

			resp.Error = response.Error{
				Status: false,
				Msg:    "Invalid token: unsupported token type",
				Code:   403,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(403)

			return
		}

		token := strings.Split(authorization, " ")
		if token[0] != "Bearer" {
			resp := &response.Response{}
			defer resp.RenderJSON(w, r)

			resp.Error = response.Error{
				Status: false,
				Msg:    "Invalid token: unsupported token type",
				Code:   403,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(403)

			return
		}

		jwtToken, err := jwt.Parse(token[1], func(_token *jwt.Token) (interface{}, error) {
			if method, ok := _token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("signing method invalid - a")
			} else if method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("signing method invalid - b")
			}

			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})
		if err != nil {
			resp := &response.Response{}
			defer resp.RenderJSON(w, r)

			resp.Error = response.Error{
				Status: false,
				Msg:    err.Error(),
				Code:   500,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)

			return
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			resp := &response.Response{}
			defer resp.RenderJSON(w, r)

			resp.Error = response.Error{
				Status: false,
				Msg:    "Invalid token: unsupported token type",
				Code:   401,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(401)

			return
		}
		if !jwtToken.Valid {
			resp := &response.Response{}
			defer resp.RenderJSON(w, r)

			resp.Error = response.Error{
				Status: false,
				Msg:    "Invalid token: unsupported token type",
				Code:   401,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(401)

			return
		}

		// do something with decoded claims
		for key, val := range claims {
			if key != "permissions" {
				continue
			}
			ctxVal := entity.ContextValue{
				M: map[string]interface{}{
					key: val,
				},
			}
			r = r.WithContext(context.WithValue(r.Context(), entity.ContextKey("claims"), ctxVal))
		}

		next.ServeHTTP(w, r)
	})
}
