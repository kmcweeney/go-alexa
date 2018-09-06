package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

//Backend interface for getting meals
type Backend interface {
	Get(date time.Time, mealType string, details bool) (Meal, error)
	Update(meals []*Meal) error
}

//Dynamo implementation of the backend
type dynamo struct{}

//Meal is the representation of a meal
type Meal struct {
	ID       string    `json:"id"`
	Date     time.Time `json:"date"`
	MealType string    `json:"mealType"`
	MainDish string    `json:"mainDish"`
	Sides    string    `json:"sides"`
}

// Get returns the meal with the given type on the date
// details are returned if the details input is true
func (d dynamo) Get(date time.Time, mealType string, details bool) (Meal, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		fmt.Println("Error creating AWS session")
		return Meal{}, err
	}
	svc := dynamodb.New(sess)
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("meal"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(fmt.Sprintf("%s-%s", date.Format(time.RFC3339), mealType)),
			},
		},
	})
	if err != nil {
		fmt.Println("error getting meal: ", err)
		return Meal{}, err
	}
	m := Meal{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &m)
	if err != nil {
		fmt.Println("error while unmarshalling: ", err)
		return Meal{}, err
	}
	return m, nil
}

// Update the db with the input slice of Meals
func (d dynamo) Update(meals []*Meal) error {
	fmt.Println("Updating DB")
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		fmt.Println("Error creating AWS session: ", err)
		return err
	}
	svc := dynamodb.New(sess)
	for _, meal := range meals {
		exists, err := ac.backend.Get(meal.Date, meal.MealType, false)
		if err != nil {
			fmt.Println("Problem checking if the meal already exists: ", err)
			return err
		}
		if exists.MainDish == "" {
			av, err := dynamodbattribute.MarshalMap(meal)
			if err != nil {
				fmt.Println("Got error marshalling the meal:", err)
				return err
			}
			input := &dynamodb.PutItemInput{Item: av, TableName: aws.String("meal")}
			_, err = svc.PutItem(input)
			if err != nil {
				fmt.Println("Error putting item in db: ", err)
				return err
			}
			fmt.Printf("successfull added item %v\n", meal)
		}
	}
	return nil
}
