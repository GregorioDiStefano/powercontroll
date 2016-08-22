package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

var db Database

type macAddress string

func setupConfig() viper.Viper {
	config := viper.New()
	config.SetConfigName("config")
	config.AddConfigPath(".")

	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	return *config
}

func setupMQTT(config viper.Viper) (MQTT.Client, error) {
	var c MQTT.Client

	opts := MQTT.NewClientOptions().AddBroker(config.GetString("broker_uri"))
	opts.SetUsername(config.GetString("username"))
	opts.SetPassword(config.GetString("password"))
	opts.SetClientID("server")
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(false)
	opts.SetOnConnectHandler(connectionHandler)
	opts.SetConnectionLostHandler(connectionLostHandler)
	c = MQTT.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, errors.New("Unable to connect to: " + config.GetString("broker_uri"))
	}
	fmt.Println("Connected.")

	c.Subscribe("announce", 0x00, announceHandler)
	c.Subscribe("pong", 0x00, pongHandler)
	c.Subscribe("lwt", 0x00, announceHandler)
	c.Subscribe("status_alert", 0x00, statusAlertHandler)

	return c, nil
}

func devicePinger(client MQTT.Client, config viper.Viper) {
	c := time.Tick(time.Second * time.Duration(config.GetInt("ping_frequency")))

	for {
		fmt.Println("Sending ping..")
		sendPing(client)
		<-c
	}
}

func main() {
	trap := make(chan os.Signal)
	signal.Notify(trap, os.Interrupt)

	config := setupConfig()
	client, err := setupMQTT(config)

	if err != nil {
		panic(err)
	}

	db = new(ormDatabase)
	db.openDatabase()

	// clean up routine on exit
	go func() {
		<-trap
		db.closeDatabase()
		fmt.Println("Exiting..")
		os.Exit(0)
	}()

	go devicePinger(client, config)
	ginEngine := setupRouting(client)
	ginEngine.Run()

	for {
		time.Sleep(1 * time.Second)
	}
}
