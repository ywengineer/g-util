package sql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/url"
	"time"
)

type MySQL struct {
	conn *sqlx.DB
	log  *zap.Logger
}

func NewMySQL(dbUser, dbPassword, dbHost, dbName, locale string,
	dbPort, dbWriteTimeout, dbReadTimeout, dbDialTimeout int,
	maxOpenConn, maxIdleConn int,
	log *zap.Logger) (*MySQL, error) {
	// "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&time_zone=%s&loc=%s&autocommit=true&writeTimeout=%ds&readTimeout=%ds&timeout=%ds"
	if len(locale) == 0 {
		locale = "Asia/Shanghai"
	}
	if _, e := time.LoadLocation(locale); e != nil {
		log.Panic("unknown locale", zap.String("loc", locale))
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&loc=%s&autocommit=true&writeTimeout=%ds&readTimeout=%ds&timeout=%ds",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
		url.QueryEscape(locale),
		dbWriteTimeout,
		dbReadTimeout,
		dbDialTimeout)
	if db, err := sqlx.Connect("mysql", dsn); err != nil {
		log.Fatal("connect mysql failed", zap.Error(err))
		return nil, err
	} else {
		//
		db.SetMaxOpenConns(maxOpenConn)
		db.SetMaxIdleConns(maxIdleConn)
		//
		if err := db.Ping(); err == nil {
			return &MySQL{conn: db, log: log}, nil
		} else {
			log.Fatal("ping mysql failed", zap.Error(err))
			return nil, err
		}
	}
}

func (mysql *MySQL) GetConn() *sqlx.DB {
	return mysql.conn
}

func (mysql *MySQL) DoExecSql(sql string) error {
	if r, err := mysql.conn.Exec(sql); err != nil {
		mysql.log.Error("execute sql failed", zap.String("tag", "SQL"), zap.String("sql", sql), zap.Error(err))
		return err
	} else if rows, re := r.RowsAffected(); re == nil {
		mysql.log.Info("execute sql succeed", zap.String("tag", "SQL"), zap.String("sql", sql), zap.Int64("rows", rows))
		return nil
	} else {
		mysql.log.Error("execute sql failed", zap.String("tag", "SQL"), zap.String("sql", sql), zap.Error(re))
		return re
	}
}
