package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
)

func setupRouting(client MQTT.Client) *gin.Engine {
	r := gin.Default()

	r.GET("/devices", func(c *gin.Context) {
		c.IndentedJSON(200, devicesEndpoint())
	})

	r.GET("/devices/:mac/status", func(c *gin.Context) {
		c.IndentedJSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/devices/:mac/off", func(c *gin.Context) {
		setPower(client, macAddress(c.Params.ByName("mac")), false)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/devices/:mac/on", func(c *gin.Context) {
		setPower(client, macAddress(c.Params.ByName("mac")), true)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.DELETE("/devices/:mac", func(c *gin.Context) {
		if err := deleteDeviceEndpoint(c.Params.ByName("mac")); err != nil {
			c.JSON(500, err.Error())
			return
		}
		c.Status(204)
	})

	r.PUT("/devices/:mac/schedule", func(c *gin.Context) {
		if err := setDeviceScheduleEndpoint(c.Request.Body); err != nil {
			c.JSON(500, err.Error())
			return
		}
		c.Status(204)
	})

	r.PUT("/devices/:mac/location", func(c *gin.Context) {
		if err := setDeviceLocationEndpoint(c.Params.ByName("mac"), "test"); err != nil {
			c.JSON(500, err.Error())
			return
		}
		c.Status(204)
	})

	return r
}
