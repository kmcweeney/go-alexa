package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kmcweeney/go-alexa/alexa"
)

const (
	tFormat string = "2006-01-02"
	region  string = "us-east-1"
)

func DispatchIntents(request alexa.Request) alexa.Response {
	var response alexa.Response
	if request.Body.Type == "LaunchRequest" {
		return alexa.NewSimpleResponse("welcome", GetTodaysLunch())
	}
	switch request.Body.Intent.Name {
	case "today":
		response = alexa.NewSimpleResponse("today", GetTodaysLunch())
	case "todayLunch":
		response = alexa.NewSimpleResponse("today", GetTodaysLunch())
	case "week":
		response = alexa.NewSimpleResponse("week", GetWeek())
	case "weekLunch":
		response = alexa.NewSimpleResponse("week", GetWeek())
	case "updateDB":
		fmt.Println("handeling case updateDB")
		meals, err := buildMeals("https://www.parkwayschools.net/site/handlers/icalfeed.ashx?MIID=4134", "lunch")
		if err != nil {
			fmt.Println("Error building meal map from ical: ", err)
			response = alexa.NewSimpleResponse("error", "couldn't build db")
		}
		err = UpdateDB(meals)
		if err != nil {
			fmt.Println("Couldn't load data:", err)
			response = alexa.NewSimpleResponse("error", "couldn't build db")
		}
		response = alexa.NewSimpleResponse("worked", "db updated")
	case "firstInsert":
		fmt.Println("handeling case updateDB")
		meals, err := buildMeals("https://www.parkwayschools.net/site/handlers/icalfeed.ashx?MIID=4134", "lunch")
		if err != nil {
			fmt.Println("Error building meal map from ical: ", err)
			response = alexa.NewSimpleResponse("error", "couldn't build db")
		}
		err = firstInsert(meals)
		if err != nil {
			fmt.Println("Couldn't load data:", err)
			response = alexa.NewSimpleResponse("error", "couldn't build db")
		}
		response = alexa.NewSimpleResponse("worked", "db updated")
	case alexa.FallbackIntent:
		response = alexa.NewSimpleResponse("today", GetTodaysLunch())
	case alexa.HelpIntent:
		response = alexa.NewSimpleResponse("help", "ask for today's lunch")
	}

	return response
}

func Handler(request alexa.Request) (alexa.Response, error) {
	return DispatchIntents(request), nil
}

func main() {
	lambda.Start(Handler)
}
