package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

// ---------------- 商品 ----------------
func handleProductList(w http.ResponseWriter, r *http.Request) {
	sc := scopeOf(getClaims(r))
	where := " WHERE 1=1"
	args := []interface{}{}
	if !sc.IsPlatform {
		if sc.ShopId > 0 {
			where += " AND pr.shop_id=?"
			args = append(args, sc.ShopId)
		} else if sc.ParkId > 0 {
			where += " AND s.park_id=?"
			args = append(args, sc.ParkId)
		}
	}
	if sid := queryInt(r, "shopId", 0); sid > 0 {
		where += " AND pr.shop_id=?"
		args = append(args, sid)
	}
	list, _ := queryMaps(`SELECT pr.id, pr.name, s.name shop, pr.price, pr.original_price,
		pr.stock, pr.sales_count, pr.status FROM product pr JOIN shop s ON s.id=pr.shop_id`+
		where+" ORDER BY pr.id", args...)
	ok(w, map[string]interface{}{"list": list, "total": len(list)})
}

func handleProductCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ShopId        int64   `json:"shopId"`
		Name          string  `json:"name"`
		Price         float64 `json:"price"`
		OriginalPrice float64 `json:"originalPrice"`
		Stock         int     `json:"stock"`
	}
	if err := decodeBody(r, &req); err != nil || req.Name == "" || req.ShopId == 0 {
		fail(w, 400, 400, "商品名称与所属店铺必填")
		return
	}
	res, err := db.Exec(`INSERT INTO product(shop_id,name,price,original_price,stock) VALUES(?,?,?,?,?)`,
		req.ShopId, req.Name, req.Price, req.OriginalPrice, req.Stock)
	if err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	id, _ := res.LastInsertId()
	okMsg(w, map[string]interface{}{"id": id}, "商品创建成功")
}

func handleProductUpdate(w http.ResponseWriter, r *http.Request) {
	id := pathInt(r, "id")
	var req struct {
		Name          string  `json:"name"`
		Price         float64 `json:"price"`
		OriginalPrice float64 `json:"originalPrice"`
		Stock         int     `json:"stock"`
		Status        string  `json:"status"`
	}
	if err := decodeBody(r, &req); err != nil {
		fail(w, 400, 400, "请求格式错误")
		return
	}
	_, err := db.Exec(`UPDATE product SET name=COALESCE(NULLIF(?,''),name),
		price=CASE WHEN ?>0 THEN ? ELSE price END,
		original_price=CASE WHEN ?>0 THEN ? ELSE original_price END,
		stock=CASE WHEN ?>=0 THEN ? ELSE stock END,
		status=COALESCE(NULLIF(?,''),status), updated_at=datetime('now','localtime') WHERE id=?`,
		req.Name, req.Price, req.Price, req.OriginalPrice, req.OriginalPrice,
		req.Stock, req.Stock, req.Status, id)
	if err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	writeLog(getClaims(r), "商品管理", "更新商品", fmt.Sprintf("product#%d", id))
	okMsg(w, nil, "更新成功")
}

func handleProductDelete(w http.ResponseWriter, r *http.Request) {
	id := pathInt(r, "id")
	if _, err := db.Exec("DELETE FROM product WHERE id=?", id); err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	writeLog(getClaims(r), "商品管理", "删除商品", r.PathValue("id"))
	okMsg(w, nil, "已删除")
}

// ---------------- 员工 / 岗位 ----------------
func handleEmployeeList(w http.ResponseWriter, r *http.Request) {
	sc := scopeOf(getClaims(r))
	where := " WHERE 1=1"
	args := []interface{}{}
	if !sc.IsPlatform {
		if sc.ShopId > 0 {
			where += " AND e.shop_id=?"
			args = append(args, sc.ShopId)
		} else {
			where += " AND e.park_id=?"
			args = append(args, sc.ParkId)
		}
	}
	list, _ := queryMaps(`SELECT e.id, e.name, e.username, COALESCE(e.phone,'') phone,
		COALESCE(p.name,'-') park, COALESCE(s.name,'-') shop,
		CASE e.role_code WHEN 'ADMIN_PLATFORM' THEN '平台超级管理员' WHEN 'PLATFORM_OPER' THEN '总部运营'
			WHEN 'PARK_ADMIN' THEN '园区管理员' WHEN 'PARK_MANAGER' THEN '园区经理'
			WHEN 'SHOP_ADMIN' THEN '商户管理员' WHEN 'SHOP_CASHIER' THEN '商户收银员' ELSE e.role_code END role,
		e.role_code, COALESCE(pos.name,'-') position, e.online_status online, e.status, e.created_at
		FROM employee e LEFT JOIN park p ON p.id=e.park_id LEFT JOIN shop s ON s.id=e.shop_id
		LEFT JOIN position pos ON pos.id=e.position_id`+where+" ORDER BY e.id", args...)
	ok(w, map[string]interface{}{"list": list, "total": len(list)})
}

func handleEmployeeCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username, Name, Phone, RoleCode, Password string
		ParkId, ShopId, PositionId               int64
	}
	if err := decodeBody(r, &req); err != nil || req.Username == "" || req.Name == "" || req.RoleCode == "" {
		fail(w, 400, 400, "用户名、姓名、角色必填")
		return
	}
	if req.Password == "" {
		req.Password = "123456"
	}
	res, err := db.Exec(`INSERT INTO employee(username,name,phone,park_id,shop_id,position_id,role_code,password_hash)
		VALUES(?,?,?,?,?,?,?,?)`, req.Username, req.Name, req.Phone, nullID(req.ParkId), nullID(req.ShopId),
		nullID(req.PositionId), req.RoleCode, hashPassword(req.Password))
	if err != nil {
		fail(w, 400, 400, "创建失败(用户名可能重复): "+err.Error())
		return
	}
	id, _ := res.LastInsertId()
	writeLog(getClaims(r), "员工管理", "新增员工", req.Name+"("+req.Username+")")
	okMsg(w, map[string]interface{}{"id": id}, "员工创建成功")
}

func handleEmployeeUpdate(w http.ResponseWriter, r *http.Request) {
	id := pathInt(r, "id")
	var req struct {
		Name, Phone, RoleCode, Status     string
		ParkId, ShopId, PositionId        int64
	}
	if err := decodeBody(r, &req); err != nil {
		fail(w, 400, 400, "请求格式错误")
		return
	}
	_, err := db.Exec(`UPDATE employee SET name=COALESCE(NULLIF(?,''),name),
		phone=COALESCE(NULLIF(?,''),phone), role_code=COALESCE(NULLIF(?,''),role_code),
		status=COALESCE(NULLIF(?,''),status), park_id=?, shop_id=?, position_id=?,
		updated_at=datetime('now','localtime') WHERE id=?`,
		req.Name, req.Phone, req.RoleCode, req.Status,
		nullID(req.ParkId), nullID(req.ShopId), nullID(req.PositionId), id)
	if err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	writeLog(getClaims(r), "员工管理", "更新员工", fmt.Sprintf("emp#%d", id))
	okMsg(w, nil, "更新成功")
}

func handleEmployeeDelete(w http.ResponseWriter, r *http.Request) {
	id := pathInt(r, "id")
	var role string
	_ = db.QueryRow("SELECT role_code FROM employee WHERE id=?", id).Scan(&role)
	if role == "ADMIN_PLATFORM" && scalarInt("SELECT COUNT(*) FROM employee WHERE role_code='ADMIN_PLATFORM'") <= 1 {
		fail(w, 400, 400, "不可删除最后一个平台超级管理员")
		return
	}
	if _, err := db.Exec("DELETE FROM employee WHERE id=?", id); err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	writeLog(getClaims(r), "员工管理", "删除员工", r.PathValue("id"))
	okMsg(w, nil, "已删除")
}

func handlePositionList(w http.ResponseWriter, r *http.Request) {
	list, _ := queryMaps(`SELECT id,name,code,COALESCE(responsibilities,'') responsibilities,status,created_at,
		(SELECT COUNT(*) FROM employee WHERE position_id=position.id) employee_count FROM position ORDER BY id`)
	ok(w, map[string]interface{}{"list": list, "total": len(list)})
}

// ---------------- 闸机 / 人脸 / 进出 ----------------
func handleGateList(w http.ResponseWriter, r *http.Request) {
	list, _ := queryMaps(`SELECT g.id, g.name, g.type, g.direction, g.device_sn, g.enabled,
		g.online_status, g.today_pass, COALESCE(p.name,'') park FROM gate g LEFT JOIN park p ON p.id=g.park_id ORDER BY g.id`)
	ok(w, map[string]interface{}{"list": list, "total": len(list)})
}

func handleGateCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ParkId     int64  `json:"parkId"`
		Name       string `json:"name"`
		Type       string `json:"type"`
		Direction  string `json:"direction"`
		DeviceSn   string `json:"deviceSn"`
	}
	if err := decodeBody(r, &req); err != nil || req.Name == "" || req.DeviceSn == "" {
		fail(w, 400, 400, "闸机名称与设备SN必填")
		return
	}
	sc := scopeOf(getClaims(r))
	if !sc.IsPlatform && sc.ParkId > 0 {
		req.ParkId = sc.ParkId // 园区视角强制归属本园区
	}
	if req.Type == "" {
		req.Type = "扫码闸机"
	}
	res, err := db.Exec(`INSERT INTO gate(park_id,name,type,direction,device_sn,enabled,online_status,today_pass)
		VALUES(?,?,?,?,?,1,0,0)`, nullID(req.ParkId), req.Name, req.Type, req.Direction, req.DeviceSn)
	if err != nil {
		fail(w, 400, 400, "创建失败(SN可能重复): "+err.Error())
		return
	}
	id, _ := res.LastInsertId()
	writeLog(getClaims(r), "闸机管理", "新增闸机", req.Name)
	okMsg(w, map[string]interface{}{"id": id}, "闸机已添加")
}

func handleGateToggle(w http.ResponseWriter, r *http.Request) {
	id := pathInt(r, "id")
	var cur int
	if err := db.QueryRow("SELECT enabled FROM gate WHERE id=?", id).Scan(&cur); err != nil {
		fail(w, 404, 404, "闸机不存在")
		return
	}
	nv := 1 - cur
	_, _ = db.Exec("UPDATE gate SET enabled=? WHERE id=?", nv, id)
	writeLog(getClaims(r), "闸机管理", "启停闸机", map[int]string{0: "禁用", 1: "启用"}[nv])
	okMsg(w, map[string]interface{}{"enabled": nv}, "状态已更新")
}

func handleGateDelete(w http.ResponseWriter, r *http.Request) {
	id := pathInt(r, "id")
	if _, err := db.Exec("DELETE FROM gate WHERE id=?", id); err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	writeLog(getClaims(r), "闸机管理", "删除闸机", r.PathValue("id"))
	okMsg(w, nil, "已删除")
}

func handleActivityUpdate(w http.ResponseWriter, r *http.Request) {
	id := pathInt(r, "id")
	var req struct {
		Title, Type, ParkScope, ShopScope, StartAt, EndAt string
	}
	if err := decodeBody(r, &req); err != nil {
		fail(w, 400, 400, "请求格式错误")
		return
	}
	_, err := db.Exec(`UPDATE activity SET title=COALESCE(NULLIF(?,''),title),
		type=COALESCE(NULLIF(?,''),type), park_scope=COALESCE(NULLIF(?,''),park_scope),
		shop_scope=COALESCE(NULLIF(?,''),shop_scope), start_at=COALESCE(NULLIF(?,''),start_at),
		end_at=COALESCE(NULLIF(?,''),end_at) WHERE id=?`,
		req.Title, req.Type, req.ParkScope, req.ShopScope, req.StartAt, req.EndAt, id)
	if err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	writeLog(getClaims(r), "活动管理", "更新活动", r.PathValue("id"))
	okMsg(w, nil, "更新成功")
}

func handleActivityDelete(w http.ResponseWriter, r *http.Request) {
	id := pathInt(r, "id")
	if _, err := db.Exec("DELETE FROM activity WHERE id=?", id); err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	writeLog(getClaims(r), "活动管理", "删除活动", r.PathValue("id"))
	okMsg(w, nil, "已删除")
}

func handleActivityAudit(w http.ResponseWriter, r *http.Request) {
	id := pathInt(r, "id")
	var req struct{ Approve bool `json:"approve"` }
	_ = decodeBody(r, &req)
	var status string
	if err := db.QueryRow("SELECT status FROM activity WHERE id=?", id).Scan(&status); err != nil {
		fail(w, 404, 404, "活动不存在")
		return
	}
	if status != "待审核" {
		fail(w, 400, 400, "仅【待审核】活动可审核，当前："+status)
		return
	}
	ns := "已驳回"
	if req.Approve {
		ns = "进行中"
	}
	_, _ = db.Exec("UPDATE activity SET status=? WHERE id=?", ns, id)
	writeLog(getClaims(r), "活动管理", "审核活动", r.PathValue("id")+" -> "+ns)
	okMsg(w, map[string]interface{}{"status": ns}, "审核完成："+ns)
}

func handleActivityEnd(w http.ResponseWriter, r *http.Request) {
	id := pathInt(r, "id")
	var status string
	if err := db.QueryRow("SELECT status FROM activity WHERE id=?", id).Scan(&status); err != nil {
		fail(w, 404, 404, "活动不存在")
		return
	}
	if status != "进行中" {
		fail(w, 400, 400, "仅【进行中】活动可结束，当前："+status)
		return
	}
	_, _ = db.Exec("UPDATE activity SET status='已结束' WHERE id=?", id)
	writeLog(getClaims(r), "活动管理", "结束活动", r.PathValue("id"))
	okMsg(w, map[string]interface{}{"status": "已结束"}, "活动已结束")
}

func handleFaceList(w http.ResponseWriter, r *http.Request) {
	list, _ := queryMaps(`SELECT f.id, f.user_id, u.phone, f.feature_id, f.authorized_gates, f.status, f.created_at
		FROM face f JOIN user u ON u.id=f.user_id ORDER BY f.id`)
	ok(w, map[string]interface{}{"list": list, "total": len(list)})
}

func handleEntryRecordList(w http.ResponseWriter, r *http.Request) {
	list, _ := queryMaps(`SELECT er.id, u.nickname user, u.phone, g.name gate, er.direction, er.method, er.passed_at
		FROM entry_record er LEFT JOIN user u ON u.id=er.user_id LEFT JOIN gate g ON g.id=er.gate_id ORDER BY er.id DESC`)
	ok(w, map[string]interface{}{"list": list, "total": len(list)})
}

// ---------------- 活动 ----------------
func handleActivityList(w http.ResponseWriter, r *http.Request) {
	where := ""
	args := []interface{}{}
	if st := r.URL.Query().Get("status"); st != "" {
		where = " WHERE status=?"
		args = append(args, st)
	}
	list, _ := queryMaps(`SELECT id,title,type,park_scope,shop_scope,start_at,end_at,status,
		exposure,joins,orders,verify_rate FROM activity`+where+" ORDER BY id", args...)
	ok(w, map[string]interface{}{"list": list, "total": len(list)})
}

func handleActivityCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title, Type, ParkScope, ShopScope, StartAt, EndAt string
	}
	if err := decodeBody(r, &req); err != nil || req.Title == "" {
		fail(w, 400, 400, "活动标题必填")
		return
	}
	res, err := db.Exec(`INSERT INTO activity(title,type,park_scope,shop_scope,start_at,end_at,status)
		VALUES(?,?,?,?,?,?, '待审核')`, req.Title, req.Type, req.ParkScope, req.ShopScope, req.StartAt, req.EndAt)
	if err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	id, _ := res.LastInsertId()
	writeLog(getClaims(r), "活动管理", "创建活动", req.Title)
	okMsg(w, map[string]interface{}{"id": id}, "活动已创建，等待审核")
}

// ---------------- 第三方平台核销 ----------------
func handlePlatformRecords(w http.ResponseWriter, r *http.Request) {
	platform := r.URL.Query().Get("platform")
	sc := scopeOf(getClaims(r))
	where := " WHERE 1=1"
	args := []interface{}{}
	if platform != "" {
		where += " AND platform=?"
		args = append(args, platform)
	}
	if !sc.IsPlatform && sc.ShopId > 0 {
		where += " AND shop_id=?"
		args = append(args, sc.ShopId)
	}
	if st := r.URL.Query().Get("status"); st != "" {
		where += " AND status=?"
		args = append(args, st)
	}
	list, _ := queryMaps(`SELECT biz_no id, platform, coupon_code, COALESCE(product,'') product,
		COALESCE(origin_price,0) origin_price, COALESCE(sale_price,0) sale_price,
		COALESCE(user_name,'') user, COALESCE(phone,'') phone, COALESCE(shop_name,'') shop,
		COALESCE(coupon_type,'团购套餐') coupon_type, verify_count, max_count, status,
		COALESCE(verify_by,'') verify_by, COALESCE(verify_at,'') verify_at, refundable
		FROM platform_verify_records`+where+" ORDER BY id DESC", args...)
	ok(w, map[string]interface{}{"list": list, "total": len(list)})
}

func handlePlatformSettlements(w http.ResponseWriter, r *http.Request) {
	platform := r.URL.Query().Get("platform")
	where := ""
	args := []interface{}{}
	if platform != "" {
		where = " WHERE platform=?"
		args = append(args, platform)
	}
	list, _ := queryMaps(`SELECT biz_no id, platform, shop_name shop, settle_date, order_count,
		total_amount, commission, net_amount, status, period, rate FROM platform_settlements`+where+" ORDER BY id DESC", args...)
	ok(w, map[string]interface{}{"list": list, "total": len(list)})
}

func handlePlatformStoreConfigs(w http.ResponseWriter, r *http.Request) {
	platform := r.URL.Query().Get("platform")
	where := ""
	args := []interface{}{}
	if platform != "" {
		where = " WHERE platform=?"
		args = append(args, platform)
	}
	list, _ := queryMaps(`SELECT id, platform, shop_name shop, COALESCE(app_auth_token,'') app_auth_token,
		status, COALESCE(bind_at,'-') bind_at, COALESCE(last_sync,'-') last_sync, verify_count, COALESCE(scope,'-') scope
		FROM platform_store_configs`+where+" ORDER BY id", args...)
	ok(w, map[string]interface{}{"list": list, "total": len(list)})
}

// 验券准备: 校验门店绑定并返回券信息 (模拟第三方应答)
func handleVerifyPrepare(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Platform   string `json:"platform"`
		CouponCode string `json:"couponCode"`
		ShopId     int64  `json:"shopId"`
	}
	if err := decodeBody(r, &req); err != nil || req.CouponCode == "" || req.Platform == "" {
		fail(w, 400, 400, "平台与券码必填")
		return
	}
	var status string
	err := db.QueryRow(`SELECT status FROM platform_store_configs WHERE platform=? AND shop_id=?`,
		req.Platform, req.ShopId).Scan(&status)
	if err != nil || status != "已绑定" {
		fail(w, 400, 400, "该门店未绑定此平台")
		return
	}
	// 已存在的券: 回放真实记录; 否则给出模拟券信息
	rec, _ := queryMaps(`SELECT product, origin_price, sale_price, coupon_type, verify_count, max_count, status
		FROM platform_verify_records WHERE platform=? AND coupon_code=? LIMIT 1`, req.Platform, req.CouponCode)
	if len(rec) > 0 {
		ok(w, rec[0])
		return
	}
	ok(w, map[string]interface{}{
		"couponCode": req.CouponCode, "product": "团购套餐(待核销)",
		"originPrice": 198.0, "salePrice": 128.0, "couponType": "团购套餐",
		"verifyCount": 0, "maxCount": 1, "status": "待核销",
	})
}

// 执行验券: 写入/更新核销记录
func handleVerifyExecute(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Platform   string `json:"platform"`
		CouponCode string `json:"couponCode"`
		ShopId     int64  `json:"shopId"`
		Product    string `json:"product"`
		SalePrice  float64 `json:"salePrice"`
	}
	if err := decodeBody(r, &req); err != nil || req.CouponCode == "" {
		fail(w, 400, 400, "券码必填")
		return
	}
	c := getClaims(r)
	var shopName string
	_ = db.QueryRow(`SELECT name FROM shop WHERE id=?`, req.ShopId).Scan(&shopName)
	var id int64
	var cur string
	err := db.QueryRow(`SELECT id,status FROM platform_verify_records WHERE platform=? AND coupon_code=?`,
		req.Platform, req.CouponCode).Scan(&id, &cur)
	bizNo := strings_ToUpper(req.Platform[:2]) + time.Now().Format("20060102150405")
	if err == nil {
		if cur == "已核销" {
			fail(w, 400, 400, "该券已核销")
			return
		}
		_, _ = db.Exec(`UPDATE platform_verify_records SET status='已核销', verify_by=?, verify_at=datetime('now','localtime') WHERE id=?`, c.Name, id)
	} else {
		if req.Product == "" {
			req.Product = "团购套餐"
		}
		_, _ = db.Exec(`INSERT INTO platform_verify_records(biz_no,platform,coupon_code,product,origin_price,sale_price,shop_id,shop_name,coupon_type,status,verify_by,verify_at,refundable)
			VALUES(?,?,?,?,?,?,?,?, '团购套餐', '已核销', ?, datetime('now','localtime'), 1)`,
			bizNo, req.Platform, req.CouponCode, req.Product, req.SalePrice, req.SalePrice, req.ShopId, shopName, c.Name)
	}
	_, _ = db.Exec(`UPDATE platform_store_configs SET verify_count=verify_count+1, last_sync=datetime('now','localtime') WHERE platform=? AND shop_id=?`, req.Platform, req.ShopId)
	writeLog(c, "第三方核销", "执行验券", req.Platform+":"+req.CouponCode)
	okMsg(w, map[string]interface{}{"couponCode": req.CouponCode, "salePrice": req.SalePrice}, "核销成功")
}

func handleVerifyRevoke(w http.ResponseWriter, r *http.Request) {
	bizNo := r.PathValue("id")
	var status string
	var refundable int
	if err := db.QueryRow(`SELECT status,refundable FROM platform_verify_records WHERE biz_no=?`, bizNo).Scan(&status, &refundable); err != nil {
		fail(w, 404, 404, "核销记录不存在")
		return
	}
	if status == "已撤销" {
		fail(w, 400, 400, "该券已撤销")
		return
	}
	if refundable == 0 {
		fail(w, 400, 400, "该券不支持撤销")
		return
	}
	_, _ = db.Exec(`UPDATE platform_verify_records SET status='已撤销' WHERE biz_no=?`, bizNo)
	writeLog(getClaims(r), "第三方核销", "撤销验券", bizNo)
	okMsg(w, nil, "已撤销")
}

// ---------------- 系统 ----------------
func handleMiniProgramList(w http.ResponseWriter, r *http.Request) {
	list, _ := queryMaps(`SELECT id,app_id,name,type,type_label,app_secret,callback_url,mch_id,status,version,published_at
		FROM mini_program ORDER BY id`)
	ok(w, map[string]interface{}{"list": list, "total": len(list)})
}

func handleMiniProgramUpdate(w http.ResponseWriter, r *http.Request) {
	id := pathInt(r, "id")
	var req struct {
		CallbackUrl, MchId, Status, Version string
	}
	if err := decodeBody(r, &req); err != nil {
		fail(w, 400, 400, "请求格式错误")
		return
	}
	_, err := db.Exec(`UPDATE mini_program SET callback_url=COALESCE(NULLIF(?,''),callback_url),
		mch_id=COALESCE(NULLIF(?,''),mch_id), status=COALESCE(NULLIF(?,''),status),
		version=COALESCE(NULLIF(?,''),version) WHERE id=?`,
		req.CallbackUrl, req.MchId, req.Status, req.Version, id)
	if err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	writeLog(getClaims(r), "系统设置", "更新小程序配置", r.PathValue("id"))
	okMsg(w, nil, "配置已更新")
}

func handleLogList(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	size := queryInt(r, "size", 30)
	total := scalarInt("SELECT COUNT(*) FROM operation_log")
	list, _ := queryMaps(`SELECT id,operator,role_code,module,action,COALESCE(detail,'') detail,created_at
		FROM operation_log ORDER BY id DESC LIMIT ? OFFSET ?`, size, (page-1)*size)
	ok(w, listData{List: list, Total: total, Page: page, Size: size})
}

// 权限矩阵: 按角色返回模块可见性
func handlePermissionList(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	where := ""
	args := []interface{}{}
	if role != "" {
		where = " WHERE role_code=?"
		args = append(args, role)
	}
	list, _ := queryMaps(`SELECT role_code, module_key, visibility FROM role_permission`+where+" ORDER BY role_code, module_key", args...)
	ok(w, map[string]interface{}{"list": list, "total": len(list)})
}

func handlePermissionUpdate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoleCode   string `json:"roleCode"`
		ModuleKey  string `json:"moduleKey"`
		Visibility string `json:"visibility"`
	}
	if err := decodeBody(r, &req); err != nil || req.RoleCode == "" || req.ModuleKey == "" {
		fail(w, 400, 400, "角色与模块必填")
		return
	}
	if req.Visibility != "可操作" && req.Visibility != "只读" && req.Visibility != "隐藏" {
		fail(w, 400, 400, "可见性取值非法")
		return
	}
	_, err := db.Exec(`INSERT INTO role_permission(role_code,module_key,visibility) VALUES(?,?,?)
		ON CONFLICT(role_code,module_key) DO UPDATE SET visibility=excluded.visibility, updated_at=datetime('now','localtime')`,
		req.RoleCode, req.ModuleKey, req.Visibility)
	if err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	writeLog(getClaims(r), "权限管理", "调整权限", req.RoleCode+"/"+req.ModuleKey+"="+req.Visibility)
	okMsg(w, nil, "权限已更新")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	var ver string
	_ = db.QueryRow("SELECT sqlite_version()").Scan(&ver)
	ok(w, map[string]interface{}{
		"status": "ok", "time": time.Now().Format(time.RFC3339), "sqlite": ver,
	})
}

// 占位以避免误用 strings 包大写转换的小工具 (platform 前缀)
func strings_ToUpper(s string) string {
	out := []rune(s)
	for i, c := range out {
		if c >= 'a' && c <= 'z' {
			out[i] = c - 32
		}
	}
	return string(out)
}

var _ = sql.ErrNoRows
