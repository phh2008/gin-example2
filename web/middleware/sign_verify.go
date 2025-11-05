package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"com.example/example/model/result"
	"com.example/example/pkg/config"
	"com.example/example/pkg/exception"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

const signKey = "sign"
const timestampKey = "timestamp"
const nonceKey = "nonce"

const expireTime int64 = 60

// SignVerify 签名验证
func SignVerify(conf *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 签名验证
		var data map[string]interface{}
		body, err := io.ReadAll(ctx.Request.Body)
		ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
		err = json.Unmarshal(body, &data)
		if err != nil {
			slog.Error("参数解析错误", "error", err)
			result.Error[any](exception.ParamError).Response(ctx)
			ctx.Abort()
			return
		}
		if v, ok := data[nonceKey]; !ok || cast.ToString(v) == "" {
			result.Error[any](errors.New("缺少参数：" + nonceKey)).Response(ctx)
			ctx.Abort()
			return
		}
		if v, ok := data[timestampKey]; !ok || cast.ToString(v) == "" {
			result.Error[any](errors.New("缺少参数：" + timestampKey)).Response(ctx)
			ctx.Abort()
			return
		}
		if _, ok := data[signKey]; !ok {
			result.Error[any](errors.New("缺少参数：" + signKey)).Response(ctx)
			ctx.Abort()
			return
		}
		// 生成签名
		sign := createSign(data, conf.Server.SignToken)
		// 签名是否一致
		if sign != cast.ToString(data[signKey]) {
			result.Error[any](exception.SignVerifyError).Response(ctx)
			ctx.Abort()
			return
		}
		// 时间戳时否过期或超前（暂定60秒内有效）
		timestamp := cast.ToString(data[timestampKey])
		reqTimestamp, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			result.Error[any](exception.SignVerifyError).Response(ctx)
			ctx.Abort()
			return
		}
		sub := time.Now().Unix() - reqTimestamp
		expired := getExpired(conf)
		if sub > expired || -sub > expired {
			result.Error[any](exception.SignVerifyError).Response(ctx)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// getExpired 获取过期时长
func getExpired(conf *config.Config) int64 {
	expire := conf.Server.ExpireTime
	if expire <= 0 {
		return expireTime
	}
	return expire
}

// createSign 生成签名串
func createSign(data map[string]interface{}, signToken string) string {
	var keys []string
	for k, _ := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	for _, k := range keys {
		if k == signKey {
			continue
		}
		val := data[k]
		if val != nil {
			valType := reflect.ValueOf(val)
			kind := valType.Type().Kind()
			if kind == reflect.Struct || kind == reflect.Slice || kind == reflect.Map {
				continue
			}
		}
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(cast.ToString(val))
		sb.WriteString("&")
	}
	values := strings.TrimSuffix(sb.String(), "&") + signToken
	sign := md5.Sum([]byte(values))
	return hex.EncodeToString(sign[:])
}
