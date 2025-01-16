package conf

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

var (
	configPath     = "."
	configFileName = "conf.json"
)

type GbeConfig struct {
	DataSource  DataSourceConfig `json:"dataSource"`
	Rest        RestServer       `json:"rest"`
	JwtSecret   string           `json:"jwtSecret"`
	Password    string           `json:"password"`
	ProjectInfo ProjectInfo      `json:"projectInfo"`
	Env         string           `json:"env"`
	RabbitMQ    RabbitMQ         `json:"rabbitmq"`
}

type DataSourceConfig struct {
	Host              string `json:"host"`
	Port              string `json:"port"`
	Database          string `json:"database"`
	User              string `json:"user"`
	Password          string `json:"password"`
	SslMode           string `json:"sslMode"`
	EnableAutoMigrate bool   `json:"enableAutoMigrate"`
	Retries           int    `json:"retries"`
	Mode              int    `json:"mode"`
}

type ProjectInfo struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type RabbitMQ struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type RestServer struct {
	Addr string `json:"addr"`
}

type Server struct {
	ServerMinTime time.Duration
	ServerTime    time.Duration
	ServerTimeout time.Duration
}

var config GbeConfig
var configOnce sync.Once

func SetConfFilePath(path string) {
	configPath = path
}

func SetConfFileName(name string) {
	configFileName = name
}

func GetConfig() *GbeConfig {
	configOnce.Do(func() {

		bytes, err := os.ReadFile(configPath + "/" + configFileName)
		log.Println(configPath + "/" + configFileName)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(bytes, &config)
		if err != nil {
			panic(err)
		}
	})

	return &config
}

func GetTestConf() *GbeConfig {
	return &testConfig
}

var testConfig = GbeConfig{
	DataSource: DataSourceConfig{
		Host:              "localhost",
		Port:              "8432",
		Database:          "postgres",
		User:              "postgres",
		Password:          "password",
		SslMode:           "disable",
		EnableAutoMigrate: true,
		Retries:           3,
		Mode:              1,
	},
	Password: "Test",
	Rest: RestServer{
		Addr: "localhost:8080",
	},
	JwtSecret: "supersecret",
	ProjectInfo: ProjectInfo{
		Name:   "TestProject",
		Domain: "example.com",
	},
	RabbitMQ: RabbitMQ{
		Host: "127.0.0.1",
		Port: "4672",
	},
	Env: "test",
}
