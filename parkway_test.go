package main

import (
	"errors"
	"strings"
	"testing"
	"time"
)

type fakecal struct{}

type fakedb struct{}

func (fd fakedb) Get(date time.Time, mealType string, details bool) (Meal, error) {
	if mealType == "error" {
		return Meal{}, errors.New("you asked for one")
	}

	m := Meal{
		Date:     date,
		ID:       "myid",
		MainDish: "BAAAAAAAACON!!!",
		Sides:    "RANCH AND MORE BAAAAAAACON!",
	}
	return m, nil
}

func (fd fakedb) Update(meals []*Meal) error {
	if len(meals) > 0 {
		return nil
	}
	return errors.New("no meals")
}

// Retrieve will retrieve the calendar
func (fc fakecal) Retrieve(mealType string) ([]*Meal, error) {
	if mealType == "error" {
		return []*Meal{}, errors.New("you asked for one")
	}
	m := Meal{
		Date:     time.Now(),
		ID:       "myid",
		MainDish: "BAAAAAAAACON!!!",
		Sides:    "RANCH AND MORE BAAAAAAACON!",
	}
	out := []*Meal{&m}
	return out, nil
}

func setup() {
	fc := fakecal{}
	fd := fakedb{}
	ac = AppContext{
		backend: fd,
		cal:     fc,
	}
}

func TestGetToday(t *testing.T) {
	setup()
	out := GetTodaysLunch()
	if out == "" {
		t.Error("expected one day of bacon, but got nil")
	}
	if !strings.Contains(out, "RANCH AND MORE") {
		t.Errorf("Should have contained details but they are missing: %s", out)
	}
}

func TestGetWeek(t *testing.T) {
	setup()
	out := GetWeek()
	if out == "" {
		t.Error("expected a bunch of bacon but got nil")
	}
	if !strings.Contains(out, "BAAAAAAAACON!!!") {
		t.Errorf("expected a bunch of bacon but got: %s", out)
	}
	if strings.Contains(out, "RANCH AND MORE") {
		t.Errorf("Shouldn't have any details but do: %s", out)
	}
}

func TestGetDay(t *testing.T) {
	setup()
	now := time.Now()
	if now.Weekday() != time.Saturday {
		// test something this week
		out := GetDay("friday")
		if out == "" {
			t.Error("should have gotten some bacon on friday")
		}
		if !strings.Contains(out, "friday") {
			t.Errorf("Should have friday in there: %s", out)
		}
		if !strings.Contains(out, "RANCH AND MORE") {
			t.Errorf("Should have details but doesnt: %s", out)
		}
	}

	out := GetDay("monday")
	if out == "" {
		t.Error("should have gotten some bacon on monday")
	}
	if !strings.Contains(out, "monday") {
		t.Errorf("Should have monday in there: %s", out)
	}
	if !strings.Contains(out, "RANCH AND MORE") {
		t.Errorf("Should have details but doesnt: %s", out)
	}

}
