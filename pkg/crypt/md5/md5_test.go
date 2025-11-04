package md5

import (
	"encoding/base64"
	"testing"
)

func TestMd5(t *testing.T) {
	src := "abc"
	dst := Md5(src)
	contentMD5 := base64.StdEncoding.EncodeToString(dst)
	t.Log(contentMD5)
	t.Log(Md5ToString(src))
}
