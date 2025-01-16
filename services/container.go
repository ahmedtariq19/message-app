package services

import (
	"message-app/conf"
	"message-app/pkg/rabbitmq"
	"message-app/repository"
	"message-app/repository/postgres"
	"os"
)

type Container struct {
	Store                 repository.Store
	GbeConfigService      GbeConfigService
	Producer              rabbitmq.Producer
	Consumer              rabbitmq.Consumer
	MessageService        MessageService
	JWTService            JWTService
	AuthenticationService AuthenticationService
}

func NewServiceContainer() *Container {
	gbeConfig := conf.GetConfig()
	gbeConfigService := NewGbeConfigService()
	postgreStore := postgres.SharedStore(gbeConfig)
	producer, err := rabbitmq.NewProducer(gbeConfig.RabbitMQ)
	if err != nil {
		os.Exit(1)
	}

	consumer, err := rabbitmq.NewConsumer(gbeConfig.RabbitMQ)
	if err != nil {
		os.Exit(1)
	}

	messageService := NewMessageService(postgreStore)
	jwtService := NewJWTService(gbeConfigService)
	authenticationService := NewAuthenticationService(jwtService, gbeConfigService)
	return &Container{
		Store:                 postgreStore,
		GbeConfigService:      gbeConfigService,
		Producer:              producer,
		Consumer:              consumer,
		MessageService:        messageService,
		JWTService:            jwtService,
		AuthenticationService: authenticationService,
	}
}
