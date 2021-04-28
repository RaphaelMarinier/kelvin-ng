package main

import (
	"testing"
	"time"
)

func TestReadOK(t *testing.T) {
	correctfiles := []string{
		"testdata/config-example.json",
		"testdata/config-example-newstyleschedule.json",
		"testdata/config-example.yaml",
	}
	for _, testFile := range correctfiles {
		c := Configuration{}
		c.ConfigurationFile = testFile
		err := c.Read()
		if err != nil {
			t.Fatalf("Could not read correct configuration file : %v with error : %v", c.ConfigurationFile, err)
		}
	}
}

type MockSunStateCalculator struct {
	MockSunrise time.Time
	MockSunset  time.Time
}

func (calculator *MockSunStateCalculator) CalculateSunset(date time.Time, latitude float64, longitude float64) time.Time {
	return calculator.MockSunset
}

func (calculator *MockSunStateCalculator) CalculateSunrise(date time.Time, latitude float64, longitude float64) time.Time {
	return calculator.MockSunrise
}

func TestLightScheduleForDay(t *testing.T) {
	c := Configuration{}
	c.ConfigurationFile = "testdata/config-example-newstyleschedule.json"
	err := c.Read()
	if err != nil {
		t.Fatalf("Could not read correct configuration file : %v with error : %v", c.ConfigurationFile, err)
	}
	location := time.UTC
	calculator := &MockSunStateCalculator{
		time.Date(2021, 4, 28, 7, 30, 0, 0, location),
		time.Date(2021, 4, 28, 20, 0, 0, 0, location)}

	s, err := c.lightScheduleForDay(1, time.Date(2021, 4, 28, 0, 0, 1, 0, location), calculator)
	if err != nil {
		t.Fatalf("Got error %v", err)
	}

	// TODO: check schedule, under different mock sunset/sunrise.
	//	[{2021-04-27 22:00:00 +0200 CEST 2000 70} {2021-04-28 04:00:00 +0200 CEST 2000 60} {2021-04-28 06:40:11 +0200 CEST 2700 60} {2021-04-28 07:10:11 +0200 CEST 5000 100} {2021-04-28 19:25:44 +0200 CEST 5000 100} {2021-04-28 19:55:44 +0200 CEST 2700 80} {2021-04-28 22:00:00 +0200 CEST 2000 70} {2021-04-29 04:00:00 +0200 CEST 2000 60}]

	parseTime := func(t string) time.Time {
		parsed, _ := time.Parse("2006-01-02 15:04:05", t)
		return parsed
	}

	expectedTimes := []TimeStamp{
		TimeStamp{parseTime("2021-04-27 22:00:00"), 2000, 70},
		TimeStamp{parseTime("2021-04-28 04:00:00"), 2000, 60},
		TimeStamp{parseTime("2021-04-28 07:30:00"), 2700, 60},
		TimeStamp{parseTime("2021-04-28 08:00:00"), 5000, 100},
		TimeStamp{parseTime("2021-04-28 19:30:00"), 5000, 100},
		TimeStamp{parseTime("2021-04-28 20:00:00"), 2700, 80},
		TimeStamp{parseTime("2021-04-28 22:00:00"), 2000, 70},
		TimeStamp{parseTime("2021-04-29 04:00:00"), 2000, 60}}

	if len(s.times) != len(expectedTimes) {
		t.Fatalf("Got schedule with unexpected length. Got %v expected %v", s.times, expectedTimes)
	}
	for i, expectedTime := range expectedTimes {
		if expectedTime != s.times[i] {
			t.Fatalf("Got unexpected timestamp at position %v. Got %v expected %v",
				i, s.times[i], expectedTime)
		}
	}
}

func TestReadError(t *testing.T) {
	wrongfiles := []string{
		"",          // no file passed
		"testdata/", // not a regular file
		"testdata/config-bad-wrongFormat.json",
		"testdata/config-bad-wrongFormat.yaml",
	}
	for _, testFile := range wrongfiles {
		c := Configuration{}
		c.ConfigurationFile = testFile
		err := c.Read()
		if err == nil {
			t.Errorf("reading [%v] file should return an error", c.ConfigurationFile)
		}
	}
}

func TestWriteOK(t *testing.T) {
	correctfiles := []string{
		"testdata/config-example.json",
		"testdata/config-example.yaml",
	}
	for _, testFile := range correctfiles {
		c := Configuration{}
		c.ConfigurationFile = testFile
		_ = c.Read()
		c.Hash = ""
		err := c.Write()
		if err != nil {
			t.Errorf("Could not write configuration to correct file : %v", c.ConfigurationFile)
		}
	}
}
