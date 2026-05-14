FROM golang:1.25-alpine AS builder

WORKDIR /app
# 设置 Go 代理和环境变量
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o webhook main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/webhook .

EXPOSE 9099
CMD ["./webhook"]
