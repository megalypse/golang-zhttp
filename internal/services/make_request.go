package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	utils "github.com/megalypse/zhttp/internal"
	"github.com/megalypse/zhttp/models"
)

func MakeRequest[Response any, Request any](method string, request models.ZRequest[Request]) models.ZResponse[Response] {
	responseHolder := new(Response)
	client := http.Client{}

	bodyBuffer, marshalErr := json.Marshal(request.Body)

	if marshalErr != nil {
		return utils.MakeFailResponse[Response](marshalErr.Error(), nil)
	}

	httpRequest, _ := http.NewRequest(
		method,
		utils.ParseUrl(request),
		bytes.NewBuffer(bodyBuffer),
	)

	for key, value := range request.Headers {
		httpRequest.Header.Set(key, value)
	}

	httpResponse, _ := client.Do(httpRequest)

	responseBuffer, readErr := io.ReadAll(httpResponse.Body)

	if readErr != nil {
		return utils.MakeFailResponse[Response](marshalErr.Error(), nil)
	}

	unmarshalError := json.Unmarshal(responseBuffer, &responseHolder)

	if unmarshalError != nil {
		return utils.MakeFailResponse[Response](unmarshalError.Error(), nil)
	}

	return models.ZResponse[Response]{
		Content:   responseHolder,
		Response:  httpResponse,
		IsSuccess: true,
	}
}
