package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	ics "github.com/PuloV/ics-golang"
)

var lunches map[time.Time]string

func getLunch(date time.Time) string {
	if lunches[date] == "" {
		buildLunchMap()
	}
	if lunches[date] == "" {
		return "nothing"
	}
	return lunches[date]
}

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
		out = out + fmt.Sprintf("%s has %s,", startDay.Weekday(), getLunch(startDay))
		startDay = startDay.AddDate(0, 0, 1)
	}
	return out

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

func buildLunchMap() {

	if lunches == nil {
		lunches = make(map[time.Time]string)
	}

	parser := ics.New()
	ics.FilePath = "/tmp/"
	pChan := parser.GetInputChan()
	pChan <- "https://www.parkwayschools.net/site/handlers/icalfeed.ashx?MIID=4134"
	outputChan := parser.GetOutputChan()
	go func() {
		for event := range outputChan {
			date := time.Date(event.GetStart().Year(), event.GetStart().Month(), event.GetStart().Day(), 0, 0, 0, 0, time.Now().Location())
			lunches[date] = event.GetSummary()
			//fmt.Printf("%v:%s\n", date, event.GetSummary())
		}
	}()
	parser.Wait()
	errs, _ := parser.GetErrors()
	for _, err := range errs {
		fmt.Println("parseError:", err.Error())
	}
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	fmt.Printf("build map %s\n", lunches[today])
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
