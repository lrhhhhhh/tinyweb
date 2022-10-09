package tinyweb

import (
	"log"
	"net/http"
)

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				c.String(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}
