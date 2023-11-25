package vm

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ipinfo 结构体用于存储 IPinfo.io 返回的 VPN 信息
type ipinfo struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Org      string `json:"org"`
}

var (
	mutex       sync.Mutex
	ipinfoCache *ipinfo
	expiration  time.Time
)

func fetchVPNInfo() (*ipinfo, error) {
	resp, err := http.Get("http://ipinfo.io")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	vpn := &ipinfo{}
	err = json.Unmarshal(body, vpn)
	if err != nil {
		return nil, err
	}

	return vpn, nil
}

func HandlerVpnInfo(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()

	if ipinfoCache != nil && time.Now().Before(expiration) {
		c.JSON(http.StatusOK, ipinfoCache)
		return
	}

	newVPNInfo, err := fetchVPNInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "get ipinfo.io failed.",
		})
		return
	}

	ipinfoCache = newVPNInfo
	expiration = time.Now().Add(24 * time.Hour)

	c.JSON(http.StatusOK, ipinfoCache)
}
