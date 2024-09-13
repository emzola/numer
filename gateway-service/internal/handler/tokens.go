package handler

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/emzola/numer/gateway-service/internal/grpcutil"
	userpb "github.com/emzola/numer/user-service/proto"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) AuthenticateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthHTTPReq
	err := h.decodeJSON(w, r, &req)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Create gRPC connection to user service
	conn, err := grpcutil.ServiceConnection(ctx, "user-service", h.registry)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)

	authReq := &userpb.AuthenticateUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	authResp, err := client.AuthenticateUser(context.Background(), authReq)
	if err != nil {
		h.invalidCredentialsResponse(w, r)
		return
	}

	// Generate JWT token upon successful authentication
	token, err := GenerateJWTToken(authResp.UserId, authResp.Email, authResp.Role)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	err = h.encodeJSON(w, http.StatusOK, envelope{"token": token}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(userID int64, email, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Struct to capture the HTTP request JSON data
type AuthHTTPReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Struct to capture the HTTP request response
type TokenResponse struct {
	Token string `json:"token"`
}
