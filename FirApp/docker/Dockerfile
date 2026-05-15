FROM node:18-alpine AS frontend-builder
WORKDIR /build/vue
COPY fir_vue/package.json fir_vue/package-lock.json ./
RUN npm ci 2>/dev/null || npm install
COPY fir_vue/ .
RUN npm run build

FROM golang:1.22-alpine AS backend-builder
WORKDIR /build/go

ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
ENV GO111MODULE=on
COPY fir_go/go.mod fir_go/go.sum ./
RUN go mod download
COPY fir_go/ .
RUN CGO_ENABLED=0 go build -o fir .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /usr/src
COPY docker/config.yaml ./config/config.yaml
COPY --from=backend-builder /build/go/fir .
COPY --from=backend-builder /build/go/template/ ./template/
COPY --from=frontend-builder /build/vue/dist ./dist
RUN mkdir -p dist/ios dist/android dist/plist

EXPOSE 80
VOLUME /usr/src/dist/ios
VOLUME /usr/src/dist/android
VOLUME /usr/src/dist/plist

ENTRYPOINT ["./fir"]
