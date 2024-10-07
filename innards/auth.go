package innards

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("very secret")

func AuthEndpoints(r *gin.Engine) {
	r.HEAD("/auth/ping", func(ctx *gin.Context) {
		// Empty func means 200 return code with no body
	})
	r.GET("/auth/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})
}

func createJWT(id, password string) (string, error) {
	var t *jwt.Token = jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"iss":         "linkShortener",
		"sub":         "linkShortenerAuth",
		"iat":         time.Now().Unix(),
		"exp":         time.Now().Add(time.Hour).Unix(),
		"jwtId":       id,
		"jwtPassword": password,
	})
	return t.SignedString(secretKey)
}

func validateJWT(tokenString string) (id, password string, success bool) {
	var token, err = jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) { return secretKey, nil })
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var claims, ok = token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		fmt.Println("claims are ok and token is valid")
		success = true
		id, ok = claims["jwtId"].(string)
		if !ok {
			fmt.Println("claims have no id")
			success = false
		}
		password, ok = claims["jwtPassword"].(string)
		if !ok {
			fmt.Println("claims have no pW")
			success = false
		}
	} else {
		fmt.Println("claims:", ok, "token valid:", token.Valid)
	}
	return
}

func isRequestValidForId(ctx *gin.Context, id string) bool {
	var authCookie, err = ctx.Cookie("linkShortenerAuthJwt")
	if err != nil {
		fmt.Println("ERROR getting cookie: ", err.Error())
		return false
	}
	return isJWTValidForId(authCookie, id)
}

func isJWTValidForId(tokenString, id string) bool {
	var jwtId, password, success = validateJWT(tokenString)
	return success && jwtId == id && validateAccount(id, password)
}
