package main

import (
	"github.com/djumanoff/amqp"
	lib "github.com/kirigaikabuto/recommendation-system-score-store"
	setdata_common "github.com/kirigaikabuto/setdata-common"
)

var (
	postgresUser         = "setdatauser"
	postgresPassword     = "123456789"
	postgresDatabaseName = "recommendation_system"
	postgresHost         = "localhost"
	postgresPort         = 5432
	postgresParams       = "sslmode=disable"
	amqpUrl              = "amqp://localhost:5672"
)

func main() {
	config := lib.Config{
		Host:     postgresHost,
		Port:     postgresPort,
		User:     postgresUser,
		Password: postgresPassword,
		Database: postgresDatabaseName,
		Params:   postgresParams,
	}
	store, err := lib.NewPostgreStore(config)
	if err != nil {
		panic(err)
		return
	}
	service := lib.NewScoreService(store)
	commandHandler := setdata_common.NewCommandHandler(service)
	scoreAmqpEndpoints := lib.NewScoreAmqpEndpoints(commandHandler)
	rabbitConfig := amqp.Config{
		AMQPUrl:  amqpUrl,
		LogLevel: 5,
	}
	serverConfig := amqp.ServerConfig{
		ResponseX: "response",
		RequestX:  "request",
	}

	sess := amqp.NewSession(rabbitConfig)
	err = sess.Connect()
	if err != nil {
		panic(err)
		return
	}
	srv, err := sess.Server(serverConfig)
	if err != nil {
		panic(err)
		return
	}
	srv.Endpoint("score.create", scoreAmqpEndpoints.CreateScoreAmqpEndpoint())
	srv.Endpoint("score.list", scoreAmqpEndpoints.ListScoreAmqpEndpoint())
	err = srv.Start()
	if err != nil {
		panic(err)
		return
	}
}
