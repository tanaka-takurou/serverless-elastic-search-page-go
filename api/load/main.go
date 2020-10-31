package main

import (
	"os"
	"fmt"
	"log"
	"context"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type APIResponse struct {
	Message  string   `json:"message"`
	List     []string `json:"list"`
}

type Response events.APIGatewayProxyResponse

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var jsonBytes []byte
	var err error
	d := make(map[string]string)
	json.Unmarshal([]byte(request.Body), &d)
	if v, ok := d["action"]; ok {
		switch v {
		case "search" :
			if w, ok := d["word"]; ok {
				list, e := search(w)
				if e != nil {
					err = e
				} else {
					jsonBytes, _ = json.Marshal(APIResponse{Message: "Success", List: list})
				}
			}
		}
	}
	log.Print(request.RequestContext.Identity.SourceIP)
	if err != nil {
		log.Print(err)
		jsonBytes, _ = json.Marshal(APIResponse{Message: fmt.Sprint(err), List: []string{}})
		return Response{
			StatusCode: http.StatusInternalServerError,
			Body: string(jsonBytes),
		}, nil
	}
	return Response {
		StatusCode: http.StatusOK,
		Body: string(jsonBytes),
	}, nil
}

func search(word string)([]string, error) {
	var res []string
	rawResponse, err := getRawResponse(os.Getenv("DOMAIN") + "/" + os.Getenv("ES_INDEX_NAME") + "/" + os.Getenv("ES_TYPE_NAME") + "/_search?q=" + word)
	if err != nil {
		return res, err
	}
	rawResponse_ := rawResponse.(map[string]interface{})["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(rawResponse_) < 1 {
		return res, fmt.Errorf("Error: %s", "No hits")
	}
	for _, v := range rawResponse_ {
		hitText := v.(map[string]interface{})["_source"].(map[string]interface{})["text"].(string)
		if len(hitText) > 0 {
			res = append(res, hitText)
		}
	}
	return res, nil
}

func getRawResponse(url string)( interface{}, error) {
	var d interface{}
	res, err := http.Get(url)
	if err != nil {
		return d, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return d, err
	}

	err = json.Unmarshal(body, &d)
	if err != nil {
		return d, err
	}
	return d, nil
}

func main() {
	lambda.Start(HandleRequest)
}
