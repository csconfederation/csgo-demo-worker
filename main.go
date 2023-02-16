package main

import (
	"github.com/csconfederation/demoScrape2/pkg/demoscrape2"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	log.Debug(logLevel)
	if !ok {
		logLevel = "info"
	}
	level, logErr := log.ParseLevel(logLevel)
	if logErr != nil {
		level = log.DebugLevel
	}
	log.Debug(level)
	log.SetLevel(level)
	r := gin.Default()
	authUser, hasUser := os.LookupEnv("DEMO_STATS_USER")
	authPass, hasPass := os.LookupEnv("DEMO_STATS_PASSWORD")
	api := r.Group("/api")
	if hasUser && hasPass {
		api = r.Group("/api", gin.BasicAuth(gin.Accounts{
			authUser: authPass,
		}))
	}
	api.POST("/parse", func(c *gin.Context) {
		if c.Request.Body == nil || c.Request.ContentLength == 0 {
			c.JSON(400, "empty request body")
			return
		}
		var game, err = demoscrape2.ProcessDemo(c.Request.Body)
		if err != nil {
			if strings.Contains(err.Error(), "ErrInvalidFileType") {
				c.JSON(400, err.Error())
				return
			}
			c.JSON(500, err.Error())
			return
		}
		c.JSON(200, game)
	})
	api.GET("/parse-remote", func(c *gin.Context) {
		url := c.Query("url")
		authStr := c.Query("auth")
		if url == "" {
			c.JSON(400, "no url specified")
			return
		}
		req, httperr := http.NewRequest("GET", url, nil)
		if httperr != nil {
			c.JSON(500, httperr)
			return
		}
		if authStr != "" {
			req.Header.Set("Authorization", authStr)
		}
		client := &http.Client{
			Timeout: time.Minute * 20,
		}
		resp, respErr := client.Do(req)
		if resp != nil && resp.StatusCode != 200 {
			c.JSON(400, "remote url returned: "+resp.Status)
			return
		}
		if respErr != nil {
			c.JSON(400, respErr)
			return
		}
		if c.Request.Body == nil {
			c.JSON(400, "empty request body")
			return
		}
		var game, err = demoscrape2.ProcessDemo(c.Request.Body)
		if err != nil {
			if strings.Contains(err.Error(), "ErrInvalidFileType") {
				c.JSON(400, err.Error())
				return
			}
			c.JSON(500, err.Error())
			return
		}
		c.JSON(200, game)
	})
	err := r.Run()
	if err != nil {
		println(err)
		return
	}
}
