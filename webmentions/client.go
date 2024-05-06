package webmentions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const baseURL = "https://webmentions.saaste.net"

var client = http.DefaultClient

type Client struct {
	client *http.Client
	token  string
}

func NewClient(token string) *Client {
	return &Client{
		client: http.DefaultClient,
		token:  token,
	}
}

func (c *Client) GetWebmentions(domain string) (*Response, error) {
	url := fmt.Sprintf("%s/webmention/%s/%s", baseURL, url.PathEscape(domain), url.PathEscape(c.token))
	request, err := c.makeRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("sending the request failed: %w", err)
	} else if response.StatusCode != http.StatusOK {
		bytes, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("sending the request returned an error with status %d: %s", response.StatusCode, string(bytes))
	}

	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading the response body failed: %w", err)
	}

	var resp Response
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling response body failed: %w", err)
	}

	return &resp, nil
}

func (c *Client) makeRequest(method string, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating a request failed: %w", err)
	}
	return request, nil
}
