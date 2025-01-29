package Network

import (
	"encoding/json"
	"graduationDesign/back_end/Initialize/Parameter"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// jwt密钥字符串
var jwtKey = []byte("dadawre")

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

// Claims JWT 的 Claims 结构
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 登录处理函数
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Mock 验证（替换为实际验证逻辑）
	if req.Username == Parameter.Account && req.Password == Parameter.Password {
		// 生成 JWT Token
		expirationTime := time.Now().Add(8 * time.Hour) // 有效期 8 小时
		claims := &Claims{
			Username: req.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		res := LoginResponse{Success: true, Message: "Login successful", Token: tokenString}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	} else {
		res := LoginResponse{Success: false, Message: "Username or password is incorrect"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}

// 验证 JWT Token 的中间件
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 验证成功，继续处理请求
		next(w, r)
	}
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a protected route"))
}
