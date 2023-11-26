package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/bagechashu/openvpn-manager/internal/config"
	"github.com/bagechashu/openvpn-manager/internal/openvpn"
	"github.com/bagechashu/openvpn-manager/internal/vm"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"

	"gopkg.in/ini.v1"
)

func main() {
	configPath := flag.String("c", "openvpn-manager.ini", "config file path")
	flag.Parse()

	cfg, err := ini.Load(*configPath)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	if cfg.Section("").Key("app_mode").In("development", []string{"development", "production"}) == "production" {
		log.Println("[Info] Run at Production mode")
		gin.SetMode(gin.ReleaseMode)

		// Logging to a file.
		f, _ := os.Create("openvpn-manager.log")
		gin.DefaultWriter = io.MultiWriter(f)
	} else {
		log.Println("[Info] Run at Development mode, change configuare 'app_mode = production' to run at Production mode")
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api/vpn/info", vm.HandlerVpnInfo)

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "web/")
	})

	r.StaticFS("/web", http.Dir("./web"))

	// Group using gin.BasicAuth() middleware
	// gin.Accounts is a shortcut for map[string]string
	r_user := r.Group("/api/user", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	{
		// /admin/secrets endpoint
		// hit "localhost:8080/admin/secrets
		r_user.GET("/secrets", func(c *gin.Context) {
			// get user, it was set by the BasicAuth middleware
			user := c.MustGet(gin.AuthUserKey).(string)
			if secret, ok := config.Secrets[user]; ok {
				c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
			} else {
				c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
			}
		})
		r_user.POST("/add", openvpn.HandlerVpnUserAdd)
		r_user.GET("/download", func(c *gin.Context) {
			// get user, it was set by the BasicAuth middleware
			user := c.MustGet(gin.AuthUserKey).(string)
			if secret, ok := config.Secrets[user]; ok {
				c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
			} else {
				c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
			}
		})
		r_user.DELETE("/revoke", func(c *gin.Context) {
			// get user, it was set by the BasicAuth middleware
			user := c.MustGet(gin.AuthUserKey).(string)
			if secret, ok := config.Secrets[user]; ok {
				c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
			} else {
				c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
			}
		})
	}

	http_str := fmt.Sprintf(":%d", cfg.Section("server").Key("http_port").MustInt(8080))
	endless.ListenAndServe(http_str, r)
}
