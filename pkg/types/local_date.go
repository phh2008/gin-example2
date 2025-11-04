package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// LocalDate 包装time，用于数据库查询与写入，对应 mysql 的 date 类型
type LocalDate struct {
	time.Time
}

func (t *LocalDate) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	req := strings.Trim(string(data), "\"")
	if req == "" {
		return nil
	}
	date, err := time.Parse(time.DateOnly, req)
	if err != nil {
		return err
	}
	t.Time = date
	return nil
}

// MarshalJSON 重写 MarshaJSON 方法，在此方法中实现自定义格式的转换；
func (t LocalDate) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", t.Format(time.DateOnly))
	return []byte(output), nil
}

// Value 实现 Value 方法，写入数据库时会调用该方法将自定义时间类型转换并写入数据库；
func (t LocalDate) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan 实现 Scan 方法，读取数据库时会调用该方法将时间数据转换成自定义时间类型；
func (t *LocalDate) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalDate{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to LocalDate", v)
}
