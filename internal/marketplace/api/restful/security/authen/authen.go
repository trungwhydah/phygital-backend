package authen

import (
	"net/http"
	"strings"

	"backend-service/pkg/common/httpresp"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

type AuthenticatorInterface interface {
	Authenticate(c *gin.Context)
}

const (
	userKey = "user"
)

type AuthenticatorDecoder struct {
	fbAuth *auth.Client
}

func NewAuthenticatorDecoder(fbAuth *auth.Client) *AuthenticatorDecoder {
	return &AuthenticatorDecoder{fbAuth: fbAuth}
}

func (d *AuthenticatorDecoder) Decode(c *gin.Context) *auth.Token {
	authHeader := c.Request.Header.Get("Authorization")
	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || !strings.EqualFold(authParts[0], "bearer") {
		httpresp.Error(
			c,
			http.StatusUnauthorized,
			httpresp.ErrKeyAuthenticationInvalidAuthTokenFormat.Error(),
			nil,
		)

		return nil
	}

	rawToken := authParts[1]

	decoded, err := d.fbAuth.VerifyIDToken(c, rawToken)
	if err != nil {
		httpresp.Error(
			c,
			http.StatusUnauthorized,
			httpresp.ErrKeyAuthenticationInvalidSignature.Error(),
			nil,
		)

		return nil
	}

	return decoded
}

func GetUserID(ctx *gin.Context) string {
	userID, ok := ctx.Value(userKey).(string)
	if !ok {
		return ""
	}

	return userID
}
