package di

import (
	"event-service/graph"
	"event-service/internal/config"
	"event-service/internal/database"
	gorminternal "event-service/internal/database/gorm"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type Container struct {
	db             *gorm.DB
	ampqConnection map[string]*amqp.Connection
	ampqChannels   map[string]*amqp.Channel
	ampqQueues     map[string]struct{}
}

var container = &Container{
	db:             nil,
	ampqConnection: make(map[string]*amqp.Connection),
	ampqChannels:   make(map[string]*amqp.Channel),
	ampqQueues:     make(map[string]struct{}),
}

func DatabaseParameters() database.Parameters {
	return database.Parameters{
		Host:     config.GetString("MYSQL.HOST"),
		Database: config.GetString("MYSQL.DATABASE"),
		Username: config.GetString("MYSQL.USERNAME"),
		Password: config.GetString("MYSQL.PASSWORD"),
		Port:     config.GetIntOrFallback("MYSQL.PORT", 3300),
	}
}

func MysqlDriver() string {
	return database.DriverName
}

func GORM() *gorm.DB {
	if container.db != nil {
		return container.db
	}

	db, err := gorminternal.Connection(DatabaseParameters())
	if err != nil {
		panic(err)
	}

	container.db = db

	return container.db
}

func DefaultGraphQLApiResolver() (r *graph.Resolver, err error) {
	r = &graph.Resolver{}

	if r.AddEventHandler, err = DefaultEventsAddHandler(); err != nil {
		return nil, err
	}

	if r.FindEventsHandler, err = DefaultEventsListHandler(); err != nil {
		return nil, err
	}

	if r.InvitationHandler, err = DefaultInvitationHandler(); err != nil {
		return nil, err
	}

	if r.UpdateEventHandler, err = DefaultEventsUpdateHandler(); err != nil {
		return nil, err
	}

	return r, nil
}
