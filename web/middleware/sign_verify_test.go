package middleware

import (
	"fmt"
	"testing"
	"time"

	"com.example/example/pkg/config"
)

// TestCreateSign 测试生成签名
func TestCreateSign(t *testing.T) {
	conf := config.NewConfig("../../config")
	var data = map[string]interface{}{
		signKey:      "",
		anonKey:      "abc",
		timestampKey: time.Now().Unix(),
		"openId":     "oylxy5HqJ630VVOw1CR3diqUzfuQ",
	}
	fmt.Println(data)
	sign := createSign(data, conf.Server.SignToken)
	fmt.Println(sign)
}
