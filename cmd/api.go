package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func handleSysCmd(c *client) gin.HandlerFunc {
	return func(g *gin.Context) {

		s, err := io.ReadAll(g.Request.Body)
		if err != nil {
			g.JSON(http.StatusInternalServerError, err.Error())
		}
		if resp, ok := c.command(string(s)); ok {
			g.JSON(http.StatusOK, strings.Split(resp, "\n"))
			return
		}
		g.JSON(http.StatusBadRequest, []byte("bad request"))
	}
}

func handleChatRequest(c *client) gin.HandlerFunc {
	return func(g *gin.Context) {
		b, err := io.ReadAll(g.Request.Body)
		if err != nil {
			g.JSON(http.StatusBadRequest, "bad request")
			return
		}
		resp, err := c.chatRequest(string(b))
		if err != nil {
			g.JSON(http.StatusInternalServerError, "error handling response")
			return
		}
		g.JSON(http.StatusOK, resp)
	}
}

func setupRoutes(c *client) *gin.Engine {
	r := gin.Default()
	r.POST("/syscmd", handleSysCmd(c))
	r.POST("/chat", handleChatRequest(c))

	return r
}
