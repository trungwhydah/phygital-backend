package authen

import (
	"github.com/gin-gonic/gin"
)

type UserAuthenticator struct {
	decoder *AuthenticatorDecoder
}

func NewUserAuthenticator(
	decoder *AuthenticatorDecoder,
) AuthenticatorInterface {
	return &UserAuthenticator{decoder: decoder}
}

func (a *UserAuthenticator) Authenticate(c *gin.Context) {
	tokenData := a.decoder.Decode(c)
	if tokenData == nil || c.IsAborted() {
		return
	}

	c.Set(userKey, tokenData.UID)
	c.Next()
}
