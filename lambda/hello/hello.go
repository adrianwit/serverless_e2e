package main

import "github.com/aws/aws-lambda-go/lambda"

type Person struct {
	Name string
}

func handle(person *Person) (string, error) {
	if person.Name != "" {
		return "Hello " + person.Name, nil
	}
	return "Hello World", nil
}

func main() {
	lambda.Start(handle)
}
