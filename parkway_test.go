package main

import (
	"testing"
)

func TestGetToday(t *testing.T) {
	out := GetTodaysLunch()
	if out != "" {
		t.Errorf("got %s", out)
	}
}

func TestGetWeek(t *testing.T) {
	buildLunchMap()
	out := GetWeek()
	if out != "" {
		t.Errorf("got %s", out)
	}
}

func TestGetCalendar(t *testing.T) {
	//getCalendar()
	out := getCalendar()
	if out != "" {
		t.Errorf("got %s", out)
	}
}
