package main

import (
	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"net/http"
	"strings"
)

var db *geoip2.Reader
var defaultLanguage = "zh-CN"

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	ERROR   = 400
	SUCCESS = 200
)

type IPInfo struct {
	Code    int    `json:"code"`
	IP      string `json:"ip"`
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
}

type QueryParams struct {
	Ip string `json:"ip" form:"ip"`
}

// 解析 Accept-Language 字符串，返回首选语言
func parseAcceptLanguage(al string) string {
	if al == "" {
		return defaultLanguage
	}
	// 将语言标签按优先级排序（忽略权重）
	languages := strings.Split(al, ",")
	if len(languages) > 0 {
		// 提取第一个语言标签并返回
		return strings.SplitN(languages[0], ";", 2)[0]
	}
	return defaultLanguage
}

func getIPInfo(ip net.IP, language string) (IPInfo, error) {
	record, err := db.City(ip)
	if err != nil {
		return IPInfo{}, err
	}
	ipInfo := IPInfo{
		Code:    SUCCESS,
		IP:      ip.String(),
		Country: record.Country.Names[language],
		City:    record.City.Names[language],
	}

	if len(record.Subdivisions) > 0 {
		ipInfo.Region = record.Subdivisions[0].Names[language]
	}

	return ipInfo, nil
}

func handleIP(ipStr string, c *gin.Context) {
	// 检查ip格式是否正确
	ip := net.ParseIP(ipStr)
	if ip == nil {
		c.JSON(http.StatusOK, Response{Code: ERROR, Message: "Invalid IP address"})
		return
	}
	acceptLanguage := c.GetHeader("Accept-Language")
	ipInfo, err := getIPInfo(ip, parseAcceptLanguage(acceptLanguage))
	if err != nil {
		c.JSON(http.StatusOK, Response{Code: ERROR, Message: "Request failed"})
		return
	}
	c.JSON(http.StatusOK, ipInfo)
}

func ipHandler(c *gin.Context) {
	var query QueryParams
	err := c.ShouldBind(&query)
	if err != nil {
		c.JSON(http.StatusOK, Response{Code: ERROR, Message: "Request failed"})
		return
	}
	if query.Ip == "" {
		query.Ip = c.ClientIP()
	}
	handleIP(query.Ip, c)
}

func ipPathHandler(c *gin.Context) {
	ipStr := c.Param("ip")
	handleIP(ipStr, c)
}

type queryOnlyIp struct {
	Json bool `json:"json" form:"json"`
}

func onlyIp(c *gin.Context) {
	var query queryOnlyIp
	err := c.ShouldBind(&query)
	if err != nil {
		c.JSON(http.StatusOK, Response{Code: ERROR, Message: "Request failed"})
		return
	}
	if query.Json {
		c.JSON(http.StatusOK, gin.H{"code": SUCCESS, "ip": c.ClientIP()})
		return
	}
	c.String(http.StatusOK, c.ClientIP())
}

func main() {
	// 直接通过代码设置 Gin 为生产模式
	gin.SetMode(gin.ReleaseMode)

	db1, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		panic(err)
	}
	db = db1
	defer db.Close()

	r := gin.Default()
	// 只获取ip信息 默认返回字符串, 加上参数 json=true 会返回json格式 { "ip":"8.8.8.8" }
	r.GET("/ip", onlyIp)
	r.POST("/ip", onlyIp)
	r.GET("/", ipHandler)
	r.POST("/", ipHandler)
	r.GET("/:ip", ipPathHandler)
	r.POST("/:ip", ipPathHandler)
	err = r.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("--------服务启动失败: %v--------\n", err)
		return
	}
}
