package api

import (
	"net/http"
	"os"

	"urlx/logging"
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

func requestLogger(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("Incoming request", "method", c.Request.Method, "path", c.Request.URL.Path, "remote", c.ClientIP())
		c.Next()
	}
}

func SetupRouter(s store.Store, shortenerSvc shortener.Shortener, logger *logging.Logger) *gin.Engine {
	r := gin.Default()

	r.Use(requestLogger(logger))

	r.GET("/", func(c *gin.Context) {
		logger.Info("Serving UI page")
		uiPath := "ui/index.html"
		if _, err := os.Stat(uiPath); err == nil {
			c.File(uiPath)
			return
		}
		logger.Error("UI not found")
		c.String(http.StatusNotFound, "UI not found")
	})

	r.GET("/healthz", func(c *gin.Context) {
		logger.Info("Health check endpoint hit")
		c.Status(http.StatusOK)
	})

	r.POST("/shorten", func(c *gin.Context) {
		var req ShortenRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.URL == "" {
			logger.Error("Invalid shorten request", "error", err)
			c.Status(http.StatusBadRequest)
			return
		}
		short := shortenerSvc.Shorten(req.URL)
		s.Set(short, req.URL)
		logger.Info("Shortened URL", "short", short, "long", req.URL)
		c.JSON(http.StatusOK, ShortenResponse{Short: short})
	})

	r.GET("/:short", func(c *gin.Context) {
		short := c.Param("short")
		if long, ok := s.Get(short); ok {
			logger.Info("Redirecting short URL", "short", short, "long", long)
			c.Redirect(http.StatusFound, long)
			return
		}
		logger.Error("Short URL not found", "short", short)
		c.Status(http.StatusNotFound)
	})

	r.DELETE("/delete/:short", func(c *gin.Context) {
		short := c.Param("short")
		s.Delete(short)
		logger.Info("Deleted short URL", "short", short)
		c.Status(http.StatusNoContent)
	})

	return r
}
