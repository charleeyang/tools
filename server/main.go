package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", envOr("ADDR", ":8080"), "监听地址")
	dsn := flag.String("db", envOr("DB_PATH", "yfsc.db"), "SQLite 文件路径")
	webDir := flag.String("web", envOr("WEB_DIR", "../web"), "静态前端目录")
	flag.Parse()

	initDB(*dsn)

	mux := http.NewServeMux()
	registerAPI(mux)

	// 静态前端 (原型: PC后台 + 三端小程序)
	resolved := resolveWebDir(*webDir)
	if resolved != "" {
		mux.Handle("/", http.FileServer(http.Dir(resolved)))
		log.Printf("静态前端目录: %s", resolved)
	} else {
		log.Printf("未找到前端目录，仅提供 API (尝试过 %s / web / ../web)", *webDir)
	}

	handler := corsMW(logMW(mux))
	log.Printf("袁夫稻田智慧园区平台已启动 → http://localhost%s", *addr)
	log.Printf("默认账号: admin / admin123 (平台超管)")
	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal(err)
	}
}

// resolveWebDir 依次尝试候选目录，兼容从仓库根或 server/ 启动
func resolveWebDir(preferred string) string {
	for _, c := range []string{preferred, "web", "../web", "./web"} {
		if c == "" {
			continue
		}
		if st, err := os.Stat(c); err == nil && st.IsDir() {
			return c
		}
	}
	return ""
}

func registerAPI(mux *http.ServeMux) {
	// 公共
	mux.HandleFunc("GET /api/health", handleHealth)
	mux.HandleFunc("POST /api/auth/login", handleLogin)
	mux.HandleFunc("GET /api/auth/me", chain(handleMe, authMW))

	auth := func(h http.HandlerFunc) http.HandlerFunc { return chain(h, authMW) }
	// 模块读权限(可见即可) / 写权限(需可操作)
	rd := func(mod string, h http.HandlerFunc) http.HandlerFunc { return chain(h, authMW, requirePerm(mod, false)) }
	wr := func(mod string, h http.HandlerFunc) http.HandlerFunc { return chain(h, authMW, requirePerm(mod, true)) }

	// 统计
	mux.HandleFunc("GET /api/statistics/overview", auth(handleOverview))

	// 财务汇总
	mux.HandleFunc("GET /api/finance/summary", rd("finance", handleFinanceSummary))

	// 园区
	mux.HandleFunc("GET /api/parks", rd("parks", handleParkList))
	mux.HandleFunc("POST /api/parks", wr("parks", handleParkCreate))
	mux.HandleFunc("PUT /api/parks/{id}", wr("parks", handleParkUpdate))
	mux.HandleFunc("DELETE /api/parks/{id}", wr("parks", handleParkDelete))

	// 店铺
	mux.HandleFunc("GET /api/shops", rd("shops", handleShopList))
	mux.HandleFunc("POST /api/shops", wr("shops", handleShopCreate))
	mux.HandleFunc("PUT /api/shops/{id}", wr("shops", handleShopUpdate))
	mux.HandleFunc("DELETE /api/shops/{id}", wr("shops", handleShopDelete))
	mux.HandleFunc("GET /api/shop-types", rd("shops", handleShopTypeList))
	mux.HandleFunc("POST /api/shop-types", wr("shops", handleShopTypeCreate))

	// 客户
	mux.HandleFunc("GET /api/users", rd("users", handleUserList))
	mux.HandleFunc("GET /api/users/stats", rd("users", handleUserStats))
	mux.HandleFunc("POST /api/users", wr("users", handleUserCreate))
	mux.HandleFunc("PUT /api/users/{id}", wr("users", handleUserUpdate))
	mux.HandleFunc("DELETE /api/users/{id}", wr("users", handleUserDelete))
	mux.HandleFunc("GET /api/recharges", rd("users", handleRechargeList))

	// 消费
	mux.HandleFunc("GET /api/orders", rd("consumption", handleOrderList))
	mux.HandleFunc("POST /api/orders/{no}/verify", wr("consumption", handleOrderVerify))
	mux.HandleFunc("POST /api/orders/{no}/refund", wr("consumption", handleOrderRefund))
	mux.HandleFunc("GET /api/refunds", rd("consumption", handleRefundList))
	mux.HandleFunc("POST /api/refunds/{id}/approve", wr("consumption", handleRefundReview(true)))
	mux.HandleFunc("POST /api/refunds/{id}/reject", wr("consumption", handleRefundReview(false)))

	// 财务 / 提现
	mux.HandleFunc("GET /api/withdraws", rd("finance", handleWithdrawList))
	mux.HandleFunc("POST /api/withdraws/{id}/approve", wr("finance", handleWithdrawReview(true)))
	mux.HandleFunc("POST /api/withdraws/{id}/reject", wr("finance", handleWithdrawReview(false)))

	// 商品
	mux.HandleFunc("GET /api/products", rd("products", handleProductList))
	mux.HandleFunc("POST /api/products", wr("products", handleProductCreate))
	mux.HandleFunc("PUT /api/products/{id}", wr("products", handleProductUpdate))
	mux.HandleFunc("DELETE /api/products/{id}", wr("products", handleProductDelete))

	// 员工 / 岗位
	mux.HandleFunc("GET /api/employees", rd("employees", handleEmployeeList))
	mux.HandleFunc("POST /api/employees", wr("employees", handleEmployeeCreate))
	mux.HandleFunc("PUT /api/employees/{id}", wr("employees", handleEmployeeUpdate))
	mux.HandleFunc("DELETE /api/employees/{id}", wr("employees", handleEmployeeDelete))
	mux.HandleFunc("GET /api/positions", rd("employees", handlePositionList))

	// 闸机
	mux.HandleFunc("GET /api/gates", rd("gates", handleGateList))
	mux.HandleFunc("POST /api/gates", wr("gates", handleGateCreate))
	mux.HandleFunc("PUT /api/gates/{id}/toggle", wr("gates", handleGateToggle))
	mux.HandleFunc("DELETE /api/gates/{id}", wr("gates", handleGateDelete))
	mux.HandleFunc("GET /api/faces", rd("gates", handleFaceList))
	mux.HandleFunc("GET /api/entry-records", rd("gates", handleEntryRecordList))

	// 活动
	mux.HandleFunc("GET /api/activities", rd("activities", handleActivityList))
	mux.HandleFunc("POST /api/activities", wr("activities", handleActivityCreate))
	mux.HandleFunc("PUT /api/activities/{id}", wr("activities", handleActivityUpdate))
	mux.HandleFunc("DELETE /api/activities/{id}", wr("activities", handleActivityDelete))
	mux.HandleFunc("POST /api/activities/{id}/audit", wr("activities", handleActivityAudit))
	mux.HandleFunc("POST /api/activities/{id}/end", wr("activities", handleActivityEnd))

	// 第三方平台 (美团/抖音)
	mux.HandleFunc("GET /api/platform/verify-records", auth(handlePlatformRecords))
	mux.HandleFunc("GET /api/platform/settlements", auth(handlePlatformSettlements))
	mux.HandleFunc("GET /api/platform/store-configs", auth(handlePlatformStoreConfigs))
	mux.HandleFunc("POST /api/platform/verify/prepare", auth(handleVerifyPrepare))
	mux.HandleFunc("POST /api/platform/verify/execute", auth(handleVerifyExecute))
	mux.HandleFunc("POST /api/platform/verify/{id}/revoke", auth(handleVerifyRevoke))

	// 系统 / 权限
	mux.HandleFunc("GET /api/system/mini-programs", rd("settings", handleMiniProgramList))
	mux.HandleFunc("PUT /api/system/mini-programs/{id}", wr("settings", handleMiniProgramUpdate))
	mux.HandleFunc("GET /api/system/logs", rd("settings", handleLogList))
	mux.HandleFunc("GET /api/permissions", rd("permissions", handlePermissionList))
	mux.HandleFunc("PUT /api/permissions", wr("permissions", handlePermissionUpdate))
}
