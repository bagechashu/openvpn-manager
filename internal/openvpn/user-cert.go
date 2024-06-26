package openvpn

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bagechashu/openvpn-manager/internal/config"
	"github.com/gin-gonic/gin"
)

func HandlerVpnUserCert(c *gin.Context) {
	filename := c.Param("filename")

	filename, err := formatFilename(filename)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"response": gin.H{
				"msg": fmt.Sprintf("%s", err),
			},
		})
		return
	}
	file := filepath.Join(config.GlobalConfig.Ovpn.CertPath, filename)
	log.Printf("[Info] Download file: %s", file)
	// check does file exist
	_, err = os.Stat(file)
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

func formatFilename(filename string) (string, error) {
	fnParts := strings.Split(filename, ".")
	if len(fnParts) >= 2 && fnParts[len(fnParts)-1] == "ovpn" {
		// Reconstruct the filename without the last part (i.e., everything after the last '.')
		fn := strings.Join(fnParts[:len(fnParts)-1], ".")
		fn = cleanUsername(fn)
		return fmt.Sprintf("%s.%s", fn, fnParts[len(fnParts)-1]), nil
	} else {
		return "", fmt.Errorf("filename have some issue")
	}
}
