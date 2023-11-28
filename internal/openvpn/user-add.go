package openvpn

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bagechashu/openvpn-manager/internal/config"
	"github.com/gin-gonic/gin"
)

type RequestUserAdd struct {
	Username string `json:"username" binding:"required"`
}

func HandlerVpnUserAdd(c *gin.Context) {
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
	var requestBody RequestUserAdd

	err := c.BindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "failed",
			"response": gin.H{
				"msg": err,
			},
		})
		return
	}

	username, err := userAdd(requestBody.Username)
	if err != nil {
		log.Printf("error: %v", err)
		c.JSON(209, gin.H{
			"status": "failed",
			"response": gin.H{
				"username": username,
				"msg":      fmt.Sprintf("error: %v", err),
			},
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"response": gin.H{
			"username": username,
			"msg":      "user added successfully",
		},
	})
}
func userAdd(unsanitizedUsername string) (client string, err error) {
	unsanitizedUsername = strings.TrimSpace(unsanitizedUsername)
	client = cleanUsername(unsanitizedUsername)

	if !checkIfCertificateExists(client) {
		generateCertificates(client)
		generateClientConfig(client)
		return client, nil
	} else {
		return client, errors.New("Certificate already exists for client: " + client)
	}
}

// cleanUsername: Replace Username nonconforming characters with an underscore _
// 's/[^0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-]/_/g'
func cleanUsername(unsanitizedUsername string) (cleanedClient string) {
	cleanedClient = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= 'A' && r <= 'Z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == '-' || r == '_':
			return r
		default:
			return '_'
		}
	}, unsanitizedUsername)

	return cleanedClient
}

func checkIfCertificateExists(client string) (isExist bool) {
	filePath := "/etc/openvpn/server/easy-rsa/pki/issued/" + client + ".crt"

	// 使用 os.Stat() 函数检查文件是否存在
	if _, err := os.Stat(filePath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		log.Println("check file exist err:", err)
		return false
	}
}

func generateCertificates(client string) {
	// cmd := exec.Command("bash", "-c", `
	// 	cd /etc/openvpn/server/easy-rsa/
	// 	./easyrsa --batch --days=3650 build-client-full "`+client+`" nopass
	// `)

	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// err := cmd.Run()
	// if err != nil {
	// 	fmt.Println("命令执行出错:", err)
	// }

	command := "/etc/openvpn/server/easy-rsa/easyrsa"
	args := []string{"--batch", "--days=3650", "build-client-full", client, "nopass"}

	cmd := execCommand(command, args)
	// init pki RSA env.
	// /etc/openvpn/server/easy-rsa/easyrsa/var.example
	cmd.Env = append(os.Environ(), "EASYRSA_PKI=/etc/openvpn/server/easy-rsa/pki")

	output, err := cmd.Output()
	if err != nil {
		log.Printf("generate cert error: %v", err)
		return
	}

	log.Printf("generate cert success: %s", string(output))
}

func generateClientConfig(client string) error {
	clientCommonPath := "/etc/openvpn/server/client-common.txt"
	caCertPath := "/etc/openvpn/server/easy-rsa/pki/ca.crt"
	clientCertPath := "/etc/openvpn/server/easy-rsa/pki/issued/" + client + ".crt"
	clientKeyPath := "/etc/openvpn/server/easy-rsa/pki/private/" + client + ".key"
	tlsCryptPath := "/etc/openvpn/server/tc.key"

	// 读取文件内容
	clientCommon, err := os.ReadFile(clientCommonPath)
	if err != nil {
		return err
	}

	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return err
	}

	clientCert, err := os.ReadFile(clientCertPath)
	if err != nil {
		return err
	}

	clientKey, err := os.ReadFile(clientKeyPath)
	if err != nil {
		return err
	}

	tlsCrypt, err := os.ReadFile(tlsCryptPath)
	if err != nil {
		return err
	}

	// 写入客户端配置文件
	clientConfig := []byte(fmt.Sprintf(
		"%s\n<ca>\n%s\n</ca>\n<cert>\n%s\n</cert>\n<key>\n%s\n</key>\n<tls-crypt>\n%s\n</tls-crypt>\n",
		clientCommon, caCert, clientCert, clientKey, tlsCrypt,
	))

	// 写入客户端配置文件到 ovpn 目录下
	ovpnDir := "ovpn" // 存放配置文件的目录名
	err = os.MkdirAll(ovpnDir, 0755)
	if err != nil {
		return err
	}

	clientConfigPath := filepath.Join(ovpnDir, client+".ovpn")
	err = os.WriteFile(clientConfigPath, clientConfig, 0644)
	if err != nil {
		return err
	}

	log.Printf("%s added. Configuration available in: %s", client, clientConfigPath)
	return nil
}
