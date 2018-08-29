package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	ics "github.com/PuloV/ics-golang"
)

func getFromDb(date time.Time, mealType string, details bool) (string, error) {
	m, err := Get(date, mealType, details)
	if err != nil {
		fmt.Println("Problem with DB: ", err)
		return "", err
	}
	if m.MainDish == "" {
		meals, err := buildMeals("https://www.parkwayschools.net/site/handlers/icalfeed.ashx?MIID=4134", mealType)
		if err != nil {
			fmt.Println("Error building meal map from ical: ", err)
			return "", err
		}
		err = UpdateDB(meals)
		if err != nil {
			fmt.Println("Couldn't load data:", err)
			return "", err
		}
	}
	m, err = Get(date, mealType, details)
	if err != nil {
		fmt.Println("Problem with DB: ", err)
		return "", err
	}
	if m.MainDish == "" {
		return "nothing", nil
	}
	if details {
		return fmt.Sprintf("%s, with %s", m.MainDish, m.Sides), nil
	}
	return m.MainDish, nil
}

// GetTodaysLunch gets today's lunch
func GetTodaysLunch() string {
	fmt.Println("Getting today")
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	result, err := getFromDb(today, "lunch", true)
	if err != nil {
		return "Sorry there was an error retrieving the lunch schedule"
	}
	return fmt.Sprintf("Today the menu is %s", result)
}

func GetTomorrow() string {
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	tomorrow := today.AddDate(0, 0, 1)
	result, err := getFromDb(tomorrow, "lunch", true)
	if err != nil {
		return "Sorry there was an error retrieving the lunch schedule"
	}
	return fmt.Sprintf("Tomorrow the menu is %s", result)
}

// GetWeek gets the whole week worth of lunches
func GetWeek() string {
	fmt.Println("getting the week")
	startDay := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	if startDay.Weekday() == 0 {
		startDay = startDay.AddDate(0, 0, 1)
	}
	if startDay.Weekday() == 5 {
		startDay = startDay.AddDate(0, 0, 2)
	}
	out := ""
	for i := startDay.Weekday(); i <= 5; i++ {
		day, err := getFromDb(startDay, "lunch", false)
		if err != nil {
			day = "trouble"
		}
		out = out + fmt.Sprintf("%s has %s,", startDay.Weekday(), day)

		startDay = startDay.AddDate(0, 0, 1)
	}
	return out

}

// GetDay gets the whole week worth of lunches
func GetDay(day string) string {
	fmt.Printf("getting the day: %v\n", day)
	d := weekDayMap[strings.ToLower(day)]
	if d == 0 {
		return fmt.Sprintf("Please specify a weekday, there is no lunch on %s", day)
	}

	dateWanted := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	cDay := dateWanted.Weekday()
	// if the day requested has already passed this week, get it for next week
	if d < cDay {
		// the day requested has already passed this week, get it for next week
		dateWanted.AddDate(0, 0, 7-(int(d)-int(cDay)))
		fmt.Printf("setting the date to %s", dateWanted.String())
	} else {
		// the day requests is still this week, so just add the number of days needed
		dateWanted.AddDate(0, 0, int(cDay)-int(d))
		fmt.Printf("setting the date to %s", dateWanted.String())
	}
	result, err := getFromDb(dateWanted, "lunch", true)
	if err != nil {
		return "Sorry there was an error retrieving the lunch schedule"
	}
	return fmt.Sprintf("On %s the menu is %s", day, result)
}

func getCalendar() string {
	fmt.Println("getting the calendar by http get")
	timeout := time.Duration(5 * time.Second)
	client := http.Client{Timeout: timeout}
	out := ""
	req, err := http.NewRequest("GET", "https://www.parkwayschools.net/site/handlers/icalfeed.ashx?MIID=4134", nil)
	if err != nil {
		fmt.Printf("error creating request %v\n", err)
	}
	req.Header.Add("Content-Type", "text/calendar")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error getting the file %v\n", err)
	}
	fmt.Println("Response code", resp.StatusCode)
	fmt.Println("Response length", resp.ContentLength)

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("error getting body bytes %v\n", err)
		}
		out = string(bodyBts)
	}
	fmt.Println("exiting calendar http get")
	//fmt.Printf("output: %s\n", out)
	return out
}

func buildMeals(url string, mealType string) ([]*Meal, error) {
	var out []*Meal
	parser := ics.New()
	ics.FilePath = "/tmp/"
	pChan := parser.GetInputChan()
	pChan <- url
	outputChan := parser.GetOutputChan()
	var mutex = &sync.Mutex{}
	go func() {
		for event := range outputChan {
			m := Meal{MealType: mealType}
			date := time.Date(event.GetStart().Year(), event.GetStart().Month(), event.GetStart().Day(), 0, 0, 0, 0, time.Now().Location())
			m.Date = date
			m.MainDish = event.GetSummary()
			m.Sides = event.GetDescription()
			m.ID = fmt.Sprintf("%s-%s", date.Format(time.RFC3339), mealType)
			mutex.Lock()
			out = append(out, &m)
			mutex.Unlock()
		}
	}()
	parser.Wait()
	fmt.Println("meals: ", out)
	return out, nil
}
