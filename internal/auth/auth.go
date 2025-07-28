package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Error string `json:"error,omitempty"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func init() {
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("default-secret-key")
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[AUTH] Login request from: %s\n", r.RemoteAddr)
	fmt.Printf("[AUTH] Request method: %s\n", r.Method)
	
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		fmt.Printf("[AUTH] Method not allowed: %s\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(LoginResponse{Error: "Method not allowed"})
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("[AUTH] Failed to decode request body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(LoginResponse{Error: "Invalid request body"})
		return
	}
	
	fmt.Printf("[AUTH] Login attempt for username: %s\n", req.Username)

	// Check credentials against environment variables
	expectedUsername := os.Getenv("ADMIN_USERNAME")
	expectedPassword := os.Getenv("ADMIN_PASSWORD")
	
	fmt.Printf("[AUTH] Expected username: '%s'\n", expectedUsername)
	fmt.Printf("[AUTH] Expected password set: %t\n", expectedPassword != "")

	if expectedUsername == "" || expectedPassword == "" {
		fmt.Printf("[AUTH] Admin credentials not configured\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(LoginResponse{Error: "Admin credentials not configured"})
		return
	}

	if req.Username != expectedUsername || req.Password != expectedPassword {
		fmt.Printf("[AUTH] Invalid credentials provided\n")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{Error: "Invalid credentials"})
		return
	}
	
	fmt.Printf("[AUTH] Authentication successful for user: %s\n", req.Username)

	// Generate JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: req.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		fmt.Printf("[AUTH] Token generation failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(LoginResponse{Error: "Could not generate token"})
		return
	}

	fmt.Printf("[AUTH] Token generated successfully\n")
	
	response := LoginResponse{Token: tokenString}
	
	// Test JSON encoding before sending
	testJSON, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("[AUTH] Failed to marshal test JSON: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	fmt.Printf("[AUTH] Response JSON: %s\n", string(testJSON))
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("[AUTH] Failed to encode response: %v\n", err)
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[AUTH] Middleware called for path: %s\n", r.URL.Path)
		
		// Skip auth for login endpoint
		if r.URL.Path == "/login" {
			fmt.Printf("[AUTH] Skipping auth for login endpoint\n")
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		fmt.Printf("[AUTH] Authorization header: %s\n", authHeader)
		if authHeader == "" {
			fmt.Printf("[AUTH] No authorization header found\n")
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			fmt.Printf("[AUTH] Invalid authorization header format: %v\n", bearerToken)
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		fmt.Printf("[AUTH] Attempting to parse token: %s...\n", bearerToken[1][:min(10, len(bearerToken[1]))])
		token, err := jwt.ParseWithClaims(bearerToken[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			fmt.Printf("[AUTH] Token parsing failed: %v\n", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			fmt.Printf("[AUTH] Token valid for user: %s\n", claims.Username)
			// Add username to context for potential use in handlers
			r.Header.Set("X-Username", claims.Username)
			next.ServeHTTP(w, r)
		} else {
			fmt.Printf("[AUTH] Token claims invalid or expired\n")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
		}
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}