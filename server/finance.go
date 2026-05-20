package main

import "net/http"

// 财务汇总: 按园区/平台聚合线上(订单)、美团、抖音收入与提现/冻结/余额
// online  = 已支付/核销订单金额
// meituan = 美团已核销/部分核销券面额
// douyin  = 抖音已核销/部分核销券面额
// withdrawn = 已支付提现; frozen = 待审核提现; balance = 总收益 - 已提现
func parkFinance(parkID int64) map[string]interface{} {
	online := scalarFloat(`SELECT COALESCE(SUM(o.amount),0) FROM "order" o JOIN shop s ON s.id=o.shop_id
		WHERE s.park_id=? AND o.status IN ('已支付','已核销','已完成')`, parkID)
	meituan := scalarFloat(`SELECT COALESCE(SUM(r.sale_price),0) FROM platform_verify_records r
		JOIN shop s ON s.id=r.shop_id WHERE s.park_id=? AND r.platform='meituan' AND r.status IN ('已核销','部分核销')`, parkID)
	douyin := scalarFloat(`SELECT COALESCE(SUM(r.sale_price),0) FROM platform_verify_records r
		JOIN shop s ON s.id=r.shop_id WHERE s.park_id=? AND r.platform='douyin' AND r.status IN ('已核销','部分核销')`, parkID)
	withdrawn := scalarFloat(`SELECT COALESCE(SUM(w.amount),0) FROM withdraw w JOIN shop s ON s.id=w.shop_id
		WHERE s.park_id=? AND w.status='已支付'`, parkID)
	frozen := scalarFloat(`SELECT COALESCE(SUM(w.amount),0) FROM withdraw w JOIN shop s ON s.id=w.shop_id
		WHERE s.park_id=? AND w.status='待审核'`, parkID)
	total := online + meituan + douyin
	balance := total - withdrawn
	if balance < 0 {
		balance = 0
	}

	shops, _ := queryMaps(`SELECT id,name FROM shop WHERE park_id=? ORDER BY id`, parkID)
	shopList := []map[string]interface{}{}
	for _, s := range shops {
		sid := s["id"]
		so := scalarFloat(`SELECT COALESCE(SUM(amount),0) FROM "order" WHERE shop_id=? AND status IN ('已支付','已核销','已完成')`, sid)
		sm := scalarFloat(`SELECT COALESCE(SUM(sale_price),0) FROM platform_verify_records WHERE shop_id=? AND platform='meituan' AND status IN ('已核销','部分核销')`, sid)
		sd := scalarFloat(`SELECT COALESCE(SUM(sale_price),0) FROM platform_verify_records WHERE shop_id=? AND platform='douyin' AND status IN ('已核销','部分核销')`, sid)
		sw := scalarFloat(`SELECT COALESCE(SUM(amount),0) FROM withdraw WHERE shop_id=? AND status='已支付'`, sid)
		sf := scalarFloat(`SELECT COALESCE(SUM(amount),0) FROM withdraw WHERE shop_id=? AND status='待审核'`, sid)
		sbal := so + sm + sd - sw
		if sbal < 0 {
			sbal = 0
		}
		shopList = append(shopList, map[string]interface{}{
			"name": s["name"], "online": so, "offline": 0.0, "meituan": sm, "douyin": sd,
			"balance": sbal, "frozen": sf, "withdrawn": sw,
		})
	}

	return map[string]interface{}{
		"totalRevenue": total, "onlineRevenue": online, "offlineRevenue": 0.0,
		"meituanRevenue": meituan, "douyinRevenue": douyin,
		"balance": balance, "frozen": frozen, "withdrawn": withdrawn,
		"totalIn": total + withdrawn, "shops": shopList,
	}
}

func handleFinanceSummary(w http.ResponseWriter, r *http.Request) {
	sc := scopeOf(getClaims(r))
	where := ""
	args := []interface{}{}
	if !sc.IsPlatform && sc.ParkId > 0 {
		where = " WHERE id=?"
		args = append(args, sc.ParkId)
	}
	parks, _ := queryMaps(`SELECT id,name FROM park`+where+" ORDER BY id", args...)

	parksOut := map[string]interface{}{}
	agg := map[string]float64{}
	add := func(m map[string]interface{}, k string) { agg[k] += m[k].(float64) }
	for _, p := range parks {
		var pid int64
		switch v := p["id"].(type) {
		case int64:
			pid = v
		case float64:
			pid = int64(v)
		}
		f := parkFinance(pid)
		f["name"] = p["name"]
		parksOut[viewKeyForPark(p["name"])] = f
		for _, k := range []string{"totalRevenue", "onlineRevenue", "offlineRevenue", "meituanRevenue", "douyinRevenue", "balance", "frozen", "withdrawn", "totalIn"} {
			add(f, k)
		}
	}
	platform := map[string]interface{}{}
	for k, v := range agg {
		platform[k] = v
	}
	platform["name"] = "平台总部"
	ok(w, map[string]interface{}{"platform": platform, "parks": parksOut})
}

// 园区名 -> 原型视角键
func viewKeyForPark(name interface{}) string {
	s, _ := name.(string)
	switch {
	case contains(s, "黄梅"):
		return "park-hm"
	case contains(s, "武汉"), contains(s, "江夏"):
		return "park-wh"
	default:
		return "park-" + s
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && indexOf(s, sub) >= 0
}
func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
