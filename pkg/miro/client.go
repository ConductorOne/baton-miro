package miro

import (
	"context"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

const (
	BaseUrl     = "https://api.miro.com"
	ScimBaseUrl = "https://miro.com/api/v1/scim"
)

type Client struct {
	httpClient *uhttp.BaseHttpClient
	baseUrl    string
}

func New(httpClient *http.Client) *Client {
	return &Client{
		httpClient: uhttp.NewBaseHttpClient(httpClient),
		baseUrl:    BaseUrl,
	}
}

func (c *Client) doRequest(
	ctx context.Context,
	url *url.URL,
	method string,
	res interface{},
	body interface{},
) (*http.Response, error) {
	l := ctxzap.Extract(ctx)
	var reqOptions []uhttp.RequestOption
	if body != nil {
		reqOptions = append(reqOptions, uhttp.WithJSONBody(body))
	}

	req, err := c.httpClient.NewRequest(ctx, method, url, reqOptions...)
	if err != nil {
		return nil, err
	}

	var doOptions []uhttp.DoOption

	if res != nil {
		doOptions = append(doOptions, uhttp.WithJSONResponse(res))
	}

	resp, err := c.httpClient.Do(req, doOptions...)
	if err != nil {
		l.Error("miro-connector: failed to execute request", zap.Error(err))
		return nil, err
	}

	return resp, nil
}
