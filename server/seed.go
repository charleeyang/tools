package main

import "log"

// 模块可见性映射 (取自原型 MENU_VISIBILITY)，按角色落库到 role_permission
// 值: 1 -> 可操作, 2 -> 只读
var permPlatform = map[string]int{
	"overview": 1, "databoard": 1, "sales-stats": 1, "flow-stats": 1,
	"parks": 1, "shops": 2, "products": 1, "users": 2, "consumption": 2,
	"meituan": 1, "douyin": 1, "finance": 2, "employees": 1,
	"gates": 2, "activities": 1, "permissions": 1, "mp-mgmt": 1, "settings": 1,
}
var permPark = map[string]int{
	"overview": 1, "databoard": 1, "sales-stats": 1, "flow-stats": 1,
	"shops": 1, "products": 1, "users": 1, "consumption": 1, "meituan": 1, "douyin": 1,
	"finance": 1, "gates": 1, "activities": 1, "mp-mgmt": 1, "employees": 2,
}
var permShopAdmin = map[string]int{
	"consumption": 1, "verify": 1, "products": 1, "meituan": 1, "douyin": 1, "finance": 1, "withdraw": 1,
}
var permShopCashier = map[string]int{
	"consumption": 1, "verify": 1, "meituan": 1, "douyin": 1,
}

func visLabel(v int) string {
	if v == 2 {
		return "只读"
	}
	return "可操作"
}

func seedPerms(role string, m map[string]int) {
	for k, v := range m {
		mustExec(`INSERT INTO role_permission(role_code,module_key,visibility) VALUES(?,?,?)`,
			role, k, visLabel(v))
	}
}

func mustExec(q string, args ...interface{}) {
	if _, err := db.Exec(q, args...); err != nil {
		log.Fatalf("种子写入失败 [%s]: %v", q, err)
	}
}

func seedAll() {
	// ---- 角色权限 ----
	seedPerms("ADMIN_PLATFORM", permPlatform)
	seedPerms("PLATFORM_OPER", permPlatform)
	seedPerms("PARK_ADMIN", permPark)
	seedPerms("PARK_MANAGER", permPark)
	seedPerms("SHOP_ADMIN", permShopAdmin)
	seedPerms("SHOP_CASHIER", permShopCashier)

	// ---- 岗位 ----
	mustExec(`INSERT INTO position(name,code,responsibilities) VALUES
		('平台运营','platform_oper','平台整体运营'),
		('园区经理','park_manager','园区日常管理'),
		('火车餐厅经理','hcct_manager','门店经营管理'),
		('前厅主管','hall_supervisor','前厅收银核销'),
		('咖啡店长','coffee_manager','咖啡店经营')`)

	// ---- 店铺类型 ----
	mustExec(`INSERT INTO shop_type(name,code,description) VALUES
		('餐饮','catering','餐厅、咖啡厅、小吃等'),
		('零售','retail','特产、纪念品、日用品销售'),
		('体验','experience','手作、农耕体验等活动'),
		('其他','other','其他类型商户')`)

	// ---- 园区 ----
	mustExec(`INSERT INTO park(name,code,address,contact_person,contact_phone,status,created_at) VALUES
		('黄梅袁夫稻田','PARK_HM_001','湖北省黄冈市黄梅县大河镇','李四','13700137000','营业中','2026-01-21 16:37'),
		('武汉袁夫稻田','PARK_WH_001','湖北省武汉市江夏区','周八','13600136000','营业中','2026-03-15 09:20')`)

	// ---- 店铺 (type: 1餐饮 2零售 3体验) ----
	mustExec(`INSERT INTO shop(park_id,name,address,type_id,status,created_at) VALUES
		(1,'火车餐厅','湖北省黄冈市黄梅县大河镇77号',1,'营业中','2026-03-10 17:48'),
		(1,'树下咖啡','湖北省黄冈市黄梅县大河镇78号',1,'营业中','2026-04-09 13:40'),
		(1,'稻田手作坊','湖北省黄冈市黄梅县大河镇79号',3,'休息中','2026-05-01 09:20'),
		(2,'稻田鲜货铺','湖北省武汉市江夏区21号',2,'营业中','2026-04-22 11:30'),
		(2,'江夏火车厨','湖北省武汉市江夏区23号',1,'营业中','2026-04-25 09:00')`)

	// ---- 客户 ----
	mustExec(`INSERT INTO user(id,nickname,phone,openid,member_level,balance,consumption_total,points,face_status,registered_at) VALUES
		(3005,'稻田守望者','187****9908','oXG-2x...JK7v','VIP2',130.00,256.50,256,'未录入','2026-05-13 14:22'),
		(3004,'田园生活家','138****5678','oXG-9P...zM3a','VIP1',50.00,89.00,89,'已录入','2026-04-20 10:15'),
		(3003,'稻香一缕','159****2345','oXG-7t...Qm4b','普通用户',0.00,22.00,22,'未录入','2026-04-15 08:30'),
		(3002,'小麦的麦','152****6677','oXG-3a...Yp9d','VIP3',580.00,942.30,942,'已录入','2026-03-22 09:50'),
		(3001,'游客小新','199****1122','oXG-5f...Wq2e','普通用户',10.00,88.00,88,'已录入','2026-02-09 16:57')`)

	// ---- 商品 ----
	mustExec(`INSERT INTO product(shop_id,name,price,original_price,stock,sales_count,status) VALUES
		(1,'稻田双人套餐',128.00,198.00,50,120,'已上架'),
		(1,'火车特色套餐',99.00,168.00,80,86,'已上架'),
		(2,'树下手冲咖啡',39.90,68.00,200,210,'已上架'),
		(3,'稻田米浆体验',19.90,38.00,100,42,'已上架'),
		(4,'袁夫稻香米5kg',59.00,79.00,300,156,'已上架')`)

	// ---- 订单 (user_id, shop_id 对应上方) ----
	mustExec(`INSERT INTO "order"(order_no,user_id,shop_id,amount,pay_method,status,source,created_at) VALUES
		('202605141422001',3005,1,86.50,'余额','已支付','self','2026-05-14 14:22'),
		('202605141315002',3004,2,32.00,'微信支付','已核销','self','2026-05-14 13:15'),
		('202605141140003',3001,3,22.00,'余额','退款处理中','self','2026-05-14 11:40'),
		('202605141025004',3002,1,128.00,'余额','已核销','self','2026-05-14 10:25'),
		('202605140915005',3003,2,46.00,'微信支付','已退款','self','2026-05-14 09:15')`)

	// ---- 充值记录 ----
	mustExec(`INSERT INTO recharge(biz_no,user_id,amount,gift_amount,pay_method,txn_no,remark,status,created_at) VALUES
		('CZ20260514102230',3005,200,30,'微信支付','4200001234202605141022','','成功','2026-05-14 10:22'),
		('TF20260513163019',3004,100,10,'微信支付','4200001234202605131630','','成功','2026-05-13 16:30'),
		('CZ20260513102211',3002,500,100,'微信支付','4200001234202605131022','','成功','2026-05-13 10:22'),
		('TK20260512091215',3005,-50,0,'微信支付','TK20260512091215001','用户申请余额退款','已退款','2026-05-12 09:12'),
		('CZ20260511183055',3001,50,0,'微信支付','4200001234202605111830','','成功','2026-05-11 18:30'),
		('CZ20260510091523',3005,100,10,'微信支付','4200001234202605100915','','成功','2026-05-10 09:15')`)

	// ---- 退款申请 ----
	mustExec(`INSERT INTO refund(biz_no,user_id,order_no,amount,type,reason,status,applied_at) VALUES
		('RF20260514114009',3001,'202605141140003',22.00,'订单退款','体验项目临时取消','待审核','2026-05-14 12:00'),
		('RF20260512091200',3005,'-',50.00,'余额退款','账户停用申请全额提现','已通过','2026-05-12 09:12'),
		('RF20260510143055',3002,'202605101040002',38.00,'订单退款','菜品口味不符','已驳回','2026-05-10 14:30'),
		('RF20260509164220',3004,'-',120.00,'余额退款','用户主动申请','已通过','2026-05-09 16:42')`)

	// ---- 提现申请 (shop_id) ----
	mustExec(`INSERT INTO withdraw(biz_no,shop_id,amount,bank_name,bank_account,bank_holder,status,applied_at,reconciliation_id,reviewer,reviewed_at,paid_at) VALUES
		('TX26051415070',1,1200.00,'招商银行','6214****6908','黄轩辉','待审核','2026-05-14 15:07','',NULL,NULL,NULL),
		('TX26051310220',2,800.00,'建设银行','6217****3320','孙七','待审核','2026-05-13 10:22','',NULL,NULL,NULL),
		('TX26051209150',1,2400.00,'招商银行','6214****6908','黄轩辉','已支付','2026-05-12 09:15','DZ20260512-FCCT','李四','2026-05-12 09:30','2026-05-12 09:35'),
		('TX26050816220',2,1500.00,'建设银行','6217****3320','孙七','已支付','2026-05-08 16:20','DZ20260508-SXKF','李四','2026-05-08 16:42','2026-05-08 16:50'),
		('TX26050714050',3,320.00,'农业银行','6228****1190','赵六','已驳回','2026-05-07 14:05','','李四','2026-05-07 14:30',NULL)`)

	// ---- 员工 (登录账号) ----
	pwAdmin := hashPassword("admin123")
	pw := hashPassword("123456")
	mustExec(`INSERT INTO employee(username,name,phone,park_id,shop_id,position_id,role_code,password_hash,status,online_status,created_at) VALUES
		('admin','张三','13800138000',NULL,NULL,1,'ADMIN_PLATFORM',?, '启用',1,'2026-01-05 15:09'),
		('park-hm','李四','13700137000',1,NULL,2,'PARK_ADMIN',?, '启用',1,'2026-02-09 16:57'),
		('hcct-mgr','王五','13912345678',1,1,3,'SHOP_ADMIN',?, '启用',1,'2026-03-11 14:38'),
		('hcct-cashier','赵六','13912345699',1,1,4,'SHOP_CASHIER',?, '启用',0,'2026-03-11 14:49'),
		('shuxia-mgr','孙七','18888888888',1,2,5,'SHOP_ADMIN',?, '启用',0,'2026-04-09 13:41'),
		('park-wh','周八','13600136000',2,NULL,2,'PARK_ADMIN',?, '启用',1,'2026-04-22 10:20')`,
		pwAdmin, pw, pw, pw, pw, pw)

	// ---- 闸机 (对齐原型展示) ----
	mustExec(`INSERT INTO gate(park_id,name,type,direction,device_sn,enabled,online_status,today_pass) VALUES
		(1,'主入口 1#','人脸闸机','进','E014370C2321',1,1,78),
		(1,'主入口 2#','人脸闸机','出','E014370C2322',1,1,50),
		(1,'侧门 3#','扫码闸机','进','E014370C2323',0,0,0)`)

	// ---- 人脸 ----
	mustExec(`INSERT INTO face(user_id,feature_id,authorized_gates,status) VALUES
		(3004,'FACE-3004-A1B2','黄梅园区全部闸机','已录入'),
		(3002,'FACE-3002-C3D4','黄梅园区全部闸机','已录入'),
		(3001,'FACE-3001-E5F6','武汉园区正门','已录入')`)

	// ---- 进出记录 ----
	mustExec(`INSERT INTO entry_record(user_id,gate_id,direction,method,passed_at) VALUES
		(3004,1,'进','人脸','2026-05-14 09:12'),
		(3002,1,'进','人脸','2026-05-14 10:01'),
		(3004,2,'出','人脸','2026-05-14 15:48'),
		(3001,1,'进','扫码','2026-05-14 11:20')`)

	// ---- 活动 ----
	mustExec(`INSERT INTO activity(title,type,park_scope,shop_scope,start_at,end_at,status,exposure,joins,orders,verify_rate) VALUES
		('🌾 夏日稻田音乐节','节日活动','全部园区','全部商户','2026-06-01','2026-06-03','进行中',1245,86,62,'72%'),
		('🎫 充值送好礼','会员专享','全部园区','全部商户','2026-05-01','2026-12-31','进行中',3580,420,380,'90%'),
		('☕ 树下咖啡满减券','满减活动','黄梅袁夫稻田','树下咖啡','2026-05-10','2026-05-25','进行中',820,160,130,'81%'),
		('🎁 五一假期满100减30','限时折扣','全部园区','全部商户','2026-05-01','2026-05-05','已结束',5320,730,680,'93%'),
		('🚂 火车餐厅 6 月限定','优惠券活动','黄梅袁夫稻田','火车餐厅','2026-06-05','2026-06-30','待审核',0,0,0,'-')`)

	// ---- 美团核销记录 ----
	mustExec(`INSERT INTO platform_verify_records(biz_no,platform,coupon_code,product,origin_price,sale_price,user_name,phone,shop_id,shop_name,coupon_type,verify_count,max_count,status,verify_by,verify_at,refundable) VALUES
		('MT20260514142201','meituan','MT-YF-8821-5023','稻田双人套餐',198.00,128.00,'稻田守望者','187****9908',1,'火车餐厅','团购套餐',1,1,'已核销','王五','2026-05-14 14:22',1),
		('MT20260514131502','meituan','MT-YF-9942-7163','树下咖啡体验券',68.00,39.90,'田园生活家','138****5678',2,'树下咖啡','团购套餐',1,1,'已核销','赵六','2026-05-14 13:15',1),
		('MT20260514102503','meituan','MT-YF-7710-5532','火车特色套餐',168.00,99.00,'小麦的麦','152****6677',1,'火车餐厅','团购套餐',1,1,'已撤销','王五','2026-05-14 10:25',0),
		('MT20260514091504','meituan','MT-YF-3308-1247','稻田米浆体验',38.00,19.90,'稻香一缕','159****2345',3,'稻田手作坊','团购套餐',1,1,'已核销','赵六','2026-05-14 09:15',1),
		('MT20260513203005','meituan','MT-YF-6621-8890','稻田双人套餐',198.00,128.00,'游客小新','199****1122',1,'火车餐厅','团购套餐',1,1,'已核销','王五','2026-05-13 20:30',0)`)

	// ---- 抖音核销记录 ----
	mustExec(`INSERT INTO platform_verify_records(biz_no,platform,coupon_code,product,origin_price,sale_price,user_name,phone,shop_id,shop_name,coupon_type,verify_count,max_count,status,verify_by,verify_at,refundable) VALUES
		('DY20260514182001','douyin','DY-GROUP-A8K2-M9P1','稻田丰收双人餐',228.00,158.00,'稻田守望者','187****9908',1,'火车餐厅','团购套餐',1,1,'已核销','王五','2026-05-14 18:20',1),
		('DY20260514163002','douyin','DY-GROUP-C3F7-R2X5','精品手冲咖啡体验',88.00,49.90,'田园生活家','138****5678',2,'树下咖啡','团购套餐',1,1,'已核销','赵六','2026-05-14 16:30',1),
		('DY20260514150003','douyin','DY-CARD-T9W4-L8H6','稻田米浆次卡 (10次)',380.00,199.00,'小麦的麦','152****6677',1,'火车餐厅','次卡',4,10,'部分核销','王五','2026-05-14 15:00',0),
		('DY20260514142004','douyin','DY-VOUCHER-B1N6-S3D8','满100减20代金券',20.00,9.90,'稻香一缕','159****2345',2,'树下咖啡','代金券',1,1,'已核销','赵六','2026-05-14 14:20',1),
		('DY20260514105005','douyin','DY-GROUP-J5M2-K7Q9','火车特色双人餐',298.00,188.00,'游客小新','199****1122',1,'火车餐厅','团购套餐',1,1,'已撤销','王五','2026-05-14 10:50',0),
		('DY20260513203006','douyin','DY-PACK-P4R9-H2V1','稻田一日游组合券包',456.00,268.00,'稻田守望者','187****9908',1,'火车餐厅','组合券包',2,3,'部分核销','王五','2026-05-13 20:30',1)`)

	// ---- 结算明细 ----
	mustExec(`INSERT INTO platform_settlements(biz_no,platform,shop_name,settle_date,order_count,total_amount,commission,net_amount,status,period,rate) VALUES
		('MT-STL-20260514-001','meituan','火车餐厅','2026-05-14',12,1420.50,142.05,1278.45,'待结算','T+1','10%'),
		('MT-STL-20260513-002','meituan','火车餐厅','2026-05-13',8,896.00,89.60,806.40,'已结算','T+1','10%'),
		('MT-STL-20260514-003','meituan','树下咖啡','2026-05-14',5,218.50,21.85,196.65,'待结算','T+1','10%'),
		('DY-STL-20260514-001','douyin','火车餐厅','2026-05-14',15,2280.50,114.03,2166.47,'待结算','T+1','5%'),
		('DY-STL-20260513-002','douyin','火车餐厅','2026-05-13',10,1420.00,71.00,1349.00,'已结算','T+1','5%'),
		('DY-STL-20260514-003','douyin','树下咖啡','2026-05-14',6,299.40,14.97,284.43,'待结算','T+1','5%')`)

	// ---- 平台门店配置 ----
	mustExec(`INSERT INTO platform_store_configs(platform,shop_id,shop_name,app_auth_token,status,bind_at,last_sync,verify_count,scope) VALUES
		('meituan',1,'火车餐厅','MT-AUTH-HCCT-xxxx-xxxx-xxxx','已绑定','2026-05-10 09:30','2026-05-14 14:25',42,'门店映射+商品管理'),
		('meituan',2,'树下咖啡','MT-AUTH-SXKF-xxxx-xxxx-xxxx','已绑定','2026-05-11 14:20','2026-05-14 13:18',18,'门店映射'),
		('meituan',3,'稻田手作坊','','待绑定',NULL,NULL,0,'-'),
		('douyin',1,'火车餐厅','DY-AUTH-HCCT-xxxx-xxxx-xxxx','已绑定','2026-05-12 10:15','2026-05-14 18:22',58,'门店映射+商品管理'),
		('douyin',2,'树下咖啡','DY-AUTH-SXKF-xxxx-xxxx-xxxx','已绑定','2026-05-13 11:40','2026-05-14 16:35',22,'门店映射'),
		('douyin',3,'稻田手作坊','','待绑定',NULL,NULL,0,'-')`)

	// ---- 小程序 ----
	mustExec(`INSERT INTO mini_program(app_id,name,type,type_label,app_secret,callback_url,mch_id,status,version,published_at) VALUES
		('wx7c2a3b4e5d6f7890','袁夫稻田·客户端','client','客户端','••••••••••••••••','https://api.yuanfu-rice.com/wx/client/callback','1900012345','active','v2.3.1','2026-05-12'),
		('wx8d3b4c5e6f7a8901','袁夫稻田·商户端','merchant','商户端','••••••••••••••••','https://api.yuanfu-rice.com/wx/merchant/callback','1900012345','active','v1.8.4','2026-05-10'),
		('wx9e4c5d6f7a8b9012','袁夫稻田·管理端','admin','管理端','••••••••••••••••','https://api.yuanfu-rice.com/wx/admin/callback','1900012345','active','v1.5.2','2026-05-08')`)

	// ---- 优惠券 ----
	mustExec(`INSERT INTO coupon(activity_id,name,type,discount_amount,threshold,total_count,used_count,valid_from,valid_until,status) VALUES
		(3,'树下咖啡满50减10','满减券',10,50,500,160,'2026-05-10','2026-05-25','启用'),
		(4,'五一满100减30','满减券',30,100,1000,680,'2026-05-01','2026-05-05','已结束')`)

	// ---- 操作日志 ----
	mustExec(`INSERT INTO operation_log(operator,role_code,module,action,detail,ip,created_at) VALUES
		('张三','ADMIN_PLATFORM','员工管理','新增员工','新增 周八(park-wh)','10.0.0.12','2026-05-14 10:20'),
		('李四','PARK_ADMIN','提现管理','审核提现','驳回 TX26050714050','10.0.0.21','2026-05-07 14:30'),
		('王五','SHOP_ADMIN','美团核销','撤销验券','撤销 MT-YF-7710-5532','10.0.0.33','2026-05-14 10:30')`)
}
