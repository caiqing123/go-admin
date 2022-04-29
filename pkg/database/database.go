// Package database 数据库操作
package database

import (
	"database/sql"
	"errors"
	"fmt"

	"api/pkg/config"
	"api/pkg/search"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DB 对象
var DB *gorm.DB
var SQLDB *sql.DB

// Connect 连接数据库
func Connect(dbConfig gorm.Dialector, _logger gormlogger.Interface) {

	// 使用 gorm.Open 连接数据库
	var err error
	DB, err = gorm.Open(dbConfig, &gorm.Config{
		Logger: _logger,
	})
	// 处理错误
	if err != nil {
		fmt.Println(err.Error())
	}

	// 获取底层的 sqlDB
	SQLDB, err = DB.DB()
	if err != nil {
		fmt.Println(err.Error())
	}
}
func CurrentDatabase() (dbname string) {
	dbname = DB.Migrator().CurrentDatabase()
	return
}

func DeleteAllTables() error {

	var err error

	switch config.Get("database.connection") {
	case "mysql":
		err = deleteMysqlDatabase()
	case "sqlite":
		_ = deleteAllSqliteTables()
	default:
		panic(errors.New("database connection not supported"))
	}

	return err
}

func deleteAllSqliteTables() error {
	var tables []string
	DB.Select(&tables, "SELECT name FROM sqlite_master WHERE type='table'")
	for _, table := range tables {
		_ = DB.Migrator().DropTable(table)
	}
	return nil
}

func deleteMysqlDatabase() error {
	dbname := CurrentDatabase()
	sqls := fmt.Sprintf("DROP DATABASE %s;", dbname)
	if err := DB.Exec(sqls).Error; err != nil {
		return err
	}
	sqls = fmt.Sprintf("CREATE DATABASE %s;", dbname)
	if err := DB.Exec(sqls).Error; err != nil {
		return err
	}
	sqls = fmt.Sprintf("USE %s;", dbname)
	if err := DB.Exec(sqls).Error; err != nil {
		return err
	}
	return nil
}

func TableName(obj interface{}) string {
	stmt := &gorm.Statement{DB: DB}
	_ = stmt.Parse(obj)
	return stmt.Schema.Table
}

// MakeCondition 生成查询语言
func MakeCondition(q interface{}, db *gorm.DB) *gorm.DB {
	condition := &search.GormCondition{
		GormPublic: search.GormPublic{},
		Join:       make([]*search.GormJoin, 0),
	}
	search.ResolveSearchQuery(config.Get("database.connection"), q, condition)

	for _, join := range condition.Join {
		if join == nil {
			continue
		}
		db = db.Joins(join.JoinOn)
		for k, v := range join.Where {
			db = db.Where(k, v...)
		}
		for k, v := range join.Or {
			db = db.Or(k, v...)
		}
	}
	for k, v := range condition.Where {
		db = db.Where(k, v...)
	}
	for k, v := range condition.Or {
		db = db.Or(k, v...)
	}
	return db
}
