package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/kanziw/howareyou/config"
)

const driverName = "mysql"

type ConnectionOption struct {
	Host     string
	Port     int
	DBName   string
	User     string
	Password string
}

func GetDB(setting config.Setting) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		setting.DBUser,
		setting.DBPassword,
		setting.DBHost,
		setting.DBPort,
		setting.DBName,
	)

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		logrus.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "ping timeout")
	}

	db.SetMaxIdleConns(setting.DBMaxIdleConns)
	db.SetMaxOpenConns(setting.DBMaxOpenConns)
	db.SetConnMaxLifetime(setting.DBConnMaxLifetime)

	return db, nil
}
