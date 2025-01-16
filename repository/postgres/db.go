package postgres

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"sync"
	"time"

	"message-app/conf"
	"message-app/models"
	"message-app/repository"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gdb *gorm.DB
var storeOnce sync.Once
var store repository.Store

type Store struct {
	db *gorm.DB
}

// SharedStore return global or single instance of postgres connection (bounded in sync once)
func SharedStore(cfg *conf.GbeConfig) repository.Store {
	storeOnce.Do(func() {
		err := initDb(cfg)
		if err != nil {
			panic(err)
		}
		store = NewStore(gdb)
	})
	return store
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

func initDb(cfg *conf.GbeConfig) error {

	retries := cfg.DataSource.Retries
	var err error

	dsn := fmt.Sprintf(
		"user=%v password=%v host=%v port=%v dbname=%v sslmode=%v",
		cfg.DataSource.User, cfg.DataSource.Password, cfg.DataSource.Host, cfg.DataSource.Port, cfg.DataSource.Database, cfg.DataSource.SslMode,
	)

	config := &gorm.Config{
		Logger: logger.Default.LogMode(getLogLevel(cfg.DataSource.Mode)),
	}

	gdb, err = gorm.Open(postgres.Open(dsn), config)
	for err != nil {
		log.Println(err, fmt.Sprintf("Failed to connect to database (%d)", retries))

		if retries > 1 {
			retries--
			time.Sleep(5 * time.Second)
			gdb, err = gorm.Open(postgres.Open(dsn), config)
			continue
		}
		os.Exit(1)
	}

	db, err := gdb.DB()
	if err != nil {
		os.Exit(1)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(50)

	if cfg.DataSource.EnableAutoMigrate {
		var tables = []interface{}{
			&models.Message{},
		}
		for _, table := range tables {
			log.Printf("migrating database, table: %v", reflect.TypeOf(table))
			if err = gdb.AutoMigrate(table); err != nil {
				return err
			}
		}
	}

	return nil
}

func getLogLevel(mode int) logger.LogLevel {
	switch mode {
	case 0:
		return logger.Silent
	case 1:
		return logger.Error
	case 2:
		return logger.Warn
	case 3:
		return logger.Info
	default:
		return logger.Silent
	}
}

func (s *Store) GetDb() *gorm.DB {
	return s.db
}
