package openvpn

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandlerVpnUserShow(c *gin.Context) {
	users, err := userShow()
	if err != nil {
		log.Printf("show users error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"response": gin.H{
				"users": nil,
				"msg":   err,
			},
		})
		return
	}
	if len(users) == 0 {
		c.JSON(200, gin.H{
			"status": "success",
			"response": gin.H{
				"users": nil,
				"msg":   "no users",
			},
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"response": gin.H{
			"users": users,
			"count": len(users),
			"msg":   "user query successfully",
		},
	})
}

func userShow() (users []User, err error) {
	filePath := "/etc/openvpn/server/easy-rsa/pki/index.txt"

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	// 检索以 V 开头但不以 server 开头的行数，并展示每行中“=”之后的内容
	for _, line := range lines {
		if strings.HasPrefix(line, "V") && !strings.Contains(line, "=server") {
			parts := strings.Split(line, "=")
			if len(parts) > 1 {
				user := User{
					Username: parts[1],
				}
				users = append(users, user)
			}
		}
	}
	// 倒序添加, 将后创建的用户放在最前面
	for i, j := 0, len(users)-1; i < j; i, j = i+1, j-1 {
		users[i], users[j] = users[j], users[i] //reverse the users
	}
	return users, nil
}
