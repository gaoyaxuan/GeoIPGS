## GeoIPGS
一个基于GeoIP的ip查询服务
需要自行下载GeoIP,命名为:GeoLite2-City.mmdb,映射到容器内/app/GeoLite2-City.mmdb

[GeoIP文件下载教程](http://www.modsecurity.cn/practice/post/15.html)
## 运行教程
```


docker-compose up -d 
请求方式:

curl http://localhost:8080/141.132.14.151

curl http://localhost:8080?ip=141.132.14.151

结果返回:

{
"ip": "141.132.14.151",
"country": "澳大利亚",
"region": "维多利亚州",
"city": "巴拉腊特"
}
```
