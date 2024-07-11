package main

import (
	"log"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
)

type IPInfo struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
}

func getClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	xff := c.GetHeader("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ", ")
		if len(ips) > 0 {
			ip = ips[0]
		}
	}
	return ip
}

func getIPInfo(db *geoip2.Reader, ip string) (IPInfo, error) {
	parsedIP := net.ParseIP(ip)
	record, err := db.City(parsedIP)
	if err != nil {
		return IPInfo{}, err
	}

	ipInfo := IPInfo{
		IP:      ip,
		Country: record.Country.Names["zh-CN"],
		City:    record.City.Names["zh-CN"],
	}

	if len(record.Subdivisions) > 0 {
		ipInfo.Region = record.Subdivisions[0].Names["zh-CN"]
	}

	return ipInfo, nil
}

func handleIP(db *geoip2.Reader, ip string, c *gin.Context) {
	if net.ParseIP(ip) == nil {
		c.JSON(400, gin.H{"error": "无效的IP地址"})
		return
	}
	ipInfo, err := getIPInfo(db, ip)
	if err != nil {
		c.JSON(500, gin.H{"error": "获取IP信息失败"})
		return
	}
	c.JSON(200, ipInfo)
}

func ipHandler(db *geoip2.Reader) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.Query("ip")
		if ip == "" {
			ip = getClientIP(c)
		}

		handleIP(db, ip, c)
	}
}

func ipPathHandler(db *geoip2.Reader) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.Param("ip")
		handleIP(db, ip, c)
	}
}

func main() {
	// 直接通过代码设置 Gin 为生产模式
	gin.SetMode(gin.ReleaseMode)

	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := gin.Default()
	r.GET("/", ipHandler(db))
	r.GET("/:ip", ipPathHandler(db))
	err = r.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("服务启动失败: %v\n", err)
		return
	}
}
