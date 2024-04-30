package listenbrainz

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

const baseURL = "https://api.listenbrainz.org"
const statsExtraItems = 10
const interval = "week" // Allowed values: this_week, this_month, this_year, week, month, quarter, year, half_yearly, all_time

var client = http.DefaultClient

type Client struct {
	client    *http.Client
	userToken string
}

func NewClient(userToken string) *Client {
	return &Client{
		client:    http.DefaultClient,
		userToken: userToken,
	}
}

func (c *Client) FetchTopArtists(count int) ([]Artist, error) {
	stats, err := c.fetchArtistStats(count + statsExtraItems)
	if err != nil {
		return nil, fmt.Errorf("fetching artist stats failed: %w", err)
	}

	statsWithData := make([]StatsArtist, 0)
	mbids := make([]string, 0)

	for _, artist := range stats.Payload.Artists {
		if artist.MBID != "" && artist.Name != "" {
			mbids = append(mbids, artist.MBID)
			statsWithData = append(statsWithData, artist)
		}
		if len(mbids) == count {
			break
		}
	}

	metaDatas, err := c.fetchArtistMetadata(mbids)
	if err != nil {
		return nil, fmt.Errorf("fetching artist metadata failed: %w", err)
	}

	output := make([]Artist, 0)
	for _, stat := range statsWithData {
		metadata := metaDatas[stat.MBID]

		sort.Slice(
			metadata.Tag.Artist,
			func(a, b int) bool {
				return metadata.Tag.Artist[a].Count > metadata.Tag.Artist[b].Count
			},
		)

		tags := make([]string, 0)
		for _, tag := range metadata.Tag.Artist {
			tags = append(tags, tag.Tag)
		}

		output = append(output, Artist{
			IBID:        metadata.ArtistMBID,
			Name:        metadata.Name,
			URL:         fmt.Sprintf("https://listenbrainz.org/artist/%s", metadata.ArtistMBID),
			ListenCount: stat.ListenCount,
			Tags:        tags,
		})
	}

	return output, nil

}

func (c *Client) fetchArtistStats(count int) (*StatsArtistsResponse, error) {
	url := fmt.Sprintf("%s/1/stats/user/saaste/artists?count=%d&range=%s", baseURL, count, interval)
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

	var stats StatsArtistsResponse
	err = json.Unmarshal(bytes, &stats)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling response body failed: %w", err)
	}

	return &stats, nil
}

func (c *Client) fetchArtistMetadata(mbids []string) (map[string]ArtistMetaData, error) {
	mbidsParam := strings.Join(mbids, ",")
	url := fmt.Sprintf("%s/1/metadata/artist?artist_mbids=%s&inc=artist+tag", baseURL, mbidsParam)
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

	var responseData []ArtistMetaData
	err = json.Unmarshal(bytes, &responseData)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling response body failed: %w", err)
	}

	output := make(map[string]ArtistMetaData)
	for _, metadata := range responseData {
		output[metadata.ArtistMBID] = metadata
	}

	return output, nil
}

func (c *Client) makeRequest(method string, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating a request failed: %w", err)
	}
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", c.userToken))
	return request, nil
}
