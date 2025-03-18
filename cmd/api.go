package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func handleSysCmd(c *chatClient) gin.HandlerFunc {
	return func(g *gin.Context) {
		s, err := io.ReadAll(g.Request.Body)
		if err != nil {
			g.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		resp := c.cmdRegistry.ExecuteCommand(c, string(s))
		if resp == "" {
			g.JSON(http.StatusBadRequest, "bad request")
			return
		}
		g.JSON(http.StatusOK, strings.Split(resp, "\n"))
	}
}

func handleChatRequest(c *chatClient) gin.HandlerFunc {
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

func setupRoutes(c *chatClient) *gin.Engine {
	r := gin.Default()
	r.POST("/syscmd", handleSysCmd(c))
	r.POST("/chat", handleChatRequest(c))

	return r
}
