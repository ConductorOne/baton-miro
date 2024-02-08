package miro

import (
	"fmt"
	"net/http"
)

type (
	queryParam struct {
		name  string
		value string
	}
	queryFunction func() queryParam
)

func WithLimit(limit int32) queryFunction {
	return func() queryParam {
		return queryParam{
			name:  "limit",
			value: fmt.Sprint(limit),
		}
	}
}

func WithCursor(cursor string) queryFunction {
	return func() queryParam {
		return queryParam{
			name:  "cursor",
			value: cursor,
		}
	}
}

func addQueryParams(req *http.Request, queries ...queryFunction) *http.Request {
	q := req.URL.Query()
	for _, query := range queries {
		param := query()
		q.Add(param.name, param.value)
	}
	req.URL.RawQuery = q.Encode()

	return req
}
