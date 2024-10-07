package innards

import (
	//"encoding/json"

	"net/http"

	"github.com/gin-gonic/gin"
)

func ApiEndpoints(r *gin.Engine) {
	r.GET("/api/v1/u/:id", func(ctx *gin.Context) {
		var id = ctx.Param("id")
		if !isRequestValidForId(ctx, id) {
			ctx.Status(http.StatusUnauthorized)
			return
		}

		var acc, err = findAccountWithId(id)
		panicIf(err)

		ctx.JSON(http.StatusOK, acc)
	})
	r.GET("/api/v1/u/:id/links", func(ctx *gin.Context) {
		var id = ctx.Param("id")
		if !isRequestValidForId(ctx, id) {
			ctx.Status(http.StatusUnauthorized)
			return
		}

		var acc, err = findAccountWithId(id)
		panicIf(err)

		ctx.JSON(http.StatusOK, acc.Links)
	})

	userAliasV1(r)
}

func userAliasV1(r *gin.Engine) {
	r.GET("/api/v1/u/:id/:alias", func(ctx *gin.Context) {
		var id = ctx.Param("id")
		var alias = ctx.Param("alias")
		if !isRequestValidForId(ctx, id) {
			ctx.Status(http.StatusUnauthorized)
			return
		}

		var acc, err = findAccountWithId(id)
		panicIf(err)

		if _, exists := acc.Links[alias]; !exists {
			ctx.Status(http.StatusNotFound)
			return
		}

		url, err := correspondingLink(id, alias)
		panicIf(err)
		ctx.String(http.StatusOK, url)
	})
	r.POST("/api/v1/u/:id/:alias", func(ctx *gin.Context) {
		var id = ctx.Param("id")
		var alias = ctx.Param("alias")
		if !isRequestValidForId(ctx, id) {
			ctx.Status(http.StatusUnauthorized)
			return
		}

		var acc, err = findAccountWithId(id)
		panicIf(err)

		if _, exists := acc.Links[alias]; exists {
			ctx.Status(http.StatusMethodNotAllowed)
			return
		}

		var reqBody LinkRequest
		err = ctx.ShouldBindBodyWithJSON(&reqBody)
		if err != nil || reqBody.Alias != alias {
			ctx.Status(http.StatusBadRequest)
			return
		}

		url, err := putLinkOnAccountAndPersist(alias, reqBody.Url, acc)
		panicIf(err)
		ctx.String(http.StatusOK, url)
	})
	r.PATCH("/api/v1/u/:id/:alias", func(ctx *gin.Context) {
		var id = ctx.Param("id")
		var alias = ctx.Param("alias")
		if !isRequestValidForId(ctx, id) {
			ctx.Status(http.StatusUnauthorized)
			return
		}

		var acc, err = findAccountWithId(id)
		panicIf(err)

		if _, exists := acc.Links[alias]; !exists {
			ctx.Status(http.StatusNotFound)
			return
		}

		var reqBody LinkRequest
		err = ctx.ShouldBindBodyWithJSON(&reqBody)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		if reqBody.Url != "" {
			reqBody.Url, err = putLinkOnAccountAndPersist(alias, reqBody.Url, acc)
			panicIf(err)
		}
		if reqBody.Alias != "" {
			acc.Links[reqBody.Alias] = acc.Links[alias]
			delete(acc.Links, alias)
		}
		ctx.String(http.StatusOK, reqBody.Url)
	})
	r.DELETE("/api/v1/u/:id/:alias", func(ctx *gin.Context) {
		var id = ctx.Param("id")
		var alias = ctx.Param("alias")
		if !isRequestValidForId(ctx, id) {
			ctx.Status(http.StatusUnauthorized)
			return
		}

		var acc, err = findAccountWithId(id)
		panicIf(err)

		if _, exists := acc.Links[alias]; !exists {
			ctx.Status(http.StatusNotFound)
			return
		}

		delete(acc.Links, alias)
		persistAccount(acc)
		// ctx.Status(http.StatusOK) // is already default
	})
}
