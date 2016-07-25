package main

import (
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
)

func main() {

	myJwtMiddleware := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("AllYourBase"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	iris.Get("/ping", pingHandler)
  iris.Get("/auth/login", loginHandler)
	iris.Get("/secured/ping", myJwtMiddleware.Serve, securedPingHandler)

	iris.Listen(":8080")

}

type response struct {
	Text string `json:"text"`
}

//JSON struct that holds generated authorization token
type AuthToken struct {
  Text string `json:"token"`
}

func createToken() (string, error) {
	mySigningKey := []byte("AllYourBase")

	claims := &jwt.StandardClaims{
		Issuer:    "test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func loginHandler(ctx *iris.Context){
  token,_ := createToken()
  authToken := AuthToken{token}
  ctx.JSON(iris.StatusOK, authToken)
}

func pingHandler(ctx *iris.Context) {
	response := response{"All good. You don't need to be authenticated to call this"}
	ctx.JSON(iris.StatusOK, response)
}

func securedPingHandler(ctx *iris.Context) {
	response := response{"All good. You only get this message if you're authenticated"}
	// get the *jwt.Token which contains user information using:
	// user:= myJwtMiddleware.Get(ctx) or context.Get("jwt").(*jwt.Token)
	ctx.JSON(iris.StatusOK, response)
}
