package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func setPower(client MQTT.Client, m macAddress, on bool) {
	topic := fmt.Sprintf("%s/%s", "device", m)
	turnOnJSON, _ := json.Marshal(map[string]string{"power": strconv.FormatBool(on)})

	if token := client.Publish(topic, 2, true, string(turnOnJSON)); token.Wait() && token.Error() != nil {

	}
}

func sendPing(client MQTT.Client) {
	if token := client.Publish("ping", 2, false, ""); token.Wait() && token.Error() != nil {

	}
}
