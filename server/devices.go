package main

import "time"

var recentDevices []deviceLastSeen

type device struct {
	DeviceMac string    `json:"device"`
	Location  string    `json:"location"`
	Date      time.Time `json:"date"`

	Powered         bool            `json:"powered"`
	LastStatus      lastMessageSent `json:"lastMessageSent"`
	ProgrammedTimes []Schedule      `json:"schedule"`
}

type lastMessageSent struct {
	Powered string `json:"powered"`
}

type deviceLastSeen struct {
	DeviceMac string `json:"device"`
	SeenEpoch int64  `json:"last_seen"`
}

func RecentDevices() []deviceLastSeen {
	return recentDevices
}

func updateDeviceLastSeen(deviceMacAddress string) {

	// remove entry in array
	for i, e := range recentDevices {
		if e.DeviceMac == deviceMacAddress {
			recentDevices = append(recentDevices[:i], recentDevices[i+1:]...)
		}
	}

	// update entry with new time
	recentDevices = append(recentDevices, deviceLastSeen{deviceMacAddress, time.Now().Unix()})
}
