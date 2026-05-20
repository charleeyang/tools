package main

import "strings"

// queryMaps 执行查询并返回 []map，列名转为 camelCase，方便前端直接消费
func queryMaps(q string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	out := []map[string]interface{}{}
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		m := make(map[string]interface{}, len(cols))
		for i, c := range cols {
			v := vals[i]
			if b, ok := v.([]byte); ok {
				v = string(b)
			}
			m[toCamel(c)] = v
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func scalarInt(q string, args ...interface{}) int64 {
	var n int64
	_ = db.QueryRow(q, args...).Scan(&n)
	return n
}

func scalarFloat(q string, args ...interface{}) float64 {
	var f float64
	_ = db.QueryRow(q, args...).Scan(&f)
	return f
}

func toCamel(s string) string {
	if !strings.Contains(s, "_") {
		return s
	}
	parts := strings.Split(s, "_")
	var b strings.Builder
	for i, p := range parts {
		if p == "" {
			continue
		}
		if i == 0 {
			b.WriteString(p)
		} else {
			b.WriteString(strings.ToUpper(p[:1]) + p[1:])
		}
	}
	return b.String()
}

func writeLog(c *employeeClaims, module, action, detail string) {
	if c == nil {
		return
	}
	_, _ = db.Exec(`INSERT INTO operation_log(operator,role_code,module,action,detail,ip) VALUES(?,?,?,?,?,?)`,
		c.Name, c.RoleCode, module, action, detail, "")
}
