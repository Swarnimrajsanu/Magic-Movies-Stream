package utils

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/Swarnimrajsanu/MagicMoviesStream/Server/MagicMoviesServer/database"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Role      string
	UserId    string
	jwt.RegisteredClaims
}

var SECRET_KEY = os.Getenv("SECRET_KEY")
var SECRET_REFRESH_KEY = os.Getenv("SECRET_REFRESH_KEY")

// ========================= JWT GENERATION =========================

func GenerateAllTokens(email, firstName, lastName, role, userId string) (string, string, error) {

	// ACCESS TOKEN
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	// REFRESH TOKEN
	refreshClaims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(SECRET_REFRESH_KEY))
	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}

// ========================= UPDATE TOKENS =========================

func UpdateAllTokens(userId, token, refreshToken string, client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updateData := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"update_at":     time.Now(),
		},
	}

	userCollection := database.OpenCollection("users", client)

	_, err := userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateData)
	return err
}

// ========================= GET TOKEN FROM COOKIE =========================

func GetAccessToken(c *gin.Context) (string, error) {
	tokenString, err := c.Cookie("access_token")
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ========================= VALIDATE ACCESS TOKEN =========================

func ValidateToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}

	if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
		return nil, errors.New("invalid signing algorithm")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}

// ========================= GET VALUES FROM CONTEXT (MIDDLEWARE SET) =========================

func GetUserIdFromContext(c *gin.Context) (string, error) {
	userId, exists := c.Get("userId")
	if !exists {
		return "", errors.New("userId does not exist in context")
	}
	id, ok := userId.(string)
	if !ok {
		return "", errors.New("unable to cast userId to string")
	}
	return id, nil
}

func GetRoleFromContext(c *gin.Context) (string, error) {
	role, exists := c.Get("role")
	if !exists {
		return "", errors.New("role does not exist in context")
	}
	r, ok := role.(string)
	if !ok {
		return "", errors.New("unable to cast role")
	}
	return r, nil
}

// ========================= VALIDATE REFRESH TOKEN =========================

func ValidateRefreshToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_REFRESH_KEY), nil
	})
	if err != nil {
		return nil, err
	}

	if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
		return nil, errors.New("invalid signing algorithm")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("refresh token has expired")
	}

	return claims, nil
}
