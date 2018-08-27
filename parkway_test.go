package main

import (
	"fmt"
	"testing"
	"time"
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

func TestLoadMeals(t *testing.T) {
	meals, err := buildMeals("https://www.parkwayschools.net/site/handlers/icalfeed.ashx?MIID=4134", "lunch")
	if err != nil {
		fmt.Println("errors! ", err)
	}
	for _, m := range meals {
		fmt.Println("meal: ", m)
	}
	t.Errorf("jdkfjfkdj %s", "")
}

func TestStuff(t *testing.T) {
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	t.Errorf("%s", today.Format(tFormat))
}
