package sensitive

import (
	"fmt"
	"testing"
)

func TestMask01(t *testing.T) {
	fmt.Println(String(Address, "广东xxx区xx街道999999号"))
	fmt.Println(String(BankCard, "622333555566668888"))
	fmt.Println(String(Name, "张三丰"))
	fmt.Println(String(Name, "张三"))
	fmt.Println(String(Name, "张三丰生"))
	fmt.Println(String(UserName, "zhangSanFen"))
	fmt.Println(String(Email, "zhangSanFen@xxx.com"))
	fmt.Println(String(Phone, "12011112222"))
	fmt.Println(String(IDCard, "10010030500102100X"))
	fmt.Println(String(Phone, "12011112222"))
	fmt.Println(String(Password, "123456789"))
	fmt.Println(String(Common, "测试字符串abc123"))
	fmt.Println(String(All, "测试字符串abc123"))
	fmt.Println(String("notExist", "12011112222"))
}

type PersonTest struct {
	Address  string `mask:"Address"`
	BankCard string `mask:"BankCard"`
	Name     string `mask:"Name"`
	UserName string `mask:"UserName"`
	Email    string `mask:"Email"`
	IDCard   string `mask:"IDCard"`
	Phone    string `mask:"Phone"`
	Password string `mask:"Password"`
	Age      int    `mask:"random100"`
}

func TestMask02(t *testing.T) {
	var person = PersonTest{
		Address:  "广东xxx区xx街道999999号",
		BankCard: "622333555566668888",
		Name:     "张三丰",
		UserName: "zhangSanFen",
		Email:    "zhangSanFen@xxx.com",
		IDCard:   "10010030500102100X",
		Phone:    "12011112222",
		Password: "123456789",
		Age:      18,
	}
	ret, err := Mask(&person)
	if err != nil {
		panic(err)
	}
	//copier.Copy(&person, ret)
	fmt.Printf("%+v\n", person)
	fmt.Printf("%+v\n", ret)
}
