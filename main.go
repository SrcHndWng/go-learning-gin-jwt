package main

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

var secretKey = "TwliyMEXai"

func validate(req *http.Request) (*jwt.Token, error) {
	return request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		b := []byte(secretKey)
		return b, nil
	})
}

func main() {
	r := gin.Default()

	r.GET("/get-token/", func(c *gin.Context) {
		// アルゴリズムの指定
		token := jwt.New(jwt.GetSigningMethod("HS256"))

		// ユーザ、有効期限を設定
		token.Claims = jwt.MapClaims{
			"user": "Guest",
			"exp":  time.Now().Add(time.Hour * 1).Unix(),
		}

		// トークンに対して署名の付与
		tokenString, err := token.SignedString([]byte(secretKey))
		if err == nil {
			c.JSON(200, gin.H{"token": tokenString})
		} else {
			c.JSON(500, gin.H{"message": "Could not generate token"})
		}
	})

	r.GET("/api/private/", func(c *gin.Context) {
		// 署名の検証
		token, err := validate(c.Request)

		if err == nil {
			claims := token.Claims.(jwt.MapClaims)
			msg := fmt.Sprintf("Hello, '%s'!", claims["user"])
			c.JSON(200, gin.H{"message": msg})
		} else {
			c.JSON(401, gin.H{"error": fmt.Sprint(err)})
		}
	})

	r.Run(":8080")
}
