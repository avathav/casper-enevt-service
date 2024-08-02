package gorm

import (
	"context"
	"errors"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"event-service/internal/database"

	mysqldriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ContextKey string

const (
	connectionKey ContextKey = "connection"
)

func Connection(params database.Parameters) (*gorm.DB, error) {
	conn, err := database.MysqlConnection(params)
	if err != nil {
		return nil, err
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	return gorm.Open(mysqldriver.New(mysqldriver.Config{Conn: conn}), &gorm.Config{
		Logger: newLogger,
	})
}

func ContextWithConnection(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, connectionKey, db)
}

func ConnectionFromContext(ctx context.Context) (*gorm.DB, error) {
	v, ok := ctx.Value(connectionKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("no gorm connection in context")
	}

	return v.WithContext(ctx), nil
}
