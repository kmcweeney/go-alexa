package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/PuloV/ics-golang"
	"github.com/kmcweeney/go-alexa/alexa"
)

var lunches map[time.Time]string

func buildLunchMap() {
	if lunches == nil {
		lunches = make(map[time.Time]string)
	}
	parser := ics.New()

	parserChan := parser.GetInputChan()
	parserChan <- "https://www.parkwayschools.net/site/handlers/icalfeed.ashx?MIID=4134"

	outputChan := parser.GetOutputChan()
	go func() {
		for event := range outputChan {
			date := time.Date(event.GetStart().Year(), event.GetStart().Month(), event.GetStart().Day(), 0, 0, 0, 0, event.GetStart().Location())
			lunches[date] = event.GetSummary()
			fmt.Printf("%v:%s\n", event.GetStart(), event.GetSummary())
		}
	}()
	parser.Wait()
}

func DispatchIntents(request alexa.Request) alexa.Response {
	var response alexa.Response
	switch request.Body.Intent.Name {
	case "today":
		response = handleToday(request)
	case alexa.HelpIntent:
		return alexa.NewSimpleResponse("help", "ask for today's lunch")
	}

	return response
}

func handleToday(request alexa.Request) alexa.Response {
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	if lunches[today] == "" {
		buildLunchMap()
	}
	if lunches[today] == "" {
		return alexa.NewSimpleResponse("error", "sorry today is not on the calendar")
	}

	response := fmt.Sprintf("Today the menu is %s", lunches[today])
	return alexa.NewSimpleResponse("today", response)
}

func Handler(request alexa.Request) (alexa.Response, error) {
	return DispatchIntents(request), nil
}

func main() {
	buildLunchMap()
	lambda.Start(Handler)
}
