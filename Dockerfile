# 多阶段构建: 编译 Go 后端 (内嵌 SQLite, 纯 Go 无需 cgo)
FROM golang:1.24 AS builder
WORKDIR /src
COPY server/go.mod server/go.sum ./server/
RUN cd server && go mod download
COPY server/ ./server/
RUN cd server && CGO_ENABLED=0 go build -o /out/yfsc-server .

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /out/yfsc-server /app/yfsc-server
COPY web/ /app/web/
EXPOSE 8080
ENV ADDR=:8080 WEB_DIR=/app/web DB_PATH=/app/data/yfsc.db
VOLUME ["/app/data"]
ENTRYPOINT ["/app/yfsc-server"]
