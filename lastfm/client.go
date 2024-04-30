package lastfm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	ns_url "net/url"
)

const baseURL = "https://ws.audioscrobbler.com/2.0/"

var client = http.DefaultClient

type Client struct {
	client *http.Client
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{
		client: http.DefaultClient,
		apiKey: apiKey,
	}
}

func (c *Client) GetArtistInfoByMBID(mbid string) (*ArtistInfo, error) {
	resp, err := c.fetchArtistInfo("", mbid)
	if err != nil {
		return nil, err
	}

	return &resp.ArtistInfo, nil
}

func (c *Client) fetchArtistInfo(artistName string, MBID string) (*ArtistInfoResponse, error) {
	var url string
	if MBID != "" {
		url = fmt.Sprintf("%s?method=artist.getinfo&mbid=%s&api_key=%s&format=json", baseURL, MBID, c.apiKey)
	} else {
		url = fmt.Sprintf("%s?method=artist.getinfo&artist=%s&api_key=%s&format=json", baseURL, ns_url.QueryEscape(artistName), c.apiKey)
	}

	fmt.Println(url)

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

	// TODO Consider doing something with rate limiting headers

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading the response body failed: %w", err)
	}

	var stats ArtistInfoResponse
	err = json.Unmarshal(bytes, &stats)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling response body failed: %w", err)
	}

	return &stats, nil
}

func (c *Client) makeRequest(method string, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating a request failed: %w", err)
	}
	return request, nil
}
