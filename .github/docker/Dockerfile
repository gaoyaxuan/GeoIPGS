# ----------------- 构建阶段 -----------------
FROM --platform=$BUILDPLATFORM golang:alpine AS build
WORKDIR /src

# 1. 下载依赖 (利用 Docker 缓存机制)
# 只有 go.mod 变动时才重新下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 2. 复制源码
COPY . .

# 3. 编译优化
# -ldflags="-s -w" : 去除符号表和调试信息 (减小体积)
# -trimpath        : 移除文件系统路径信息 (安全性)
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 \
    go build -ldflags="-s -w" -trimpath -o /app/GeoIPGS ./main.go

# 4. 准备只有必要文件的 "根文件系统"
# 创建一个非 root 用户
RUN echo "appuser:x:10001:10001::/app:/sbin/nologin" > /app/passwd && \
    echo "appuser:x:10001:" > /app/group

# 安装时区数据和CA证书 (从 alpine 提取)
RUN apk add --no-cache tzdata ca-certificates && \
    cp /usr/share/zoneinfo/Asia/Shanghai /app/localtime && \
    echo "Asia/Shanghai" > /app/timezone

# ----------------- 运行阶段 (Scratch) -----------------
FROM scratch

WORKDIR /app

# 1. 从构建层只复制必要文件
COPY --from=build /app/passwd /etc/passwd
COPY --from=build /app/group /etc/group
COPY --from=build /app/localtime /etc/localtime
COPY --from=build /app/timezone /etc/timezone
# 如果程序需要访问 https 外部链接，取消下面这行的注释
# COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 2. 复制二进制文件
COPY --from=build /app/GeoIPGS /app/GeoIPGS

# 3. 复制数据文件 (如果有默认文件想打进镜像)
# COPY ip2region.xdb /app/ip2region.xdb

# 4. 设置环境变量
ENV TZ=Asia/Shanghai

# 5. 指定非 root 用户运行
USER appuser:appuser

EXPOSE 8080

CMD [ "/app/GeoIPGS" ]
