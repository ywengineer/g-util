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
	dsn  string
	conn *sqlx.DB
	log  *zap.Logger
}

func NewMySQL(dbUser, dbPassword, dbHost, dbName, locale string,
	dbPort, dbWriteTimeout, dbReadTimeout, dbDialTimeout int,
	maxOpenConn, maxIdleConn int,
	log *zap.Logger) *MySQL {
	// "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&time_zone=%s&loc=%s&autocommit=true&writeTimeout=%ds&readTimeout=%ds&timeout=%ds"
	if len(locale) == 0 {
		locale = "Asia/Shanghai"
	}
	if _, e := time.LoadLocation(locale); e != nil {
		log.Panic("unknown locale", zap.String("loc", locale))
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&loc=%s&autocommit=true&writeTimeout=%ds&readTimeout=%ds&timeout=%ds&interpolateParams=true",
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
		return nil
	} else {
		//
		db.SetMaxOpenConns(maxOpenConn)
		db.SetMaxIdleConns(maxIdleConn)
		//
		if err := db.Ping(); err == nil {
			return &MySQL{conn: db, log: log, dsn: dsn}
		} else {
			log.Fatal("ping mysql failed", zap.Error(err))
			return nil
		}
	}
}

func (mysql *MySQL) String() string {
	return mysql.dsn
}

func (mysql *MySQL) Info() map[string]interface{} {
	s := mysql.conn.Stats()
	return map[string]interface{}{
		"MaxOpenConnections(Maximum number of open connections to the database)": s.MaxOpenConnections,
		//
		"OpenConnections(The number of established connections both in use and idle)": s.OpenConnections,
		"Idle(The number of idle connections)":                                        s.Idle,
		"InUse(The number of connections currently in use)":                           s.InUse,
		//
		"WaitCount(The total number of connections waited for)":             s.WaitCount,
		"WaitDuration(The total time blocked waiting for a new connection)": s.WaitDuration.Seconds(),
		//
		"MaxIdleClosed(The total number of connections closed due to SetMaxIdleConns)":        s.MaxIdleClosed,
		"MaxLifetimeClosed(The total number of connections closed due to SetConnMaxLifetime)": s.MaxLifetimeClosed,
	}
}

func (mysql *MySQL) GetConn() *sqlx.DB {
	mysql.conn.Stats()
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
