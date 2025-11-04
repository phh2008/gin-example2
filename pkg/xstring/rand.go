package xstring

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var src = rand.NewSource(time.Now().UnixNano())
var randCommon = rand.New(src)

// RandCode 生成随机验证码数字
func RandCode() int {
	min := 100000
	max := 999999
	return randCommon.Intn(max-min+1) + min
}

// RandCodeStr 生成随机验证码字符串
func RandCodeStr() string {
	return strconv.Itoa(RandCode())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandStr 随机字符串
func RandStr(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randCommon.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randCommon.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// UUID 生成UUID
func UUID() string {
	u, _ := uuid.NewUUID()
	return u.String()
}
