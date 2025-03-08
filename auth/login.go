package auth

import (
	"errors"
	"golist/database"
	"golist/models"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

// Define secret JWT key (asap put a real key in .env)
var jwtKey = []byte("ma_clé_secrète")

// Data of the token
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 1 week token
		},
	}

	// Create new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func CheckJWT(c *fiber.Ctx) (string, error) {

	tokenString := c.Cookies("token")
	if tokenString == "" {
		return "", errors.New("Token missing")
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("token invalid ou expired")
	}

	//
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Error when parsing claims")
	}

	// Check if the token is expired
	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return "", errors.New("Token expired")
	}

	return claims["username"].(string), nil
}

func LoginUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).SendString("Parsing error")
	}

	var dbUser models.User
	if err := database.DB.Where("username = ?", user.Username).First(&dbUser).Error; err != nil {
		return c.Status(401).SendString("Username or password incorrect")
	}

	err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return c.Status(401).SendString("Username or password incorrect")
	}

	// Generate JWT when logged in
	token, err := GenerateJWT(dbUser.Username)
	if err != nil {
		return c.Status(500).SendString("Error generating token")
	}

	// Set token in cookies
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Secure:   true,
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"token": token,
	})
}
