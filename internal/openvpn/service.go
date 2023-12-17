package openvpn

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func HandlerOvpnRestart(c *gin.Context) {
	err := openvpnRestart()
	if err != nil {
		c.JSON(500, gin.H{
			"status": "failed",
			"response": gin.H{
				"msg": fmt.Sprintf("failed to restart openvpn: %s", err),
			},
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"response": gin.H{
			"msg": "openvpn service restart successfuly",
		},
	})
}
func HandlerOvpnStatus(c *gin.Context) {
	status, err := openvpnStatus()
	if err != nil {
		c.JSON(500, gin.H{
			"status": "failed",
			"response": gin.H{
				"msg": fmt.Sprintf("check openvpn service status error: %s", err),
			},
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"response": gin.H{
			"msg": status,
		},
	})
}

func openvpnRestart() error {
	command := "systemctl"
	args := []string{"restart", "openvpn-server@server.service"}

	cmd := execCommand(command, args)

	_, err := cmd.Output()
	if err != nil {
		log.Printf("openvpn service restart err: %s", err)
		return err
	}

	log.Println("openvpn service restart success")
	return nil
}

func openvpnStatus() (result string, err error) {
	command := "systemctl"
	args := []string{"status", "openvpn-server@server.service"}

	cmd := execCommand(command, args)

	output, err := cmd.Output()
	if err != nil {
		log.Printf("openvpn service check status err: %s", err)
		return string(output), err
	}

	log.Printf("Cmd Output: %s", string(output))
	return string(output), err
}
