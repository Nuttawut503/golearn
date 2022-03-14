package main

import (
	"context"
	"errors"
	"fmt"
	"gojwt/auth"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func saveDataInRedis(rdb *redis.Client, key, value string, expire int64) error {
	return rdb.Set(context.Background(), key, value, time.Unix(expire, 0).Sub(time.Now())).Err()
}

func AuthMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		splits := strings.Split(ctx.Request.Header.Get("Authorization"), " ")
		if len(splits) != 2 {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("No auth token"))
			return
		}
		token := splits[1]
		claims, err := auth.ExtractToken(token, false)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("Invalid token"))
			return
		}
		uuid, ok := claims["access_uuid"].(string)
		if !ok {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("Token info is missing"))
			return
		}
		username, err := rdb.Get(context.Background(), uuid).Result()
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("Token is expired"))
			return
		}
		ctx.Set("accessuuid", uuid)
		ctx.Set("username", username)
	}
}

type (
	LoginForm struct {
		Username string
	}
	RefreshRequestForm struct {
		RefreshToken string
	}
	LogoutRequestForm struct {
		RefreshToken string
	}
)

func main() {
	r := gin.Default()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	r.POST("/login", func(ctx *gin.Context) {
		var form LoginForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.String(http.StatusUnprocessableEntity, "Wrong format")
			return
		}

		tokens, err := auth.GeneratePairToken(form.Username)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Something went wrong")
			return
		}

		if err := saveDataInRedis(rdb, tokens[0].UUID, form.Username, tokens[0].Expire); err != nil {
			ctx.String(http.StatusInternalServerError, "Something went wrong")
			return
		}
		if err := saveDataInRedis(rdb, tokens[1].UUID, form.Username, tokens[1].Expire); err != nil {
			ctx.String(http.StatusInternalServerError, "Something went wrong")
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"AccessToken":  tokens[0].Token,
			"RefreshToken": tokens[1].Token,
		})
	})
	authRoute := r.Group("/", AuthMiddleware(rdb))
	authRoute.GET("/resource", func(ctx *gin.Context) {
		var username string
		if v, ok := ctx.Get("username"); !ok {
			ctx.String(http.StatusForbidden, "Something went wrong")
			return
		} else if username, ok = v.(string); !ok {
			ctx.String(http.StatusForbidden, "Some info is missing")
			return
		}
		ctx.String(http.StatusOK, fmt.Sprintf("Hello %s!", username))
	})
	r.POST("/refresh", func(ctx *gin.Context) {
		var form RefreshRequestForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.String(http.StatusUnprocessableEntity, "Wrong format")
			return
		}
		claims, err := auth.ExtractToken(form.RefreshToken, true)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("Invalid token"))
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("Token info is missing"))
			return
		}
		accessToken, err := auth.GenerateToken(username, false)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Something went wrong")
			return
		}
		if err := saveDataInRedis(rdb, accessToken.UUID, username, accessToken.Expire); err != nil {
			ctx.String(http.StatusInternalServerError, "Something went wrong")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"AccessToken": accessToken.Token,
		})
	})
	authRoute.POST("/logout", func(ctx *gin.Context) {
		var form LogoutRequestForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.String(http.StatusUnprocessableEntity, "Wrong format")
			return
		}
		var accessuuid string
		if v, ok := ctx.Get("accessuuid"); !ok {
			ctx.String(http.StatusForbidden, "Something went wrong")
			return
		} else if accessuuid, ok = v.(string); !ok {
			ctx.String(http.StatusForbidden, "Some info is missing")
			return
		}
		refreshClaims, err := auth.ExtractToken(form.RefreshToken, true)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("Invalid token"))
			return
		}
		refreshuuid, ok := refreshClaims["refresh_uuid"].(string)
		if !ok {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("Token info is missing"))
			return
		}
		rdb.Del(context.Background(), accessuuid)
		rdb.Del(context.Background(), refreshuuid)
		ctx.String(http.StatusOK, "Logged out successfully")
	})
	r.Run()
}
