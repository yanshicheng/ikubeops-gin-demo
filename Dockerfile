FROM golang:1.22.3 as builder

WORKDIR /app
COPY .  .
RUN go env -w GOPROXY=https://goproxy.cn,direct && env make build

FROM alpine:3.17
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/dist/ikubeops-gin-dem .
COPY config /app/config
CMD ["./ikubeops-gin-dem", "start"]

# docker buildx build --platform linux/amd64  -t harbor.ikubeops.local/public/ingress-manager:latest .