package main

import (
	"os"
	"log"
	"bytes"
	"context"
	"net/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, event events.DynamoDBEvent) error {
	for _, record := range event.Records {
		for name, value := range record.Change.NewImage {
			if value.DataType() == events.DataTypeString {
				err := save(name, value.String())
				if err != nil {
					log.Print(err)
					return err
				}
			}
		}
	}
	return nil
}

func save(name string, value string) error {
	jsonStr := `{"` + name + `":"` + value + `"}`

	req, err := http.NewRequest(
		"POST",
		os.Getenv("DOMAIN") + "/" + os.Getenv("ES_INDEX_NAME") + "/" + os.Getenv("ES_TYPE_NAME") + "/",
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		return err
	}
    req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
