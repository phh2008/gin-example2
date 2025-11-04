package sensitive

import "strings"

const MaskLength = 6
const MinMaskLength = 1

func maskValueMinMaxMask(val string, prefixNoMaskLen int, suffixNoMaskLen int, maskStr string, minLength int, maxLength int) string {
	if val == "" {
		return val
	}
	value := []rune(val)
	prefixLength := Max(prefixNoMaskLen, 0)
	suffixLength := Max(suffixNoMaskLen, 0)
	if prefixLength == 0 && suffixLength == 0 {
		// 全遮罩
		return strings.Repeat(maskStr, Max(MaskLength, minLength))
	}
	// 最小脱敏长度
	minLen := Max(minLength, MinMaskLength)
	length := len(value)
	if length <= minLen {
		return maskStr
	}
	// 前后保留不脱敏是否有重合
	total := prefixLength + suffixLength
	if total > length {
		overLength := total - length
		if overLength >= length {
			// 全遮罩
			return strings.Repeat(maskStr, Max(MaskLength, minLength))
		}
		if prefixLength > suffixLength {
			prefixLength = prefixLength - overLength - 1
		} else {
			suffixLength = suffixLength - overLength - 1
		}
	} else if total == length {
		if prefixLength > suffixLength {
			prefixLength--
		} else {
			suffixLength--
		}
	}
	prefix := ""
	if prefixLength > 0 {
		prefix = string(value[0:Min(prefixLength, length)])
	}
	suffix := ""
	if suffixLength != 0 {
		suffix = string(value[Min(Abs(length-suffixLength), length):])
	}
	maskText := strings.Repeat(maskStr, Max(Min(maxLength, length-prefixLength-suffixLength), minLen))
	return prefix + maskText + suffix
}

func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
