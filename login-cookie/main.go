package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}
var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Allow requests from your frontend domain
        c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
        
        // Allow credentials (cookies) to be sent with requests
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        
        // Allow specific HTTP methods
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        
        // Allow specific headers
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
        
        // Handle preflight requests (OPTIONS method)
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusOK)
            return
        }
        
        // Pass control to the next middleware or handler
        c.Next()
    }
}

func main() {
	dsn := "hadi:admin123@tcp(127.0.0.1:3306)/cookie_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	db.AutoMigrate(&User{})

	r := gin.Default()
	r.Use(CORSMiddleware())
	// Public routes
	r.POST("/register", Register)
	r.POST("/login", Login)
	r.GET("/logout", Logout)

	// Protected routes
	protected := r.Group("/protected")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/profile", Profile)
	}

	r.Run(":8080")
}

func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password before storing it (you should use bcrypt or similar in production)
	hashedPassword := user.Password // In production, hash this password
	user.Password = hashedPassword

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundUser User
	if err := db.Where("username = ?", user.Username).First(&foundUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password (in production, compare hashed passwords)
	if foundUser.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: foundUser.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Set token in cookie
	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
// Gunakan struct UserResponse untuk response tanpa password
	userResponse := UserResponse{
		ID:       foundUser.ID,
		Username: foundUser.Username,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
		"user":    userResponse,
	})

}

func Logout(c *gin.Context) {
	// Clear the token cookie
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// AuthMiddleware checks for the JWT token in the cookie and validates it
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the cookie
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
			c.Abort()
			return
		}

		// Parse and validate the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Attach the claims to the context for further use
		c.Set("username", claims.Username)
		c.Next()
	}
}

// Profile is a protected route that returns user information
func Profile(c *gin.Context) {
  // Ambil username dari context (disetel di AuthMiddleware)
  username, exists := c.Get("username")
  if !exists {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
    return
  }

  // Kirim respons berisi username
  c.JSON(http.StatusOK, gin.H{"message": "Welcome to your profile", "username": username})
}


