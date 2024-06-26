package openvpn

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bagechashu/openvpn-manager/internal/config"
	"github.com/gin-gonic/gin"
)

func HandlerVpnUserCert(c *gin.Context) {
	filename := c.Param("filename")

	file := filepath.Join(config.GlobalConfig.Ovpn.CertPath, cleanUsername(filename))

	// check does file exist
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"response": gin.H{
				"msg": fmt.Sprintf("%s is not exist", filename),
			},
		})
		return
	}

	// 提供文件下载
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(file)
}
