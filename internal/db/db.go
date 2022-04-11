package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDB(opts ...Option) *gorm.DB {
	dsn := defaultDsn
	dsn.ApplyOptions(opts...)
	logger.Println("Opening Postgres DB with options:", dsn.CensoredString())

	// We want default transactions to be enabled,
	// data consistency is much more important than performance (for now)
	// skipcq: GO-W1004
	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{
		SkipDefaultTransaction: false,
	})

	if err != nil {
		logger.Println("error opening Postgres DB:", err)
		return nil
	}
	return db
}

func InitializeDB(db *gorm.DB) error {
	// essentially doing CREATE TABLE IF NOT EXISTS for each model, nothing more exotic than that

	// remember to add new Models here when they're added
	// do one by one to avoid a set of problems with models depending on each other
	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&Record{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&Machine{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&LoginInfo{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&MachineToken{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&UserToken{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&Organization{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&Schedule{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&Task{}); err != nil {
		return err
	}

	return nil
}

var defaultDsn = dsn{
	"host":     "localhost",
	"user":     "postgres",
	"dbname":   "postgres",
	"port":     5432,
	"sslmode":  "disable",
	"TimeZone": "Etc/UTC",
}

type dsn map[string]interface{}

func (dsn dsn) ApplyOptions(opts ...Option) {
	for _, o := range opts {
		o(dsn)
	}
}

func (dsn dsn) String() string {
	if connString, present := dsn["raw"]; present {
		return connString.(string)
	}

	var r string
	for k, v := range dsn {
		r += fmt.Sprintf("%s=%v ", k, v)
	}

	return r
}

func (dsn dsn) CensoredString() string {
	var r string
	for k, v := range dsn {
		if k == "password" {
			v = "******"
		}
		r += fmt.Sprintf("%s='%v' ", k, v)
	}

	return r
}

type Option func(dsn)

func WithConnString(connString string) Option {
	return func(d dsn) {
		d["raw"] = connString
	}
}

func WithHost(host string) Option {
	return func(d dsn) {
		d["host"] = host
	}
}

func WithUsername(u string) Option {
	return func(d dsn) {
		d["user"] = u
	}
}

func WithPassword(pw string) Option {
	return func(d dsn) {
		d["password"] = pw
	}
}

func WithDBName(db string) Option {
	return func(d dsn) {
		d["dbname"] = db
	}
}

func WithPort(p int) Option {
	return func(d dsn) {
		d["port"] = p
	}
}

func WithSSLMode(m string) Option {
	return func(d dsn) {
		d["sslmode"] = m
	}
}

func WithTimeZone(tz string) Option {
	return func(d dsn) {
		d["TimeZone"] = tz
	}
}
