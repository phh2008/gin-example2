package orm

import (
	"fmt"
	"log/slog"
	"time"

	"com.example/example/model"
	"com.example/example/pkg/config"
	"com.example/example/pkg/logger/glog"
	"com.example/example/pkg/xstring"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// NewDB 创建 gorm.DB
func NewDB(conf *config.Config) *gorm.DB {
	cfg := glog.NewConfig(slog.Default().Handler()).
		//WithGroupKey("db").
		WithSlowThreshold(time.Second).
		WithIgnoreRecordNotFoundError(true).
		WithTraceAll(true)
	glogger := glog.NewWithConfig(cfg)
	var gdb, err = gorm.Open(mysql.Open(conf.Db.Dsn1), &gorm.Config{
		Logger: glogger,
	})
	if err != nil {
		panic(err)
	}
	err = gdb.Use(dbresolver.Register(dbresolver.Config{
		Sources: []gorm.Dialector{mysql.Open(conf.Db.Dsn2)},
	}, "sys_permission").Register(dbresolver.Config{
		Sources: []gorm.Dialector{mysql.Open(conf.Db.Dsn3)},
	}, "sys_role").
		SetConnMaxIdleTime(time.Hour).
		SetConnMaxLifetime(24 * time.Hour).
		SetMaxIdleConns(10).
		SetMaxOpenConns(50))
	if err != nil {
		panic(err)
	}
	// 设置连接池参数
	//sqlDB, err := gdb.DB()
	//if err != nil {
	//	panic(err)
	//}
	//sqlDB.SetMaxIdleConns(10)           // 空闲最大连接数
	//sqlDB.SetMaxOpenConns(60)           // 最大打开连接数
	//sqlDB.SetConnMaxLifetime(time.Hour) // 连接可重用的时长
	return gdb
}

func emptyScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

func OrderScope(field string, direct string) func(db *gorm.DB) *gorm.DB {
	if field == "" {
		return emptyScope()
	}
	if !xstring.ColumnReg.MatchString(field) {
		// 非法字段
		return emptyScope()
	}
	if direct == "" || !xstring.DirectReg.MatchString(direct) {
		direct = "asc"
	}
	sort := fmt.Sprintf("%s %s", field, direct)
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(sort)
	}
}

func PageScope(pageNo, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNo <= 0 {
			pageNo = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		offset := (pageNo - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// QueryPage 分页查询
func QueryPage[T any](db *gorm.DB, pageNo, pageSize int) (model.PageData[T], error) {
	return QueryOrderPage[T](db, "", "", pageNo, pageSize)
}

func OrderPageScope[T any](page model.QueryPage, data *model.PageData[T]) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageNo := page.GetPageNo()
		pageSize := page.GetPageSize()
		data.PageNo = pageNo
		data.PageSize = pageSize
		offset := (pageNo - 1) * pageSize
		return db.Scopes(OrderScope(page.Sort, page.Direction)).Offset(offset).Limit(pageSize)
	}
}

// QueryOrderPage 排序分页查询
func QueryOrderPage[T any](db *gorm.DB, sortField string, direct string, pageNo, pageSize int) (model.PageData[T], error) {
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	var pageData model.PageData[T]
	pageData.PageNo = pageNo
	pageData.PageSize = pageSize
	offset := (pageNo - 1) * pageSize
	err := db.Count(&pageData.Count).Scopes(OrderScope(sortField, direct)).Offset(offset).Limit(pageSize).Find(&pageData.Data).Error
	return pageData, err
}
