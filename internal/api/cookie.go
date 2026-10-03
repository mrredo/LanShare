package api

import (
	"crypto/rand"
	"encoding/hex"
	"lanshare/config"

	"github.com/gin-gonic/gin"
)

func setOwnerCookie(c *gin.Context) error {
	_, err := c.Cookie("lanshare_owner")
	if err == nil {
		return nil
	}

	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return err
	}

	id := hex.EncodeToString(bytes)

	c.SetCookie(
		"lanshare_owner",
		id,
		60*60*24*365,
		"/",
		"",
		config.CookiesOnHTTPSOnly,
		true,
	)

	return nil
}
