package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kmcweeney/go-alexa/alexa"
)

const (
	tFormat string = "2006-01-02"
	region  string = "us-east-1"
)

var weekDayMap = map[string]time.Weekday{
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
}

//AppContext contains pointers to the implementations
type AppContext struct {
	cal     Calendar
	backend Backend
}

var ac AppContext

//DispatchIntents answer some alexa questions
func DispatchIntents(request alexa.Request) alexa.Response {
	var response alexa.Response
	if request.Body.Type == "LaunchRequest" {
		return alexa.NewSimpleResponse("welcome", GetTodaysLunch())
	}
	fmt.Printf("Intent name: %s\n", request.Body.Intent.Name)
	switch request.Body.Intent.Name {
	case "todayLunch":
		fmt.Println("handeling case today intent")
		response = alexa.NewSimpleResponse("today", GetTodaysLunch())
	case "weekLunch":
		fmt.Println("handeling case week intent")
		response = alexa.NewSimpleResponse("week", GetWeek())
	case "tomorrow":
		fmt.Println("handeling case tomorrow intent")
		response = alexa.NewSimpleResponse("tomorrow", GetTomorrow())
	case "dayOfWeek":
		fmt.Println("dayOfWeek intent")
		day := request.Body.Intent.Slots["day"]
		fmt.Printf("Looking up lunch for %s\n", day)
		response = alexa.NewSimpleResponse("day", GetDay(day.Value))
	case alexa.FallbackIntent:
		fmt.Println("handeling case fallback")
		response = alexa.NewSimpleResponse("today", GetTodaysLunch())
	case alexa.HelpIntent:
		fmt.Println("handeling case help")
		response = alexa.NewSimpleResponse("help", "ask for today's lunch")
	}

	return response
}

//Handler the handler
func Handler(request alexa.Request) (alexa.Response, error) {
	return DispatchIntents(request), nil
}

func setup() {
	cals := make(map[string]string)
	cals["lunch"] = "https://www.parkwayschools.net/site/handlers/icalfeed.ashx?MIID=4134"
	be := dynamo{}
	cal := parkway{
		cals: cals,
	}
	ac = AppContext{
		backend: be,
		cal:     cal,
	}
}

func main() {
	setup()
	lambda.Start(Handler)
}
