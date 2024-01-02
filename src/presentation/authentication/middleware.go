package authentication

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Middleware(tokenAudience string) gin.HandlerFunc {
	authZeroSettings := AuthZeroSettings{
		JwksURI:  "https://bexs.auth0.com/.well-known/jwks.json",
		Audience: []string{tokenAudience},
		Issuer:   "https://bexs.auth0.com/",
	}

	auth0Provider := NewAuthZeroProvider(authZeroSettings)
	Setup(auth0Provider)

	return func(c *gin.Context) {
		//TODO check in cache
		err := Validate(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Next()
	}
}
