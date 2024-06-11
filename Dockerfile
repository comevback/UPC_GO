# 第一个阶段：使用 Go 官方镜像构建 Go 可执行文件
FROM golang:1.22-alpine as builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载所有依赖项
RUN go mod download

# 复制源代码
COPY . .

# 构建 Go 应用程序
RUN go build -o upc-go .

# 第二个阶段：使用最小化的基础镜像运行应用程序
FROM alpine:latest

# 设置工作目录
WORKDIR /usr/src/app

# 安装必要的软件包并清理缓存
RUN apk update && apk add --no-cache \
    curl \
    zsh \
    git \
    && sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" \
    && (curl -sSL "https://github.com/buildpacks/pack/releases/download/v0.32.1/pack-v0.32.1-linux.tgz" | tar -C /usr/local/bin/ --no-same-owner -xzv pack) \
    && apk del git \
    && rm -rf /var/cache/apk/* /tmp/* /var/tmp/* /usr/share/man /usr/share/doc /usr/share/licenses

# 从构建阶段复制预构建的二进制文件和其他需要的文件
COPY --from=builder /app/upc-go .
COPY --from=builder /app/views ./views
COPY --from=builder /app/public ./public

# 暴露端口 4000
EXPOSE 4000

# 设置 Zsh 为默认 shell
SHELL ["/bin/zsh", "-c"]

# 运行可执行文件
CMD ["./upc-go"]
