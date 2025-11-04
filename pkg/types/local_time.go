package types

import (
	"database/sql/driver"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// LocalTime 包装time，用于数据库查询与写入，对应 mysql 的 time 类型
type LocalTime struct {
	time.Time
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	req := strings.Trim(string(data), "\"")
	if req == "" {
		return nil
	}
	date, err := time.Parse(time.TimeOnly, req)
	if err != nil {
		return err
	}
	t.Time = date
	return nil
}

// MarshalJSON 重写 MarshaJSON 方法，在此方法中实现自定义格式的转换；
func (t LocalTime) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", t.Format(time.TimeOnly))
	return []byte(output), nil
}

// Value 实现 Value 方法，写入数据库时会调用该方法将自定义时间类型转换并写入数据库；
func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	if t.Year() == 0 {
		// 解决只有时间年份为0时插入报错问题：year is not in the range [1, 9999]: 0
		tm := t.Time.AddDate(1, 0, 0)
		// 如果 00:00:00 ,会插入失败
		if tm.Hour() == 0 && tm.Minute() == 0 && tm.Second() == 0 {
			return "00:00:00", nil
		}
		return tm, nil
	}
	return t.Time, nil
}

// Scan 实现 Scan 方法，读取数据库时会调用该方法将时间数据转换成自定义时间类型；
func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.([]uint8); ok {
		// 时间格式
		timeStr := string(value)
		tm, err := time.Parse(time.TimeOnly, timeStr)
		if err != nil {
			slog.Error("mysql time 转换为 time.Time 错误", "error", err)
			return fmt.Errorf("can not convert %v to LocalTime", v)
		}
		*t = LocalTime{Time: tm}
		return nil
	}
	return fmt.Errorf("can not convert %v to LocalTime", v)
}
