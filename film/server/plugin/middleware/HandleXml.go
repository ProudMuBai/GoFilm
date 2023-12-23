package middleware

import (
	"encoding/xml"
	"github.com/gin-gonic/gin"
)

func AddXmlHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.NegotiateFormat(gin.MIMEXML, gin.MIMEJSON) == gin.MIMEXML {
			_, _ = c.Writer.Write([]byte(xml.Header))
		}
		c.Next()
	}
}
