package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ---------------- 认证 ----------------
func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := decodeBody(r, &req); err != nil {
		fail(w, http.StatusBadRequest, 400, "请求格式错误")
		return
	}
	var (
		id                       int64
		name, roleCode, pwHash, status string
		parkId, shopId           sql.NullInt64
	)
	err := db.QueryRow(
		`SELECT id,name,role_code,password_hash,status,park_id,shop_id FROM employee WHERE username=?`,
		req.Username).Scan(&id, &name, &roleCode, &pwHash, &status, &parkId, &shopId)
	if err != nil {
		fail(w, http.StatusUnauthorized, 401, "用户名或密码错误")
		return
	}
	if status != "启用" {
		fail(w, http.StatusForbidden, 403, "账号已禁用")
		return
	}
	if !checkPassword(pwHash, req.Password) {
		fail(w, http.StatusUnauthorized, 401, "用户名或密码错误")
		return
	}
	claims := &employeeClaims{
		Id: id, Username: req.Username, Name: name, RoleCode: roleCode,
		ParkId: parkId.Int64, ShopId: shopId.Int64,
	}
	token, err := signToken(claims)
	if err != nil {
		fail(w, http.StatusInternalServerError, 500, "生成令牌失败")
		return
	}
	perms := map[string]string{}
	rows, _ := db.Query(`SELECT module_key,visibility FROM role_permission WHERE role_code=?`, roleCode)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var k, v string
			_ = rows.Scan(&k, &v)
			perms[k] = v
		}
	}
	writeLog(claims, "认证", "登录", "登录成功")
	ok(w, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id": id, "username": req.Username, "name": name,
			"roleCode": roleCode, "parkId": parkId.Int64, "shopId": shopId.Int64,
		},
		"permissions": perms,
	})
}

func handleMe(w http.ResponseWriter, r *http.Request) {
	c := getClaims(r)
	perms := map[string]string{}
	rows, _ := db.Query(`SELECT module_key,visibility FROM role_permission WHERE role_code=?`, c.RoleCode)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var k, v string
			_ = rows.Scan(&k, &v)
			perms[k] = v
		}
	}
	ok(w, map[string]interface{}{
		"user": map[string]interface{}{
			"id": c.Id, "username": c.Username, "name": c.Name,
			"roleCode": c.RoleCode, "parkId": c.ParkId, "shopId": c.ShopId,
		},
		"permissions": perms,
	})
}

// ---------------- 统计 ----------------
func handleOverview(w http.ResponseWriter, r *http.Request) {
	c := getClaims(r)
	sc := scopeOf(c)

	// 按视角范围过滤
	shopWhere, shopArgs := "", []interface{}{}
	if !sc.IsPlatform {
		if sc.ShopId > 0 {
			shopWhere, shopArgs = " WHERE id=?", []interface{}{sc.ShopId}
		} else {
			shopWhere, shopArgs = " WHERE park_id=?", []interface{}{sc.ParkId}
		}
	}
	totalShops := scalarInt("SELECT COUNT(*) FROM shop"+shopWhere, shopArgs...)

	totalParks := scalarInt("SELECT COUNT(*) FROM park")
	if !sc.IsPlatform {
		totalParks = 1
	}
	totalUsers := scalarInt("SELECT COUNT(*) FROM user")

	// 销售额 (含核销/已支付/已完成)
	paidStates := "('已支付','已核销','已完成')"
	orderWhere := " WHERE o.status IN " + paidStates
	orderArgs := []interface{}{}
	if !sc.IsPlatform {
		if sc.ShopId > 0 {
			orderWhere += " AND o.shop_id=?"
			orderArgs = append(orderArgs, sc.ShopId)
		} else {
			orderWhere += " AND s.park_id=?"
			orderArgs = append(orderArgs, sc.ParkId)
		}
	}
	totalSales := scalarFloat(`SELECT COALESCE(SUM(o.amount),0) FROM "order" o JOIN shop s ON s.id=o.shop_id`+orderWhere, orderArgs...)
	orderCount := scalarInt(`SELECT COUNT(*) FROM "order" o JOIN shop s ON s.id=o.shop_id`+orderWhere, orderArgs...)

	// 近7天销售趋势 (按天聚合，无数据补 0)
	dates := []string{}
	values := []float64{}
	base, _ := time.Parse("2006-01-02", "2026-05-14")
	for i := 6; i >= 0; i-- {
		d := base.AddDate(0, 0, -i).Format("2006-01-02")
		dates = append(dates, d[5:])
		v := scalarFloat(`SELECT COALESCE(SUM(o.amount),0) FROM "order" o JOIN shop s ON s.id=o.shop_id WHERE substr(o.created_at,1,10)=?`+
			strings.Replace(orderWhere, " WHERE o.status IN "+paidStates, " AND o.status IN "+paidStates, 1), append([]interface{}{d}, orderArgs...)...)
		values = append(values, v)
	}

	// 园区销售额分布
	parkDist, _ := queryMaps(`SELECT p.name, COALESCE(SUM(o.amount),0) AS value
		FROM park p LEFT JOIN shop s ON s.park_id=p.id
		LEFT JOIN "order" o ON o.shop_id=s.id AND o.status IN `+paidStates+`
		GROUP BY p.id ORDER BY value DESC`)

	ok(w, map[string]interface{}{
		"totalParks": totalParks, "totalShops": totalShops, "totalUsers": totalUsers,
		"todayNewUsers": 1, "todaySales": totalSales, "monthSales": totalSales,
		"orderCount": orderCount,
		"salesTrend": map[string]interface{}{"dates": dates, "values": values},
		"parkDistribution": parkDist,
	})
}

// ---------------- 园区 ----------------
func handleParkList(w http.ResponseWriter, r *http.Request) {
	list, err := queryMaps(`SELECT p.id,p.name,p.code,COALESCE(p.address,'') address,
		COALESCE(p.contact_person,'') contact_person, COALESCE(p.contact_phone,'') contact_phone,
		p.status, p.created_at, COUNT(s.id) shop_count
		FROM park p LEFT JOIN shop s ON s.park_id=p.id GROUP BY p.id ORDER BY p.id`)
	if err != nil {
		fail(w, 500, 500, err.Error())
		return
	}
	ok(w, map[string]interface{}{"list": list, "total": len(list)})
}

func handleParkCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name, Code, Address, ContactPerson, ContactPhone, BusinessHours, Description string
	}
	if err := decodeBody(r, &req); err != nil || req.Name == "" || req.Code == "" {
		fail(w, 400, 400, "园区名称与编码必填")
		return
	}
	res, err := db.Exec(`INSERT INTO park(name,code,address,contact_person,contact_phone,business_hours,description)
		VALUES(?,?,?,?,?,?,?)`, req.Name, req.Code, req.Address, req.ContactPerson, req.ContactPhone, req.BusinessHours, req.Description)
	if err != nil {
		fail(w, 400, 400, "创建失败(编码可能重复): "+err.Error())
		return
	}
	id, _ := res.LastInsertId()
	writeLog(getClaims(r), "园区管理", "新增园区", req.Name)
	okMsg(w, map[string]interface{}{"id": id}, "园区创建成功")
}

func handleParkUpdate(w http.ResponseWriter, r *http.Request) {
	id := pathInt(r, "id")
	var req struct {
		Name, Address, ContactPerson, ContactPhone, BusinessHours, Description, Status string
	}
	if err := decodeBody(r, &req); err != nil {
		fail(w, 400, 400, "请求格式错误")
		return
	}
	_, err := db.Exec(`UPDATE park SET name=COALESCE(NULLIF(?,''),name),
		address=?, contact_person=?, contact_phone=?, business_hours=?, description=?,
		status=COALESCE(NULLIF(?,''),status), updated_at=datetime('now','localtime') WHERE id=?`,
		req.Name, req.Address, req.ContactPerson, req.ContactPhone, req.BusinessHours, req.Description, req.Status, id)
	if err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	writeLog(getClaims(r), "园区管理", "更新园区", fmt.Sprintf("park#%d", id))
	okMsg(w, nil, "更新成功")
}

func handleParkDelete(w http.ResponseWriter, r *http.Request) {
	id := pathInt(r, "id")
	if scalarInt("SELECT COUNT(*) FROM shop WHERE park_id=?", id) > 0 {
		fail(w, 400, 400, "该园区下存在店铺，无法删除")
		return
	}
	if _, err := db.Exec("DELETE FROM park WHERE id=?", id); err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	writeLog(getClaims(r), "园区管理", "删除园区", fmt.Sprintf("park#%d", id))
	okMsg(w, nil, "已删除")
}

// ---------------- 店铺 ----------------
func handleShopList(w http.ResponseWriter, r *http.Request) {
	c := getClaims(r)
	sc := scopeOf(c)
	where := " WHERE 1=1"
	args := []interface{}{}
	if !sc.IsPlatform {
		if sc.ShopId > 0 {
			where += " AND s.id=?"
			args = append(args, sc.ShopId)
		} else {
			where += " AND s.park_id=?"
			args = append(args, sc.ParkId)
		}
	}
	if pid := queryInt(r, "parkId", 0); pid > 0 {
		where += " AND s.park_id=?"
		args = append(args, pid)
	}
	if st := r.URL.Query().Get("status"); st != "" {
		where += " AND s.status=?"
		args = append(args, st)
	}
	if kw := r.URL.Query().Get("keyword"); kw != "" {
		where += " AND s.name LIKE ?"
		args = append(args, "%"+kw+"%")
	}
	page := queryInt(r, "page", 1)
	size := queryInt(r, "size", 20)
	total := scalarInt("SELECT COUNT(*) FROM shop s"+where, args...)
	q := `SELECT s.id,s.park_id,p.name park_name,s.name,COALESCE(s.address,'') address,
		COALESCE(s.type_id,0) type_id, COALESCE(t.name,'') type_name, s.status, s.created_at
		FROM shop s JOIN park p ON p.id=s.park_id LEFT JOIN shop_type t ON t.id=s.type_id` +
		where + " ORDER BY s.id LIMIT ? OFFSET ?"
	list, err := queryMaps(q, append(args, size, (page-1)*size)...)
	if err != nil {
		fail(w, 500, 500, err.Error())
		return
	}
	ok(w, listData{List: list, Total: total, Page: page, Size: size})
}

func handleShopCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ParkId  int64  `json:"parkId"`
		Name    string `json:"name"`
		Address string `json:"address"`
		TypeId  int64  `json:"typeId"`
		Status  string `json:"status"`
	}
	if err := decodeBody(r, &req); err != nil || req.Name == "" || req.ParkId == 0 {
		fail(w, 400, 400, "店铺名称与所属园区必填")
		return
	}
	if req.Status == "" {
		req.Status = "营业中"
	}
	res, err := db.Exec(`INSERT INTO shop(park_id,name,address,type_id,status) VALUES(?,?,?,?,?)`,
		req.ParkId, req.Name, req.Address, nullID(req.TypeId), req.Status)
	if err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	id, _ := res.LastInsertId()
	writeLog(getClaims(r), "店铺管理", "新增店铺", req.Name)
	okMsg(w, map[string]interface{}{"id": id}, "店铺创建成功")
}

func handleShopUpdate(w http.ResponseWriter, r *http.Request) {
	id := pathInt(r, "id")
	var req struct {
		Name, Address, Status string
		TypeId                int64 `json:"typeId"`
	}
	if err := decodeBody(r, &req); err != nil {
		fail(w, 400, 400, "请求格式错误")
		return
	}
	_, err := db.Exec(`UPDATE shop SET name=COALESCE(NULLIF(?,''),name), address=?,
		type_id=COALESCE(NULLIF(?,0),type_id), status=COALESCE(NULLIF(?,''),status),
		updated_at=datetime('now','localtime') WHERE id=?`,
		req.Name, req.Address, req.TypeId, req.Status, id)
	if err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	writeLog(getClaims(r), "店铺管理", "更新店铺", fmt.Sprintf("shop#%d", id))
	okMsg(w, nil, "更新成功")
}

func handleShopDelete(w http.ResponseWriter, r *http.Request) {
	id := pathInt(r, "id")
	if scalarInt(`SELECT COUNT(*) FROM "order" WHERE shop_id=?`, id) > 0 {
		fail(w, 400, 400, "该店铺存在订单，无法删除（可改为已关闭）")
		return
	}
	_, _ = db.Exec("DELETE FROM product WHERE shop_id=?", id)
	if _, err := db.Exec("DELETE FROM shop WHERE id=?", id); err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	writeLog(getClaims(r), "店铺管理", "删除店铺", r.PathValue("id"))
	okMsg(w, nil, "已删除")
}

func handleUserUpdate(w http.ResponseWriter, r *http.Request) {
	id := pathInt(r, "id")
	var req struct {
		Nickname, Phone, MemberLevel, FaceStatus string
	}
	if err := decodeBody(r, &req); err != nil {
		fail(w, 400, 400, "请求格式错误")
		return
	}
	_, err := db.Exec(`UPDATE user SET nickname=COALESCE(NULLIF(?,''),nickname),
		phone=COALESCE(NULLIF(?,''),phone), member_level=COALESCE(NULLIF(?,''),member_level),
		face_status=COALESCE(NULLIF(?,''),face_status), updated_at=datetime('now','localtime') WHERE id=?`,
		req.Nickname, req.Phone, req.MemberLevel, req.FaceStatus, id)
	if err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	writeLog(getClaims(r), "客户管理", "更新客户", fmt.Sprintf("user#%d", id))
	okMsg(w, nil, "更新成功")
}

func handleUserDelete(w http.ResponseWriter, r *http.Request) {
	id := pathInt(r, "id")
	if scalarInt(`SELECT COUNT(*) FROM "order" WHERE user_id=?`, id) > 0 {
		fail(w, 400, 400, "该客户存在订单，无法删除")
		return
	}
	_, _ = db.Exec("DELETE FROM recharge WHERE user_id=?", id)
	_, _ = db.Exec("DELETE FROM face WHERE user_id=?", id)
	if _, err := db.Exec("DELETE FROM user WHERE id=?", id); err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	writeLog(getClaims(r), "客户管理", "删除客户", r.PathValue("id"))
	okMsg(w, nil, "已删除")
}

func handleShopTypeList(w http.ResponseWriter, r *http.Request) {
	list, _ := queryMaps(`SELECT id,name,code,COALESCE(description,'') description,status,created_at,
		(SELECT COUNT(*) FROM shop WHERE type_id=shop_type.id) shop_count FROM shop_type ORDER BY id`)
	ok(w, map[string]interface{}{"list": list, "total": len(list)})
}

func handleShopTypeCreate(w http.ResponseWriter, r *http.Request) {
	var req struct{ Name, Code, Description string }
	if err := decodeBody(r, &req); err != nil || req.Name == "" || req.Code == "" {
		fail(w, 400, 400, "类型名称与编码必填")
		return
	}
	if _, err := db.Exec(`INSERT INTO shop_type(name,code,description) VALUES(?,?,?)`, req.Name, req.Code, req.Description); err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	okMsg(w, nil, "创建成功")
}

// ---------------- 客户 ----------------
func handleUserList(w http.ResponseWriter, r *http.Request) {
	where := " WHERE 1=1"
	args := []interface{}{}
	if kw := r.URL.Query().Get("keyword"); kw != "" {
		where += " AND (nickname LIKE ? OR phone LIKE ?)"
		args = append(args, "%"+kw+"%", "%"+kw+"%")
	}
	if lv := r.URL.Query().Get("level"); lv != "" {
		where += " AND member_level=?"
		args = append(args, lv)
	}
	page := queryInt(r, "page", 1)
	size := queryInt(r, "size", 20)
	total := scalarInt("SELECT COUNT(*) FROM user"+where, args...)
	list, err := queryMaps(`SELECT id,nickname,phone,COALESCE(openid,'') openid,member_level,balance,
		consumption_total,points,face_status,registered_at FROM user`+where+" ORDER BY id DESC LIMIT ? OFFSET ?",
		append(args, size, (page-1)*size)...)
	if err != nil {
		fail(w, 500, 500, err.Error())
		return
	}
	ok(w, listData{List: list, Total: total, Page: page, Size: size})
}

func handleUserStats(w http.ResponseWriter, r *http.Request) {
	byLevel, _ := queryMaps(`SELECT member_level level, COUNT(*) count FROM user GROUP BY member_level`)
	ok(w, map[string]interface{}{
		"total":           scalarInt("SELECT COUNT(*) FROM user"),
		"vipCount":        scalarInt("SELECT COUNT(*) FROM user WHERE member_level LIKE 'VIP%'"),
		"faceRegistered":  scalarInt("SELECT COUNT(*) FROM user WHERE face_status='已录入'"),
		"totalBalance":    scalarFloat("SELECT COALESCE(SUM(balance),0) FROM user"),
		"totalConsumption": scalarFloat("SELECT COALESCE(SUM(consumption_total),0) FROM user"),
		"byLevel":         byLevel,
	})
}

func handleUserCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Nickname, Phone, MemberLevel string
	}
	if err := decodeBody(r, &req); err != nil || req.Phone == "" {
		fail(w, 400, 400, "手机号必填")
		return
	}
	if req.MemberLevel == "" {
		req.MemberLevel = "普通用户"
	}
	res, err := db.Exec(`INSERT INTO user(nickname,phone,member_level) VALUES(?,?,?)`, req.Nickname, req.Phone, req.MemberLevel)
	if err != nil {
		fail(w, 400, 400, err.Error())
		return
	}
	id, _ := res.LastInsertId()
	okMsg(w, map[string]interface{}{"id": id}, "客户创建成功")
}

// ---------------- 充值记录 ----------------
func handleRechargeList(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	size := queryInt(r, "size", 20)
	total := scalarInt("SELECT COUNT(*) FROM recharge")
	list, _ := queryMaps(`SELECT rc.biz_no id, u.phone, rc.amount, rc.gift_amount, rc.pay_method,
		rc.txn_no, COALESCE(rc.remark,'') remark, rc.status, rc.created_at at
		FROM recharge rc JOIN user u ON u.id=rc.user_id ORDER BY rc.id DESC LIMIT ? OFFSET ?`, size, (page-1)*size)
	ok(w, listData{List: list, Total: total, Page: page, Size: size})
}

// ---------------- 订单 ----------------
func orderScope(r *http.Request) (string, []interface{}) {
	sc := scopeOf(getClaims(r))
	where := " WHERE 1=1"
	args := []interface{}{}
	if !sc.IsPlatform {
		if sc.ShopId > 0 {
			where += " AND o.shop_id=?"
			args = append(args, sc.ShopId)
		} else {
			where += " AND s.park_id=?"
			args = append(args, sc.ParkId)
		}
	}
	return where, args
}

func handleOrderList(w http.ResponseWriter, r *http.Request) {
	where, args := orderScope(r)
	if st := r.URL.Query().Get("status"); st != "" {
		where += " AND o.status=?"
		args = append(args, st)
	}
	if kw := r.URL.Query().Get("keyword"); kw != "" {
		where += " AND o.order_no LIKE ?"
		args = append(args, "%"+kw+"%")
	}
	page := queryInt(r, "page", 1)
	size := queryInt(r, "size", 20)
	total := scalarInt(`SELECT COUNT(*) FROM "order" o JOIN shop s ON s.id=o.shop_id`+where, args...)
	list, _ := queryMaps(`SELECT o.order_no no, u.nickname user, s.name shop, o.amount, o.pay_method pay,
		o.status, o.source, o.created_at time
		FROM "order" o JOIN user u ON u.id=o.user_id JOIN shop s ON s.id=o.shop_id`+
		where+" ORDER BY o.id DESC LIMIT ? OFFSET ?", append(args, size, (page-1)*size)...)
	ok(w, listData{List: list, Total: total, Page: page, Size: size})
}

func handleOrderVerify(w http.ResponseWriter, r *http.Request) {
	no := r.PathValue("no")
	var st string
	if err := db.QueryRow(`SELECT status FROM "order" WHERE order_no=?`, no).Scan(&st); err != nil {
		fail(w, 404, 404, "订单不存在")
		return
	}
	if st != "已支付" {
		fail(w, 400, 400, "仅【已支付】订单可核销，当前状态："+st)
		return
	}
	_, _ = db.Exec(`UPDATE "order" SET status='已核销', updated_at=datetime('now','localtime') WHERE order_no=?`, no)
	writeLog(getClaims(r), "消费管理", "核销订单", no)
	okMsg(w, nil, "核销成功")
}

func handleOrderRefund(w http.ResponseWriter, r *http.Request) {
	no := r.PathValue("no")
	var req struct {
		Reason string `json:"reason"`
	}
	_ = decodeBody(r, &req)
	var (
		oid    int64
		uid    int64
		amount float64
		st     string
	)
	if err := db.QueryRow(`SELECT id,user_id,amount,status FROM "order" WHERE order_no=?`, no).Scan(&oid, &uid, &amount, &st); err != nil {
		fail(w, 404, 404, "订单不存在")
		return
	}
	if st == "已退款" || st == "退款处理中" {
		fail(w, 400, 400, "该订单已在退款流程中")
		return
	}
	tx, _ := db.Begin()
	bizNo := "RF" + time.Now().Format("20060102150405")
	if _, err := tx.Exec(`INSERT INTO refund(biz_no,user_id,order_id,order_no,amount,type,reason,status)
		VALUES(?,?,?,?,?, '订单退款', ?, '待审核')`, bizNo, uid, oid, no, amount, req.Reason); err != nil {
		tx.Rollback()
		fail(w, 500, 500, err.Error())
		return
	}
	if _, err := tx.Exec(`UPDATE "order" SET status='退款处理中', updated_at=datetime('now','localtime') WHERE id=?`, oid); err != nil {
		tx.Rollback()
		fail(w, 500, 500, err.Error())
		return
	}
	tx.Commit()
	writeLog(getClaims(r), "消费管理", "发起退款", no)
	okMsg(w, map[string]interface{}{"refundNo": bizNo}, "退款申请已提交")
}

// ---------------- 退款审核 ----------------
func handleRefundList(w http.ResponseWriter, r *http.Request) {
	where := ""
	args := []interface{}{}
	if st := r.URL.Query().Get("status"); st != "" {
		where = " WHERE rf.status=?"
		args = append(args, st)
	}
	list, _ := queryMaps(`SELECT rf.biz_no id, COALESCE(rf.order_no,'-') "order", u.nickname user,
		rf.amount, rf.type, COALESCE(rf.reason,'') reason, rf.status, rf.applied_at apply_at,
		COALESCE(rf.reviewer,'') reviewer FROM refund rf JOIN user u ON u.id=rf.user_id`+
		where+" ORDER BY rf.id DESC", args...)
	ok(w, map[string]interface{}{"list": list, "total": len(list)})
}

func handleRefundReview(approve bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bizNo := r.PathValue("id")
		var req struct{ Remark string `json:"remark"` }
		_ = decodeBody(r, &req)
		var (
			id, uid       int64
			rtype, status string
			amount        float64
			orderNo       sql.NullString
		)
		if err := db.QueryRow(`SELECT id,user_id,type,status,amount,order_no FROM refund WHERE biz_no=?`, bizNo).
			Scan(&id, &uid, &rtype, &status, &amount, &orderNo); err != nil {
			fail(w, 404, 404, "退款单不存在")
			return
		}
		if status != "待审核" {
			fail(w, 400, 400, "该退款单已处理，当前状态："+status)
			return
		}
		c := getClaims(r)
		tx, _ := db.Begin()
		newStatus := "已驳回"
		if approve {
			newStatus = "已通过"
		}
		_, _ = tx.Exec(`UPDATE refund SET status=?, reviewer=?, reviewed_at=datetime('now','localtime'), review_remark=? WHERE id=?`,
			newStatus, c.Name, req.Remark, id)
		if approve {
			if rtype == "余额退款" {
				tx.Exec(`UPDATE user SET balance=balance-? WHERE id=? AND balance>=?`, amount, uid, amount)
			} else if orderNo.Valid {
				tx.Exec(`UPDATE "order" SET status='已退款', updated_at=datetime('now','localtime') WHERE order_no=?`, orderNo.String)
			}
		} else if orderNo.Valid {
			tx.Exec(`UPDATE "order" SET status='已支付', updated_at=datetime('now','localtime') WHERE order_no=? AND status='退款处理中'`, orderNo.String)
		}
		tx.Commit()
		writeLog(c, "退款处理", "审核退款", bizNo+" -> "+newStatus)
		okMsg(w, nil, "审核完成："+newStatus)
	}
}

// ---------------- 提现审核 ----------------
func handleWithdrawList(w http.ResponseWriter, r *http.Request) {
	sc := scopeOf(getClaims(r))
	where := " WHERE 1=1"
	args := []interface{}{}
	if !sc.IsPlatform {
		if sc.ShopId > 0 {
			where += " AND w.shop_id=?"
			args = append(args, sc.ShopId)
		} else {
			where += " AND s.park_id=?"
			args = append(args, sc.ParkId)
		}
	}
	if st := r.URL.Query().Get("status"); st != "" {
		where += " AND w.status=?"
		args = append(args, st)
	}
	list, _ := queryMaps(`SELECT w.biz_no id, s.name shop, p.name shop_park, w.amount, w.bank_name,
		w.bank_account, w.bank_holder, w.status, w.applied_at apply_at, COALESCE(w.reconciliation_id,'') reconciliation_id,
		COALESCE(w.reviewer,'') reviewer, COALESCE(w.reviewed_at,'') reviewed_at, COALESCE(w.paid_at,'') paid_at
		FROM withdraw w JOIN shop s ON s.id=w.shop_id JOIN park p ON p.id=s.park_id`+
		where+" ORDER BY w.id DESC", args...)
	ok(w, map[string]interface{}{"list": list, "total": len(list)})
}

func handleWithdrawReview(approve bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bizNo := r.PathValue("id")
		var req struct{ Remark string `json:"remark"` }
		_ = decodeBody(r, &req)
		var status string
		var id int64
		if err := db.QueryRow(`SELECT id,status FROM withdraw WHERE biz_no=?`, bizNo).Scan(&id, &status); err != nil {
			fail(w, 404, 404, "提现单不存在")
			return
		}
		if status != "待审核" {
			fail(w, 400, 400, "该提现单已处理，当前状态："+status)
			return
		}
		c := getClaims(r)
		newStatus := "已驳回"
		if approve {
			newStatus = "已支付"
		}
		paidAt := "NULL"
		_ = paidAt
		if approve {
			_, _ = db.Exec(`UPDATE withdraw SET status=?, reviewer=?, reviewed_at=datetime('now','localtime'),
				paid_at=datetime('now','localtime'), review_remark=? WHERE id=?`, newStatus, c.Name, req.Remark, id)
		} else {
			_, _ = db.Exec(`UPDATE withdraw SET status=?, reviewer=?, reviewed_at=datetime('now','localtime'), review_remark=? WHERE id=?`,
				newStatus, c.Name, req.Remark, id)
		}
		writeLog(c, "提现管理", "审核提现", bizNo+" -> "+newStatus)
		okMsg(w, nil, "审核完成："+newStatus)
	}
}

func nullID(v int64) interface{} {
	if v == 0 {
		return nil
	}
	return v
}
