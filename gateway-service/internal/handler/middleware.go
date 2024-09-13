package handler

import (
	"context"
	"net/http"
	"os"
	"strings"

	userpb "github.com/emzola/numer/user-service/proto"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
)

func (h *Handler) authMiddleware(next http.HandlerFunc, userServiceConn *grpc.ClientConn) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			h.invalidCredentialsResponse(w, r)
			return
		}

		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

		// Parse JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			h.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Extract user claims (e.g., user_id)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			h.invalidAuthenticationTokenResponse(w, r)
			return
		}

		userID := claims["user_id"].(int64)

		// Fetch user details
		userClient := userpb.NewUserServiceClient(userServiceConn)
		resp, err := userClient.GetUser(context.Background(), &userpb.GetUserRequest{UserId: userID})
		if err != nil {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}

		// Add user details to context
		r = h.contextSetUser(r, resp.User)

		// Proceed to next handler
		next.ServeHTTP(w, r)
	})
}
