package authen

import (
	"github.com/gin-gonic/gin"
)

type AdminAuthenticator struct {
	decoder *AuthenticatorDecoder
}

func NewAdminAuthenticator(
	decoder *AuthenticatorDecoder,
) AuthenticatorInterface {
	return &AdminAuthenticator{decoder: decoder}
}

func (a *AdminAuthenticator) Authenticate(c *gin.Context) {
	tokenData := a.decoder.Decode(c)
	if tokenData == nil || c.IsAborted() {
		return
	}

	// TODO: verify admin

	c.Set(userKey, tokenData.UID)
	c.Next()
}
