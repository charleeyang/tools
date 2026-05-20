package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

type ctxKey string

const claimsKey ctxKey = "claims"

func getClaims(r *http.Request) *employeeClaims {
	c, _ := r.Context().Value(claimsKey).(*employeeClaims)
	return c
}

// 视角范围: 总部 -> 全部; 园区 -> 本园区; 商户 -> 本店铺
type viewScope struct {
	IsPlatform bool
	ParkId     int64
	ShopId     int64
}

func scopeOf(c *employeeClaims) viewScope {
	return viewScope{
		IsPlatform: c.ParkId == 0 && c.ShopId == 0,
		ParkId:     c.ParkId,
		ShopId:     c.ShopId,
	}
}

type middleware func(http.HandlerFunc) http.HandlerFunc

func chain(h http.HandlerFunc, ms ...middleware) http.HandlerFunc {
	for i := len(ms) - 1; i >= 0; i-- {
		h = ms[i](h)
	}
	return h
}

// JWT 认证
func authMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tok := bearerToken(r)
		if tok == "" {
			fail(w, http.StatusUnauthorized, 401, "未登录")
			return
		}
		claims, err := parseToken(tok)
		if err != nil {
			fail(w, http.StatusUnauthorized, 401, "Token无效或已过期")
			return
		}
		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// RBAC: 检查角色对模块的可见性
//   editable=true  要求 "可操作"
//   editable=false 要求 不为 "隐藏" (即可见: 可操作/只读)
func requirePerm(moduleKey string, editable bool) middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			c := getClaims(r)
			if c == nil {
				fail(w, http.StatusUnauthorized, 401, "未登录")
				return
			}
			var vis string
			_ = db.QueryRow(
				"SELECT visibility FROM role_permission WHERE role_code=? AND module_key=?",
				c.RoleCode, moduleKey,
			).Scan(&vis)
			if vis == "" || vis == "隐藏" {
				fail(w, http.StatusForbidden, 403, "无访问权限")
				return
			}
			if editable && vis == "只读" {
				fail(w, http.StatusForbidden, 403, "当前视角无操作权限")
				return
			}
			next.ServeHTTP(w, r)
		}
	}
}

func corsMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func logMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		if len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api" {
			log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start).Round(time.Millisecond))
		}
	})
}
