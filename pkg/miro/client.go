package miro

import (
	"net/http"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

const BaseUrl = "https://api.miro.com"

type Client struct {
	uhttp.BaseHttpClient

	baseUrl string
}

func New(accessToken string, httpClient *http.Client) *Client {
	return &Client{
		BaseHttpClient: *uhttp.NewBaseHttpClient(httpClient),
		baseUrl:        BaseUrl,
	}
}
