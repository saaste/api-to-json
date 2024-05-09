package webmentions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

	filteredMentions := make([]Webmention, 0)
	for _, mention := range resp.Webmentions {
		// Ignore local links
		if strings.Contains(mention.Source, "saaste.net") {
			c.deleteLocalMention(domain, mention)
			continue
		}
		filteredMentions = append(filteredMentions, mention)
	}
	resp.Webmentions = filteredMentions

	return &resp, nil
}

func (c *Client) deleteLocalMention(domain string, mention Webmention) error {
	url := fmt.Sprintf("%s/webmention/%s/%s?source=%s&target=%s", baseURL, url.PathEscape(domain), url.PathEscape(c.token), url.QueryEscape(mention.Source), url.QueryEscape(mention.Target))
	request, err := c.makeRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to make delete request: %w", err)
	}

	_, err = c.client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send delete request: %w", err)
	}

	return nil
}

func (c *Client) makeRequest(method string, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating a request failed: %w", err)
	}
	return request, nil
}
