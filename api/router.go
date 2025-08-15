package api

import (
	"net/http"

	"urlx/shortener"
	"urlx/store"

	"github.com/gin-gonic/gin"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Short string `json:"short"`
}

func SetupRouter(s store.Store, shortenerSvc shortener.Shortener) *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.POST("/shorten", func(c *gin.Context) {
		var req ShortenRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.URL == "" {
			c.Status(http.StatusBadRequest)
			return
		}
		short := shortenerSvc.Shorten(req.URL)
		s.Set(short, req.URL)
		c.JSON(http.StatusOK, ShortenResponse{Short: short})
	})

	r.GET("/:short", func(c *gin.Context) {
		short := c.Param("short")
		if long, ok := s.Get(short); ok {
			c.Redirect(http.StatusFound, long)
			return
		}
		c.Status(http.StatusNotFound)
	})

	r.DELETE("/delete/:short", func(c *gin.Context) {
		short := c.Param("short")
		s.Delete(short)
		c.Status(http.StatusNoContent)
	})

	return r
}
