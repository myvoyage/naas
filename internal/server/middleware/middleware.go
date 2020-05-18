package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/nilorg/naas/pkg/tools/key"
	"github.com/nilorg/oauth2"
	"github.com/nilorg/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Header 头处理
func Header() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Server", "naas(https://github.com/nilorg/naas)")
		c.Next()
	}
}

// AuthRequired 身份验证
func AuthRequired(ctx *gin.Context) {
	clientID := ctx.Query("client_id")
	session := sessions.Default(ctx)
	currentAccount := session.Get(key.SessionAccount)
	if currentAccount == nil {
		uri := *ctx.Request.URL
		redirectURI, _ := url.Parse("/oauth2/login")
		redirectURIQuery := url.Values{}
		redirectURIQuery.Set("client_id", clientID)
		redirectURIQuery.Set("login_redirect_uri", uri.String())
		redirectURI.RawQuery = redirectURIQuery.Encode()
		ctx.Redirect(302, redirectURI.String())
		return
	}
	ctx.Next()
}

func parseAuth(auth string) (token string, ok bool) {
	const prefix = "Bearer "
	// Case insensitive prefix match. See Issue 22736.
	if auth == "" || len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	return auth[len(prefix):], true
}

// AuthUserinfoRequired 身份验证
func AuthUserinfoRequired(key interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tok, ok := parseAuth(ctx.GetHeader("Authorization"))
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization Is Empty",
			})
			return
		}
		var (
			idTokenClaims *oauth2.JwtClaims
			err           error
		)
		idTokenClaims, err = oauth2.ParseJwtClaimsToken(tok, key)
		if err != nil {
			logger.Errorf("oauth2.ParseJwtClaimsToken: %s", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": oauth2.ErrAccessDenied.Error(),
			})
			return
		}
		if idTokenClaims == nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": oauth2.ErrAccessDenied.Error(),
			})
			return
		}
		if err = idTokenClaims.Valid(); err != nil {
			logger.Errorf("token valid: %s", err)
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": oauth2.ErrAccessDenied.Error(),
			})
			return
		}
		if !idTokenClaims.VerifyScope("openid", false) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": oauth2.ErrInvalidScope.Error(),
			})
			return
		}
		ctx.Set("idToken", idTokenClaims)
		ctx.Next()
	}
}
