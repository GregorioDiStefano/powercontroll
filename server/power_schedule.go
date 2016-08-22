package main

import (
	"errors"
	"fmt"
	"time"
)

type Schedule struct {
	Dates         []string `json:"dates"`
	Days          []string `json:"days"`
	SwitchOnTime  string   `json:"switch_on_time"`
	SwitchOffTime string   `json:"switch_off_time"`
}

type PowerSchedule interface {
	validateSchedule(Schedule) (Schedule, error)
	shouldDeviceBeOn(Schedule) bool
}

func (s Schedule) validateSchedule() error {
	var atLeastOneSet bool
	var validDates []time.Time
	var validDays []time.Time

	fmt.Println(s)

	if len(s.Dates) > 0 {
		for _, date := range s.Dates {
			t, err := time.Parse("02/01/2006", date)

			if err != nil {
				return errors.New("Date invalid: does not follow DD/MM/YYYY format")
			}
			validDates = append(validDates, t)
		}
		atLeastOneSet = true
	}

	if len(s.Days) > 0 {
		for _, day := range s.Days {
			t, err := time.Parse("Monday", day)

			if err != nil {
				return errors.New("Day invalid: does not follow 'Monday' format")
			}
			validDays = append(validDays, t)
		}
		atLeastOneSet = true
	}

	if len(s.SwitchOnTime) > 0 {
		_, err := time.Parse("15:04:05", s.SwitchOnTime)
		if err != nil {
			return errors.New("Switch on time invalid: does not follow HH:MM:SS")
		}
		atLeastOneSet = true
	}

	if len(s.SwitchOffTime) > 0 {
		_, err := time.Parse("15:04:05", s.SwitchOffTime)
		if err != nil {
			return errors.New("Switch off time invalid: does not follow HH:MM:SS")
		}
		atLeastOneSet = true
	}

	if atLeastOneSet == false {
		return errors.New("No valid schedule data")
	}

	return nil
}
