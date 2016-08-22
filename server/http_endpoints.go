package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

func calcuateLastSeenDevices() []deviceLastSeen {
	var calculatedRecentlySeen []deviceLastSeen

	for _, e := range RecentDevices() {
		calculatedRecentlySeen = append(calculatedRecentlySeen, deviceLastSeen{e.DeviceMac, time.Now().Unix() - e.SeenEpoch})
	}
	return calculatedRecentlySeen
}

func allSeenDevices() []byte {
	allDevices, _ := json.Marshal(db.getDevices())
	return allDevices
}

func devicesEndpoint() interface{} {
	type response struct {
		AllDevices       []device         `json:"all_devices"`
		ConnectedDevices []deviceLastSeen `json:"connected_devices"`
	}
	return response{db.getDevices(), calcuateLastSeenDevices()}
}

func setDeviceLocationEndpoint(mac, location string) error {
	return db.updateDeviceLocation(mac, location)
}

func deleteDeviceEndpoint(mac string) error {
	mac = strings.ToLower(mac)
	if _, err := net.ParseMAC(mac); err != nil {
		return errors.New("Invalid mac address")
	}
	return db.deleteDevice(mac)
}

func setDeviceScheduleEndpoint(body io.ReadCloser) error {
	b, err := ioutil.ReadAll(body)

	if err != nil {
		return err
	}

	var s []Schedule

	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	for _, e := range s {
		if err := e.validateSchedule(); err != nil {
			return err
		}
		fmt.Println("Valid schedule: ", e)
	}

	return nil
}
