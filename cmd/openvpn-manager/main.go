package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/bagechashu/openvpn-manager/internal/config"
	"github.com/bagechashu/openvpn-manager/internal/openvpn"
	"github.com/bagechashu/openvpn-manager/web"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	if config.GlobalConfig.Server.AppMode == "production" {
		log.Println("[Info] Run at Production mode")
		gin.SetMode(gin.ReleaseMode)

		// Logging to a file.
		f, _ := os.Create(config.GlobalConfig.LogPath + "openvpn-manager.log")
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

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "web/")
	})

	if config.GlobalConfig.Server.AppMode == "production" {
		r.StaticFS("/web", http.FS(web.Static))
		log.Println("[Info] Webui wse embedd static file")
	} else {
		r.StaticFS("/web", http.Dir("./web/"))
	}

	r.GET("/api/vm/info", openvpn.HandlerVmIpInfo)

	// Group using gin.BasicAuth() middleware
	// gin.Accounts is a shortcut for map[string]string
	r_user := r.Group("/api/user", gin.BasicAuth(config.GlobalAccounts))

	{
		r_user.POST("/add", openvpn.HandlerVpnUserAdd)
		r_user.DELETE("/revoke/:username", openvpn.HandlerVpnUserRevoke)
		r_user.DELETE("/delete/:username", openvpn.HandlerVpnUserRevoke)
		r_user.GET("/show", openvpn.HandlerVpnUserShow)
		r_user.GET("/cert/:filename", openvpn.HandlerVpnUserCert)
	}

	r_ovpn := r.Group("/api/ovpn", gin.BasicAuth(config.GlobalAccounts))
	{
		r_ovpn.GET("/restart", openvpn.HandlerOvpnRestart)
		r_ovpn.GET("/status", openvpn.HandlerOvpnStatus)
	}

	http_str := fmt.Sprintf(":%d", config.GlobalConfig.HttpPort)
	endless.ListenAndServe(http_str, r)
}
