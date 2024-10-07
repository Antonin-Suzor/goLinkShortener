package innards

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const baseDir = "exposableFiles/"

func UrlMapping(r *gin.Engine) {
	r.StaticFile("/", baseDir+"publicPages/index.html")
	linkShortenerRedirect(r)

	scripts(r)
	media(r)

	r.LoadHTMLGlob(baseDir + "privatePages/*")

	signup(r)
	login(r)
	myaccount(r)
}

func linkShortenerRedirect(r *gin.Engine) {
	r.Any("/u/:id/:alias", func(ctx *gin.Context) {
		var id = ctx.Param("id")
		var alias = ctx.Param("alias")
		var corresponding, err = correspondingLink(id, alias)
		if err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}
		if !strings.Contains(corresponding, "://") {
			corresponding = "https://" + corresponding
		}
		ctx.Redirect(http.StatusTemporaryRedirect, corresponding)
	})
}

func scripts(r *gin.Engine) {
	r.Static("scripts", baseDir+"scripts")
}

func media(r *gin.Engine) {
	r.Static("media", baseDir+"media")
}

func signup(r *gin.Engine) {
	r.GET("/signup", func(ctx *gin.Context) {
		ctx.File(baseDir + "publicPages/signup.html")
	})
	r.POST("/signup/post", func(ctx *gin.Context) {
		var reqBody SignupRequest
		var err = ctx.ShouldBind(&reqBody)
		panicIf(err)

		if existsAccountWithId(reqBody.Id) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"reason": "id"})
			return
		}
		if existsAccountWithEmail(reqBody.Email) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"reason": "email"})
			return
		}

		var acc = Account{Id: reqBody.Id, Email: reqBody.Email, Password: reqBody.Password}
		err = persistAccount(acc)
		panicIf(err)

		jwt, err := createJWT(reqBody.Id, reqBody.Password)
		panicIf(err)

		ctx.JSON(http.StatusCreated, gin.H{"jwtCookie": "linkShortenerAuthJwt=" + jwt + "; Path=/"})
	})
}

func login(r *gin.Engine) {
	r.GET("/login", func(ctx *gin.Context) {
		ctx.File(baseDir + "publicPages/login.html")
	})
	r.POST("/login/post", func(ctx *gin.Context) {
		var reqBody LoginRequest
		var err = ctx.ShouldBind(&reqBody)
		panicIf(err)

		if !validateAccount(reqBody.Id, reqBody.Password) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"reason": "password"})
			return
		}

		jwt, err := createJWT(reqBody.Id, reqBody.Password)
		panicIf(err)

		ctx.JSON(http.StatusOK, gin.H{"jwtCookie": "linkShortenerAuthJwt=" + jwt + "; Path=/"})
	})
}

func myaccount(r *gin.Engine) {
	r.GET("/myaccount", func(ctx *gin.Context) {
		var authCookie, err = ctx.Cookie("linkShortenerAuthJwt")
		if err != nil {
			fmt.Println("ERROR getting cookie: ", err.Error())
			ctx.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}
		var id, password, success = validateJWT(authCookie)
		if !success || !validateAccount(id, password) {
			fmt.Println("JWT not validated, success: ", success, "account Valid: ", validateAccount(id, password))
			ctx.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		acc, err := findAccountWithId(id)
		panicIf(err)
		ctx.HTML(http.StatusOK, "myaccount.html", gin.H{"acc": acc})
	})
}
