package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
)

var TokenAuth *jwtauth.JWTAuth

func init() {
	godotenv.Load("../.env")
	secretKey := os.Getenv("SECRET_KEY")
	TokenAuth = jwtauth.New("HS256", []byte(secretKey), nil)
}

func CreateToken(userId uint, username string) string {
	_, tokenString, err := TokenAuth.Encode(map[string]interface{}{"userID": userId, "username": username})
	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}

type User struct {
	Username string
	Id       uint
}

func ExtractToken(r *http.Request) *User {
	_, claims, _ := jwtauth.FromContext(r.Context())

	var data = &User{
		Username: "",
		Id:       0,
	}
	if claims["username"] != nil && claims["userID"] != nil {
		data = &User{
			Username: claims["username"].(string),
			Id:       uint(claims["userID"].(float64)),
		}
	}
	return data
}
