.PHONY: run build test smoke clean fresh

# 启动平台 (API + 静态前端) → http://localhost:8080
run:
	cd server && go run . -addr :8080 -web ../web -db yfsc.db

# 编译单一可执行文件到 bin/yfsc-server
build:
	cd server && go build -o ../bin/yfsc-server .

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
