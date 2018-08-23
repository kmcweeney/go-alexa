package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kmcweeney/go-alexa/alexa"
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
