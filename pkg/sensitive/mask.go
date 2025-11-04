package sensitive

import "github.com/showa-93/go-mask"

const (
	Address  = "Address"
	BankCard = "BankCard"
	Name     = "Name"
	UserName = "UserName"
	Email    = "Email"
	IDCard   = "IDCard"
	Phone    = "Phone"
	Password = "Password"
	Common   = "Common"
	All      = "All"
)

const maskChar = "*"

func init() {
	mask.SetMaskChar(maskChar)
	// 地址脱敏：保留前6位
	mask.RegisterMaskStringFunc(Address, func(arg string, value string) (string, error) {
		return maskValueMinMaxMask(value, 6, 0, mask.MaskChar(), 3, len(value)), nil
	})
	// 银行卡脱敏：前6位和后4位
	mask.RegisterMaskStringFunc(BankCard, func(arg string, value string) (string, error) {
		return maskValueMinMaxMask(value, 6, 4, mask.MaskChar(), 6, len(value)), nil
	})
	// 中文姓名：保留姓(前一个字)+名(最后一字)
	mask.RegisterMaskStringFunc(Name, func(arg string, value string) (string, error) {
		return maskValueMinMaxMask(value, 1, 1, mask.MaskChar(), 1, len(value)), nil
	})
	// 用户名：保留前一位，后一位
	mask.RegisterMaskStringFunc(UserName, func(arg string, value string) (string, error) {
		return maskValueMinMaxMask(value, 1, 1, mask.MaskChar(), 3, len(value)), nil
	})
	// 邮箱：仅展示前3个字符及后7个字符
	mask.RegisterMaskStringFunc(Email, func(arg string, value string) (string, error) {
		return maskValueMinMaxMask(value, 3, 7, mask.MaskChar(), 3, len(value)), nil
	})
	// 身份证：显示前4位与后4位
	mask.RegisterMaskStringFunc(IDCard, func(arg string, value string) (string, error) {
		return maskValueMinMaxMask(value, 4, 4, mask.MaskChar(), 6, len(value)), nil
	})
	// 手机号：显示前4位与后4位
	mask.RegisterMaskStringFunc(Phone, func(arg string, value string) (string, error) {
		return maskValueMinMaxMask(value, 4, 4, mask.MaskChar(), 4, len(value)), nil
	})
	// 密码：全不显示
	mask.RegisterMaskStringFunc(Password, func(arg string, value string) (string, error) {
		return maskValueMinMaxMask(value, 0, 0, mask.MaskChar(), 6, 6), nil
	})
	// 通用：只显示前一后一位
	mask.RegisterMaskStringFunc(Common, func(arg string, value string) (string, error) {
		return maskValueMinMaxMask(value, 1, 1, mask.MaskChar(), 1, 10), nil
	})
	// 全不显示：有多少位就显示多位脱敏字符
	mask.RegisterMaskStringFunc(All, func(arg string, value string) (string, error) {
		return maskValueMinMaxMask(value, 0, 0, mask.MaskChar(), 1, len(value)), nil
	})
}

// Mask returns an object with the mask applied to any given object.
// The function's argument can accept any type, including pointer, map, and slice types, in addition to struct.
// from default masker.
func Mask[T any](target T) (ret T, err error) {
	return mask.Mask(target)
}

// String masks the given argument string
// from default masker.
func String(tag, value string) (string, error) {
	return mask.String(tag, value)
}

// Int masks the given argument int
// from default masker.
func Int(tag string, value int) (int, error) {
	return mask.Int(tag, value)
}

// Uint masks the given argument int
// from default masker.
func Uint(tag string, value uint) (uint, error) {
	return mask.Uint(tag, value)
}

// Float64 masks the given argument float64
// from default masker.
func Float64(tag string, value float64) (float64, error) {
	return mask.Float64(tag, value)
}
