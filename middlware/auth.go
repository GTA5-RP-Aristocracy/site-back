package middlware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	siteback "github.com/GTA5-RP-Aristocracy/site-back"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func AuthMiddleware(secret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check the token.
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			siteback.WriteError(w, errors.New("missing token"), http.StatusUnauthorized)
			return
		}

		tokenStr = tokenStr[7:]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(secret), nil
		})
		if err != nil {
			siteback.WriteError(w, fmt.Errorf("pars jwt: %w", err), http.StatusUnauthorized)
			return
		}

		// Check if the token is valid.
		if !token.Valid {
			siteback.WriteError(w, errors.New("invalid token"), http.StatusUnauthorized)
			return
		}

		// check if the token is expired
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			siteback.WriteError(w, errors.New("invalid token claims"), http.StatusUnauthorized)
			return
		}

		if err := claims.Valid(); err != nil {
			siteback.WriteError(w, fmt.Errorf("invalid claims: %w", err), http.StatusUnauthorized)
			return
		}

		exp := int64(claims["exp"].(float64))
		if exp < 0 {
			siteback.WriteError(w, errors.New("token expired"), http.StatusUnauthorized)
			return
		}

		expDuration := time.Unix(exp, 0)
		if time.Now().After(expDuration) {
			siteback.WriteError(w, errors.New("token expired"), http.StatusUnauthorized)
			return
		}

		// Add the user ID to the context.
		userID, err := uuid.Parse(claims["id"].(string))
		if err != nil {
			siteback.WriteError(w, fmt.Errorf("parse id err: %w", err), http.StatusUnauthorized)
			return
		}

		r = r.WithContext(siteback.ContextWithUserID(r.Context(), userID))

		// Call the next handler.
		next.ServeHTTP(w, r)
	}
}
