package miro

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

const (
	BaseUrl     = "https://api.miro.com"
	ScimBaseUrl = "https://miro.com/api/v1/scim/"
)

// Client is the Miro client.
type Client struct {
	httpClient *uhttp.BaseHttpClient
	scimClient *uhttp.BaseHttpClient
}

// New creates a new Miro client.
func New(httpClient *http.Client, scimClient *http.Client) *Client {
	c := &Client{
		httpClient: uhttp.NewBaseHttpClient(httpClient),
	}

	if scimClient != nil {
		c.scimClient = uhttp.NewBaseHttpClient(scimClient)
	}

	return c
}

// doRequest executes a request to the Miro API.
func (c *Client) doRequest(
	ctx context.Context,
	endpointUrl string,
	method string,
	res interface{},
	body interface{},
	opts ...ReqOpt,
) (http.Header, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var reqOptions []uhttp.RequestOption
	if body != nil {
		reqOptions = append(reqOptions, uhttp.WithJSONBody(body))
	}

	baseUrl, err := url.Parse(BaseUrl)
	if err != nil {
		return nil, nil, err
	}

	endpointParsed, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, nil, err
	}

	urlAddress := baseUrl.ResolveReference(endpointParsed)

	req, err := c.httpClient.NewRequest(ctx, method, urlAddress, reqOptions...)
	if err != nil {
		return nil, nil, err
	}

	for _, opt := range opts {
		req = opt(req)
	}

	var doOptions []uhttp.DoOption
	var ratelimitData v2.RateLimitDescription

	if res != nil {
		doOptions = append(doOptions, uhttp.WithJSONResponse(res))
	}
	doOptions = append(doOptions, uhttp.WithRatelimitData(&ratelimitData))

	resp, err := c.httpClient.Do(req, doOptions...)
	if err != nil {
		l.Error("miro-connector: failed to execute request", zap.Error(err))
		return nil, nil, err
	}
	defer resp.Body.Close()

	annos := annotations.Annotations{}
	annos.WithRateLimiting(&ratelimitData)

	return resp.Header, annos, nil
}

// doScimRequest executes a request to the Miro SCIM API.
func (c *Client) doScimRequest(
	ctx context.Context,
	endpointUrl string,
	method string,
	res interface{},
	body interface{},
	opts ...ReqOpt,
) (http.Header, annotations.Annotations, error) {
	if c.scimClient == nil {
		return nil, nil, fmt.Errorf("SCIM client not configured: SCIM access token is required for this operation")
	}

	l := ctxzap.Extract(ctx)
	var reqOptions []uhttp.RequestOption
	if body != nil {
		reqOptions = append(reqOptions, uhttp.WithJSONBody(body))
	}

	baseUrl, err := url.Parse(ScimBaseUrl)
	if err != nil {
		return nil, nil, err
	}

	endpointParsed, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, nil, err
	}

	urlAddress := baseUrl.ResolveReference(endpointParsed)
	req, err := c.scimClient.NewRequest(ctx, method, urlAddress, reqOptions...)
	if err != nil {
		return nil, nil, err
	}

	for _, opt := range opts {
		req = opt(req)
	}

	var doOptions []uhttp.DoOption
	var ratelimitData v2.RateLimitDescription

	if res != nil {
		doOptions = append(doOptions, uhttp.WithJSONResponse(res))
	}
	doOptions = append(doOptions, uhttp.WithRatelimitData(&ratelimitData))

	resp, err := c.scimClient.Do(req, doOptions...)
	if err != nil {
		l.Error("miro-connector: failed to execute request", zap.Error(err))
		return nil, nil, err
	}
	defer resp.Body.Close()

	annos := annotations.Annotations{}
	annos.WithRateLimiting(&ratelimitData)

	return resp.Header, annos, nil
}
