package lambdautils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/jaydp17/movie-ticket-watcher/pkg/httperror"
	"net/http"
)

type Response = events.APIGatewayProxyResponse

var commonHeaders = map[string]string{
	"Content-Type":                     "application/json",
	"Access-Control-Allow-Origin":      "*",
}

// ToResponse takes any type of object, could be data or any error and tries to convert it to a proper APIGatewayProxyResponse
// It also supports httperror, where there's a StatusCode along with the error
// If it's not a httperror and a regular error, it'll send 500 as the status, along with the error message
func ToResponse(result interface{}) (Response, error) {
	errorWithCode, ok := result.(httperror.Error)
	// Http Error with statusCode
	if ok {
		jsonBody, toJsonErr := toJsonError(errorWithCode.Error())
		if toJsonErr != nil {
			return InternalServerErrorResp(), toJsonErr
		}
		resp := Response{
			StatusCode: errorWithCode.Code(),
			Body:       jsonBody,
			Headers:    commonHeaders,
		}
		return resp, nil
	}

	regularError, ok := result.(error)
	// regular error
	if ok {
		jsonBody, marshalError := toJsonError(regularError.Error())
		if marshalError != nil {
			return InternalServerErrorResp(), marshalError
		}
		resp := Response{
			StatusCode: http.StatusInternalServerError,
			Body:       jsonBody,
			Headers:    commonHeaders,
		}
		return resp, nil
	}

	jsonBody, marshalError := toJson(result)
	if marshalError != nil {
		return InternalServerErrorResp(), marshalError
	}
	resp := Response{
		StatusCode: http.StatusOK,
		Body:       jsonBody,
		Headers:    commonHeaders,
	}
	return resp, nil
}

func toJson(data interface{}) (string, error) {
	body, marshalError := json.Marshal(data)
	if marshalError != nil {
		return "", marshalError
	}

	var jsonBuf bytes.Buffer
	json.HTMLEscape(&jsonBuf, body)

	return jsonBuf.String(), nil
}

func toJsonError(errMsg string) (string, error) {
	return toJson(struct{ Error string `json:"error"` }{Error: errMsg})
}

func InternalServerErrorResp() Response {
	return Response{
		StatusCode: http.StatusInternalServerError,
		Body:       fmt.Sprintf("{\"error\": \"%s\"}", http.StatusText(http.StatusInternalServerError)),
		Headers:    commonHeaders,
	}
}
