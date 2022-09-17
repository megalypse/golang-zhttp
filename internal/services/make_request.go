package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/megalypse/zhttp/internal/models"
)

func makeRequest[Response any, Request any](method string, request models.ZRequest[Request]) models.ZResponse[Response] {
	responseHolder := new(Response)
	client := http.Client{}

	bodyBuffer, marshalErr := json.Marshal(request.Body)

	if marshalErr != nil {
		return models.MakeFailResponse[Response](marshalErr.Error(), nil)
	}

	httpRequest, _ := http.NewRequest(
		method,
		parseUrl(request),
		bytes.NewBuffer(bodyBuffer),
	)

	for _, header := range request.Headers {
		httpRequest.Header.Set(header.Key, header.Value)
	}

	httpResponse, _ := client.Do(httpRequest)

	responseBuffer, readErr := io.ReadAll(httpResponse.Body)

	if readErr != nil {
		return models.MakeFailResponse[Response](marshalErr.Error(), nil)
	}

	unmarshalError := json.Unmarshal(responseBuffer, &responseHolder)

	if unmarshalError != nil {
		return models.MakeFailResponse[Response](unmarshalError.Error(), nil)
	}

	return models.ZResponse[Response]{
		Content:   responseHolder,
		Response:  httpResponse,
		IsSuccess: true,
	}
}

func parseUrl[T any](request models.ZRequest[T]) string {
	url := request.Url
	urlParams := request.UrlParams
	queryParams := request.QueryParams

	for _, param := range urlParams {
		paramInterpolation := fmt.Sprintf("{%v}", param.Key)
		strings.ReplaceAll(url, paramInterpolation, param.Value)
	}

	urlLastIndex := len(url) - 1

	if string(url[urlLastIndex]) == "/" {
		url = url[:urlLastIndex]
	}

	url += "?"

	for i, v := range queryParams {
		var param string

		if i > 0 {
			param += "&"
		}

		param += fmt.Sprintf("%v=%v", v.Key, v.Value)

		url += param
	}

	return url
}

// func generateRequestUrl[T any](request models.ZRequest[T]) string {
// 	if request.Url != nil {
// 		return *request.Url
// 	}

// 	context := request.Context
// 	uri := request.Uri

// 	contextLastIndex := len(context) - 1

// 	if string(context[contextLastIndex]) == "/" {
// 		context = context[:contextLastIndex]
// 	}

// 	if string(uri[0]) != "/" {
// 		uri = "/" + uri
// 	}

// 	return context + uri
// }
