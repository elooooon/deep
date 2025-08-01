package config

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/logx"
)

// DatabaseConf stores database configurations.
type DatabaseConf struct {
	Host        string `json:",env=DATABASE_HOST"`
	Port        int    `json:",env=DATABASE_PORT"`
	Username    string `json:",default=root,env=DATABASE_USERNAME"`
	Password    string `json:",optional,env=DATABASE_PASSWORD"`
	DBName      string `json:",default=simple_admin,env=DATABASE_DBNAME"`
	SSLMode     string `json:",optional,env=DATABASE_SSL_MODE"`
	Type        string `json:",default=postgres,options=[mysql,postgres],env=DATABASE_TYPE"`
	MaxOpenConn int    `json:",optional,default=100,env=DATABASE_MAX_OPEN_CONN"`
	MysqlConfig string `json:",optional,env=DATABASE_MYSQL_CONFIG"`
	PGConfig    string `json:",optional,env=DATABASE_PG_CONFIG"`
}

// NewNoCacheDriver returns an Ent driver without cache.
func (c DatabaseConf) NewNoCacheDriver() *entsql.Driver {
	db, err := sql.Open(c.Type, c.GetDSN())
	logx.Must(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	logx.Must(err)

	db.SetMaxOpenConns(c.MaxOpenConn)
	return entsql.OpenDB(c.Type, db)
}

// MysqlDSN returns mysql DSN.
func (c DatabaseConf) MysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True%s", c.Username, c.Password, c.Host, c.Port, c.DBName, c.MysqlConfig)
}

// PostgresDSN returns Postgres DSN.
func (c DatabaseConf) PostgresDSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s%s",
		c.Username, c.Password, c.Host, c.Port, c.DBName, c.SSLMode, c.PGConfig)
}

// GetDSN returns DSN according to the database type.
func (c DatabaseConf) GetDSN() string {
	switch c.Type {
	case "mysql":
		return c.MysqlDSN()
	case "postgres":
		return c.PostgresDSN()
	default:
		return c.PostgresDSN()
	}
}
