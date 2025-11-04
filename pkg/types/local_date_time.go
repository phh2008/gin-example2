package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// LocalDateTime 包装time，用于数据库查询与写入，对应 mysql 的 datetime 类型
type LocalDateTime struct {
	time.Time
}

func (t *LocalDateTime) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	req := strings.Trim(string(data), "\"")
	if req == "" {
		return nil
	}
	date, err := time.Parse(time.DateTime, req)
	if err != nil {
		return err
	}
	t.Time = date
	return nil
}

// MarshalJSON 重写 MarshaJSON 方法，在此方法中实现自定义格式的转换；
func (t LocalDateTime) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", t.Format(time.DateTime))
	return []byte(output), nil
}

// Value 实现 Value 方法，写入数据库时会调用该方法将自定义时间类型转换并写入数据库；
func (t LocalDateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan 实现 Scan 方法，读取数据库时会调用该方法将时间数据转换成自定义时间类型；
func (t *LocalDateTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalDateTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to LocalDateTime", v)
}
