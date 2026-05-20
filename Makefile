.PHONY: run build test smoke clean fresh

# 启动平台 (API + 静态前端) → http://localhost:8080
run:
	cd server && go run . -addr :8080 -web ../web -db yfsc.db

# 编译单一可执行文件到 bin/yfsc-server (从磁盘 web/ 提供前端)
build:
	cd server && go build -o ../bin/yfsc-server .

# 发布构建: 前端内嵌, 生成各平台自包含单文件到 dist/
release:
	rm -rf server/webembed && cp -r web server/webembed
	mkdir -p dist
	cd server && CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags embed -trimpath -ldflags "-s -w" -o ../dist/yfsc-server-linux-amd64 .
	cd server && CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -tags embed -trimpath -ldflags "-s -w" -o ../dist/yfsc-server-macos-arm64 .
	cd server && CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -tags embed -trimpath -ldflags "-s -w" -o ../dist/yfsc-server-macos-amd64 .
	cd server && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags embed -trimpath -ldflags "-s -w" -o ../dist/yfsc-server-windows-amd64.exe .
	rm -rf server/webembed
	@echo "发布产物已生成到 dist/ (自包含, 双击或命令行运行即可)"

# Go 静态检查
test:
	cd server && go vet ./...

# 端到端 API 冒烟测试 (需先 make run 或后台启动服务)
smoke:
	bash server/smoke_test.sh

# 清空本地数据库 (重新灌入种子数据)
fresh:
	rm -f server/yfsc.db server/yfsc.db-wal server/yfsc.db-shm

clean: fresh
	rm -rf bin
