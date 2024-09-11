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

func (h *Handler) AuthMiddleware(next http.Handler, userServiceConn *grpc.ClientConn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			http.Error(w, "authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

		// Parse JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extract user claims (e.g., user_id)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
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
