// models.go
package handlers

import "github.com/golang-jwt/jwt/v4"
type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Adm      string `json:"adm"`
    Phone    string `json:"phone"`
    Role     string `json:"role"`
    password string `json:"password"`  // Added password field
    Fee      string `json:"fee"` 
    jwt.RegisteredClaims
}
