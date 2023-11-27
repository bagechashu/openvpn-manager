package openvpn

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func HandlerVpnUserCert(c *gin.Context) {
	filename := c.Param("filename")
	file := filepath.Join("./ovpn", filename)

	// 检查文件是否存在
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
