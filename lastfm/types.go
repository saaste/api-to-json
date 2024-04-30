package lastfm

type ImageSize string

const (
	ImageSizeSmall      ImageSize = "small"
	ImageSizeMedium     ImageSize = "medium"
	ImageSizeLarge      ImageSize = "large"
	ImageSizeExtraLarge ImageSize = "extralarge"
	ImageSizeMega       ImageSize = "mega"
	ImageSizeUnknown    ImageSize = ""
)

type ArtistInfoResponse struct {
	ArtistInfo ArtistInfo `json:"artist"`
}

type ArtistInfo struct {
	Name       string            `json:"name"`
	MBID       string            `json:"mbid"`
	URL        string            `json:"url"`
	Images     []ArtistInfoImage `json:"image"`
	Streamable string            `json:"streamable"`
	OnTour     string            `json:"ontour"`
	Stats      ArtistInfoStats   `json:"stats"`
	Similar    ArtistInfoSimilar `json:"similar"`
	Tags       ArtistInfoTags    `json:"tags"`
	Bio        ArtistInfoBio     `json:"bio"`
}

type ArtistInfoImage struct {
	URL  string    `json:"#text"`
	Size ImageSize `json:"size"`
}

type ArtistInfoStats struct {
	Listeners string `json:"listeners"`
	PlayCount string `json:"playcount"`
}

type ArtistInfoSimilar struct {
	Artists []ArtistInfoSimilarData `json:"artist"`
}

type ArtistInfoSimilarData struct {
	Name   string `json:"string"`
	URL    string `json:"url"`
	Images []ArtistInfoImage
}

type ArtistInfoTags struct {
	Tag []ArtistInfoTag `json:"tag"`
}

type ArtistInfoTag struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ArtistInfoBio struct {
	Links     ArtistInfoLinks `json:"links"`
	Published string          `json:"published"`
	Summary   string          `json:"summary"`
	Content   string          `json:"content"`
}

type ArtistInfoLinks struct {
	Link ArtistInfoLink `json:"link"`
}

type ArtistInfoLink struct {
	Text string `json:"#text"`
	Rel  string `json:"rel"`
	Href string `json:"href"`
}
