package config

import "time"

const Md5HashKey = "sdskj2HHDjs1"
const GRPCPortBackend = "8083"
const GRPCPort = "8082"
const HTTPPort = "8081"
const ShortDuration = 500 * time.Minute //time.Millisecond

const (
	// Database config
	Host     = "localhost"
	Port     = 6432
	User     = "user"
	Password = "password"
	DBname   = "gohw"

	MaxConnIdleTime = time.Minute
	MaxConnLifetime = time.Hour
	MinConns        = 2
	MaxConns        = 4
)

const (
	DefaultRecPerPage   = 5
	DefaultPageNum      = 1
	DefaultSortingField = "id"
)

var Brokers = []string{"localhost:19091", "localhost:29091", "localhost:39091"}

const (
	ConsumerGroupUI     = "uiResponseConsuming"
	ConsumerGroupClient = "clientRequestConsuming"
)

const (
	TopicClientRequest = "client_requests"
	TopicUIResponse    = "ui_response"
	TopicUIRequest     = "ui_request"
)

type RedisCfg struct {
	Addr     string
	Password string
	DbNum    int
}

var RedisConfig = RedisCfg{"127.0.0.1:6379", "", 1}
