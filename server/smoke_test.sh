#!/usr/bin/env bash
# 端到端 API 冒烟测试 — 覆盖每个模块 + RBAC + 视角过滤
set -u
BASE="${BASE:-http://localhost:8080}"
PASS=0; FAIL=0
jqv() { python3 -c "import sys,json;d=json.load(sys.stdin);print(d$1)"; }

login() { curl -s -X POST "$BASE/api/auth/login" -H 'Content-Type: application/json' \
  -d "{\"username\":\"$1\",\"password\":\"$2\"}" | jqv "['data']['token']"; }

# check NAME EXPECTED ACTUAL
chk() {
  if [ "$2" = "$3" ]; then PASS=$((PASS+1)); printf "  \033[32m✓\033[0m %s\n" "$1";
  else FAIL=$((FAIL+1)); printf "  \033[31m✗\033[0m %s (期望 %s, 实际 %s)\n" "$1" "$2" "$3"; fi
}
# code of an authed GET
code() { curl -s "$BASE$2" -H "Authorization: Bearer $1" | jqv "['code']"; }
hcode() { curl -s -o /dev/null -w '%{http_code}' "$BASE$2" -H "Authorization: Bearer $1"; }

echo "==== 1. 认证 ===="
ADMIN=$(login admin admin123)
[ -n "$ADMIN" ] && chk "平台超管登录拿到 token" "yes" "yes" || chk "平台超管登录拿到 token" "yes" "no"
BAD=$(curl -s -X POST "$BASE/api/auth/login" -H 'Content-Type: application/json' -d '{"username":"admin","password":"wrong"}' | jqv "['code']")
chk "错误密码被拒" "401" "$BAD"
PARK=$(login park-hm 123456)
SHOP=$(login hcct-mgr 123456)
CASHIER=$(login hcct-cashier 123456)

echo "==== 2. 读取各模块列表 (平台视角) ===="
for ep in "/api/parks" "/api/shops" "/api/shop-types" "/api/users" "/api/users/stats" \
  "/api/recharges" "/api/orders" "/api/refunds" "/api/withdraws" "/api/products" \
  "/api/employees" "/api/positions" "/api/gates" "/api/faces" "/api/entry-records" \
  "/api/activities" "/api/platform/verify-records?platform=meituan" \
  "/api/platform/verify-records?platform=douyin" "/api/platform/settlements?platform=meituan" \
  "/api/platform/store-configs?platform=douyin" "/api/system/mini-programs" \
  "/api/system/logs" "/api/permissions"; do
  chk "GET $ep" "0" "$(code "$ADMIN" "$ep")"
done

echo "==== 3. RBAC 权限 ===="
# 平台超管对 shops 只读 -> 创建店铺应 403
C=$(curl -s -X POST "$BASE/api/shops" -H "Authorization: Bearer $ADMIN" -H 'Content-Type: application/json' -d '{"name":"x","parkId":1}' | jqv "['code']")
chk "平台超管(只读)创建店铺被拒" "403" "$C"
# 园区管理员对 shops 可操作 -> 创建成功
C=$(curl -s -X POST "$BASE/api/shops" -H "Authorization: Bearer $PARK" -H 'Content-Type: application/json' -d '{"name":"测试店铺","parkId":1,"typeId":1}' | jqv "['code']")
chk "园区管理员创建店铺成功" "0" "$C"
# 收银员不可提现审核 (finance 隐藏) -> 403
C=$(hcode "$CASHIER" "/api/withdraws")
chk "收银员访问提现列表被拒(HTTP403)" "403" "$C"
# 商户管理员可访问提现
chk "商户管理员可访问提现列表" "0" "$(code "$SHOP" "/api/withdraws")"
# 未登录 -> 401
chk "无 token 访问被拒(HTTP401)" "401" "$(curl -s -o /dev/null -w '%{http_code}' "$BASE/api/parks")"

echo "==== 4. 视角数据过滤 (订单模块) ===="
PLAT_ORD=$(curl -s "$BASE/api/orders" -H "Authorization: Bearer $ADMIN" | jqv "['data']['total']")
PARK_ORD=$(curl -s "$BASE/api/orders" -H "Authorization: Bearer $PARK" | jqv "['data']['total']")
SHOP_ORD=$(curl -s "$BASE/api/orders" -H "Authorization: Bearer $SHOP" | jqv "['data']['total']")
chk "平台视角订单总数(=5)" "5" "$PLAT_ORD"
chk "园区视角仅本园区订单(=5)" "5" "$PARK_ORD"
chk "商户视角仅本店铺订单(=2)" "2" "$SHOP_ORD"

echo "==== 5. 业务流程: 核销 / 退款 / 提现审核 ===="
# 核销 已支付订单
C=$(curl -s -X POST "$BASE/api/orders/202605141422001/verify" -H "Authorization: Bearer $SHOP" | jqv "['code']")
chk "核销已支付订单成功" "0" "$C"
# 重复核销失败
C=$(curl -s -X POST "$BASE/api/orders/202605141422001/verify" -H "Authorization: Bearer $SHOP" | jqv "['code']")
chk "重复核销被拒" "400" "$C"
# 退款审核(通过)
C=$(curl -s -X POST "$BASE/api/refunds/RF20260514114009/approve" -H "Authorization: Bearer $PARK" -H 'Content-Type: application/json' -d '{"remark":"同意"}' | jqv "['code']")
chk "退款审核通过" "0" "$C"
# 提现审核(通过)
C=$(curl -s -X POST "$BASE/api/withdraws/TX26051415070/approve" -H "Authorization: Bearer $PARK" -H 'Content-Type: application/json' -d '{"remark":"打款"}' | jqv "['code']")
chk "提现审核通过" "0" "$C"

echo "==== 6. 第三方验券流程 ===="
C=$(curl -s -X POST "$BASE/api/platform/verify/prepare" -H "Authorization: Bearer $SHOP" -H 'Content-Type: application/json' -d '{"platform":"meituan","couponCode":"MT-NEW-0001","shopId":1}' | jqv "['code']")
chk "美团验券准备(已绑定门店)" "0" "$C"
C=$(curl -s -X POST "$BASE/api/platform/verify/execute" -H "Authorization: Bearer $SHOP" -H 'Content-Type: application/json' -d '{"platform":"meituan","couponCode":"MT-NEW-0001","shopId":1,"product":"测试套餐","salePrice":58}' | jqv "['code']")
chk "美团执行核销" "0" "$C"
C=$(curl -s -X POST "$BASE/api/platform/verify/prepare" -H "Authorization: Bearer $SHOP" -H 'Content-Type: application/json' -d '{"platform":"meituan","couponCode":"X","shopId":3}' | jqv "['code']")
chk "未绑定门店验券被拒" "400" "$C"

echo "==== 7. 创建类操作 ===="
C=$(curl -s -X POST "$BASE/api/parks" -H "Authorization: Bearer $ADMIN" -H 'Content-Type: application/json' -d '{"name":"测试园区","code":"PARK_TEST_'"$RANDOM"'"}' | jqv "['code']")
chk "平台超管新增园区" "0" "$C"
C=$(curl -s -X POST "$BASE/api/users" -H "Authorization: Bearer $PARK" -H 'Content-Type: application/json' -d '{"nickname":"新客户","phone":"13000000000"}' | jqv "['code']")
chk "新增客户" "0" "$C"
C=$(curl -s -X POST "$BASE/api/activities" -H "Authorization: Bearer $PARK" -H 'Content-Type: application/json' -d '{"title":"测试活动","type":"满减活动"}' | jqv "['code']")
chk "新增活动" "0" "$C"

echo "==== 8. 财务汇总 ===="
chk "财务汇总(平台,只读可见)" "0" "$(code "$ADMIN" "/api/finance/summary")"
chk "财务汇总(园区)" "0" "$(code "$PARK" "/api/finance/summary")"
HASPLAT=$(curl -s "$BASE/api/finance/summary" -H "Authorization: Bearer $ADMIN" | jqv "['data']['platform']['totalRevenue']" 2>/dev/null)
[ -n "$HASPLAT" ] && chk "汇总含平台总收益字段" "yes" "yes" || chk "汇总含平台总收益字段" "yes" "no"

echo "==== 9. 闸机写操作 (园区视角可操作) ===="
# 平台超管 gates 只读 -> 创建应 403
chk "平台超管(只读)创建闸机被拒" "403" "$(curl -s -X POST "$BASE/api/gates" -H "Authorization: Bearer $ADMIN" -H 'Content-Type: application/json' -d '{"name":"x","deviceSn":"SN-X-1"}' | jqv "['code']")"
GID=$(curl -s -X POST "$BASE/api/gates" -H "Authorization: Bearer $PARK" -H 'Content-Type: application/json' -d '{"name":"测试闸机","deviceSn":"SN-TEST-'"$RANDOM"'","type":"扫码闸机","direction":"进"}' | jqv "['data']['id']")
[ -n "$GID" ] && chk "园区管理员新增闸机" "yes" "yes" || chk "园区管理员新增闸机" "yes" "no"
chk "切换闸机启停" "0" "$(curl -s -X PUT "$BASE/api/gates/$GID/toggle" -H "Authorization: Bearer $PARK" | jqv "['code']")"
chk "删除闸机" "0" "$(curl -s -X DELETE "$BASE/api/gates/$GID" -H "Authorization: Bearer $PARK" | jqv "['code']")"

echo "==== 10. 活动审核 / 结束 ===="
# 找一个待审核活动 (种子: 火车餐厅6月限定)
AID=$(curl -s "$BASE/api/activities?status=待审核" -H "Authorization: Bearer $ADMIN" | python3 -c "import sys,json;l=json.load(sys.stdin)['data']['list'];print(l[0]['id'] if l else '')")
chk "存在待审核活动" "yes" "$([ -n "$AID" ] && echo yes || echo no)"
chk "审核通过活动" "0" "$(curl -s -X POST "$BASE/api/activities/$AID/audit" -H "Authorization: Bearer $ADMIN" -H 'Content-Type: application/json' -d '{"approve":true}' | jqv "['code']")"
chk "重复审核被拒" "400" "$(curl -s -X POST "$BASE/api/activities/$AID/audit" -H "Authorization: Bearer $ADMIN" -H 'Content-Type: application/json' -d '{"approve":true}' | jqv "['code']")"
# 结束一个进行中活动
EID=$(curl -s "$BASE/api/activities?status=进行中" -H "Authorization: Bearer $ADMIN" | python3 -c "import sys,json;l=json.load(sys.stdin)['data']['list'];print(l[0]['id'] if l else '')")
chk "结束进行中活动" "0" "$(curl -s -X POST "$BASE/api/activities/$EID/end" -H "Authorization: Bearer $ADMIN" | jqv "['code']")"

echo ""
echo "================================"
printf "通过: \033[32m%d\033[0m  失败: \033[31m%d\033[0m\n" "$PASS" "$FAIL"
echo "================================"
[ "$FAIL" -eq 0 ]
