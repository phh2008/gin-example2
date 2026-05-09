package cache

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestMemeryCache_GetSet(t *testing.T) {
	c := NewMemeryCache()

	// 测试字符串类型
	t.Run("string", func(t *testing.T) {
		key := "test:str"
		err := Set(c, key, "hello", time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
		val, err := Get[string](c, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if val != "hello" {
			t.Fatalf("expected hello, got %s", val)
		}
	})

	// 测试结构体类型
	t.Run("struct", func(t *testing.T) {
		key := "test:struct"
		want := User{Name: "alice", Age: 20}
		err := Set(c, key, want, time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
		val, err := Get[User](c, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if val != want {
			t.Fatalf("expected %+v, got %+v", want, val)
		}
	})

	// 测试 int 类型
	t.Run("int", func(t *testing.T) {
		key := "test:int"
		err := Set(c, key, 42, time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
		val, err := Get[int](c, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if val != 42 {
			t.Fatalf("expected 42, got %d", val)
		}
	})

	// 测试 float64 类型
	t.Run("float64", func(t *testing.T) {
		key := "test:float"
		err := Set(c, key, 3.14, time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
		val, err := Get[float64](c, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if val != 3.14 {
			t.Fatalf("expected 3.14, got %f", val)
		}
	})

	// 测试 byte 类型
	t.Run("byte", func(t *testing.T) {
		key := "test:byte"
		err := Set(c, key, byte('A'), time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
		val, err := Get[byte](c, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if val != 'A' {
			t.Fatalf("expected %d, got %d", 'A', val)
		}
	})

	// 测试 []byte 类型
	t.Run("[]byte", func(t *testing.T) {
		key := "test:byteslice"
		want := []byte{1, 2, 3}
		err := Set(c, key, want, time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
		val, err := Get[[]byte](c, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if !bytes.Equal(val, want) {
			t.Fatalf("expected %v, got %v", want, val)
		}
	})

	// 测试 slice 类型
	t.Run("slice", func(t *testing.T) {
		key := "test:slice"
		want := []string{"a", "b", "c"}
		err := Set(c, key, want, time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
		val, err := Get[[]string](c, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if len(val) != len(want) || val[0] != want[0] || val[1] != want[1] || val[2] != want[2] {
			t.Fatalf("expected %v, got %v", want, val)
		}
	})

	// 测试指针类型
	t.Run("pointer", func(t *testing.T) {
		key := "test:pointer"
		want := &User{Name: "bob", Age: 30}
		err := Set(c, key, want, time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
		val, err := Get[*User](c, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if val.Name != want.Name || val.Age != want.Age {
			t.Fatalf("expected %+v, got %+v", want, val)
		}
	})

	// 测试 key 不存在
	t.Run("not found", func(t *testing.T) {
		_, err := Get[string](c, "nonexistent")
		if err != ErrNotFound {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestRedisCache_GetSet(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Skipf("redis not available: %v", err)
	}
	c := NewRedisCache(client)

	// 测试字符串类型
	t.Run("string", func(t *testing.T) {
		key := "test:str"
		err := Set(c, key, "hello", time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
		val, err := Get[string](c, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if val != "hello" {
			t.Fatalf("expected hello, got %s", val)
		}
	})

	// 测试结构体类型
	t.Run("struct", func(t *testing.T) {
		key := "test:struct"
		want := User{Name: "alice", Age: 20}
		err := Set(c, key, want, time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
		val, err := Get[User](c, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if val != want {
			t.Fatalf("expected %+v, got %+v", want, val)
		}
	})

	// 测试 int 类型
	t.Run("int", func(t *testing.T) {
		key := "test:int"
		err := Set(c, key, 42, time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
		val, err := Get[int](c, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if val != 42 {
			t.Fatalf("expected 42, got %d", val)
		}
	})

	// 测试 float64 类型
	t.Run("float64", func(t *testing.T) {
		key := "test:float"
		err := Set(c, key, 3.14, time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
		val, err := Get[float64](c, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if val != 3.14 {
			t.Fatalf("expected 3.14, got %f", val)
		}
	})

	// 测试 byte 类型
	t.Run("byte", func(t *testing.T) {
		key := "test:byte"
		err := Set(c, key, byte('A'), time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
		val, err := Get[byte](c, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if val != 'A' {
			t.Fatalf("expected %d, got %d", 'A', val)
		}
	})

	// 测试 []byte 类型
	t.Run("[]byte", func(t *testing.T) {
		key := "test:byteslice"
		want := []byte{1, 2, 3}
		err := Set(c, key, want, time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
		val, err := Get[[]byte](c, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if !bytes.Equal(val, want) {
			t.Fatalf("expected %v, got %v", want, val)
		}
	})

	// 测试 slice 类型
	t.Run("slice", func(t *testing.T) {
		key := "test:slice"
		want := []string{"a", "b", "c"}
		err := Set(c, key, want, time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
		val, err := Get[[]string](c, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if len(val) != len(want) || val[0] != want[0] || val[1] != want[1] || val[2] != want[2] {
			t.Fatalf("expected %v, got %v", want, val)
		}
	})

	// 测试指针类型
	t.Run("pointer", func(t *testing.T) {
		key := "test:pointer"
		want := &User{Name: "bob", Age: 30}
		err := Set(c, key, want, time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
		val, err := Get[*User](c, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if val.Name != want.Name || val.Age != want.Age {
			t.Fatalf("expected %+v, got %+v", want, val)
		}
	})

	// 测试 key 不存在
	t.Run("not found", func(t *testing.T) {
		_, err := Get[string](c, "nonexistent")
		if err != ErrNotFound {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	// 清理测试数据
	c.Delete("test:str", "test:struct", "test:int", "test:float", "test:byte", "test:byteslice", "test:slice", "test:pointer")
}
