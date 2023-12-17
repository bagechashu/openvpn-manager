package openvpn

import (
	"encoding/json"
	"errors"
	"io"
	"log"
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

func fetchIpInfo() (*ipinfo, error) {
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
	if err = json.Unmarshal(body, vpn); err != nil {
		return nil, err
	}
	log.Printf("vpn unmarshal result: %s, \nfetch vpn result: \n%s", vpn, body)

	// 	{
	//   "status": 429,
	//   "error": {
	//     "title": "Rate limit exceeded",
	//     "message": "You've hit the daily limit for the unauthenticated API.  Create an API access token by signing up to get 50k req/month."
	//   }
	//  }
	if vpn.IP == "" {
		return nil, errors.New("ipinfo.io rate limit exceeded")
	}

	return vpn, nil
}

func HandlerVmIpInfo(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()

	if ipinfoCache != nil && time.Now().Before(expiration) {
		c.JSON(http.StatusOK, ipinfoCache)
		return
	}

	newIpInfo, err := fetchIpInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "get ipinfo.io failed.",
		})
		return
	}

	ipinfoCache = newIpInfo
	expiration = time.Now().Add(24 * time.Hour)

	c.JSON(http.StatusOK, ipinfoCache)
}
