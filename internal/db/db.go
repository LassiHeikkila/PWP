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
	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{})
	if err != nil {
		logger.Println("error opening Postgres DB:", err)
		return nil
	}
	return db
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
	var r string
	for k, v := range dsn {
		r += fmt.Sprintf("%s='%v' ", k, v)
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
