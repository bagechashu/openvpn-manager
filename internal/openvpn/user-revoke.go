package openvpn

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bagechashu/openvpn-manager/internal/config"
	"github.com/gin-gonic/gin"
)

var easyrsaCmd = "/etc/openvpn/server/easy-rsa/easyrsa"
var easyrsaEnv = "EASYRSA_PKI=/etc/openvpn/server/easy-rsa/pki"

// var groupName = vm.GetNobodyGroup()

func HandlerVpnUserRevoke(c *gin.Context) {
	// get user, it was set by the BasicAuth middleware
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := config.Secrets[user]; !ok {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "failed",
			"response": gin.H{
				"msg": fmt.Sprintf("user %s didn't have permission", user),
			},
		})
		return
	}
	username := c.Param("username")
	if username == "" {
		c.JSON(400, gin.H{
			"status": "failed",
			"response": gin.H{
				"msg": "username is required",
			},
		})
		return
	}
	err := userRevoke(username)
	if err != nil {
		c.JSON(209, gin.H{
			"status": "failed",
			"response": gin.H{
				"username": username,
				"msg":      fmt.Sprintf("failed to revoke user %s", username),
			},
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"response": gin.H{
			"username": username,
			"msg":      "user revoke successfully",
		},
	})
}

func userRevoke(client string) (err error) {
	if err = revokeCertificates(client); err != nil {
		return err
	}
	if err = genCrlCertificates(client); err != nil {
		return err
	}
	if err = resetCrlFile(); err != nil {
		return err
	}

	log.Printf("%s revoked!\n", client)
	return nil
}

// Revoke specified client using ./easyrsa --batch revoke "$client" command
func revokeCertificates(client string) (err error) {
	args := []string{"--batch", "revoke", client}
	cmd := execCommand(easyrsaCmd, args)
	cmd.Env = append(os.Environ(), easyrsaEnv)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("revoke cert error: %v", err)
		return err
	}
	log.Printf("revoke cert success: %s", string(output))
	return nil
}

// Generate a new CRL with a validity of 3650 days using ./easyrsa --batch --days=3650 gen-crl command
func genCrlCertificates(client string) (err error) {
	args := []string{"--batch", "--days=3650", "gen-crl"}
	cmd := execCommand(easyrsaCmd, args)
	cmd.Env = append(os.Environ(), easyrsaEnv)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("gen CRL error: %v", err)
		return err
	}
	log.Printf("gen CRL success: %s", string(output))
	return nil
}

// rm -f /etc/openvpn/server/crl.pem
// cp /etc/openvpn/server/easy-rsa/pki/crl.pem /etc/openvpn/server/crl.pem
// # CRL is read with each client connection, when OpenVPN is dropped to nobody
// chown nobody:"$group_name" /etc/openvpn/server/crl.pem
func resetCrlFile() (err error) {
	sourcePath := "/etc/openvpn/server/easy-rsa/pki/crl.pem"
	destinationPath := "/etc/openvpn/server/crl.pem"

	// Read the content of the source file
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		log.Println("Error reading source file:", err)
		return err
	}

	// Write the content to the destination file
	err = os.WriteFile(destinationPath, content, 0644)
	if err != nil {
		log.Println("Error writing to destination file:", err)
		return err
	}

	log.Println("File content copied successfully!")
	return nil
}
