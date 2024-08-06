package listenbrainz

// Stats
type StatsArtistsResponse struct {
	Payload StatsPayload `json:"payload"`
}
type StatsPayload struct {
	Artists          []StatsArtist `json:"artists"`
	Count            int64         `json:"count"`
	TotalArtistCount int64         `json:"total_artist_count"`
	Range            string        `json:"range"`
	LastUpdated      int64         `json:"last_updated"`
	UserID           string        `json:"user_id"`
	FromTS           int64         `json:"from_ts"`
	ToTS             int64         `json:"to_ts"`
}
type StatsArtist struct {
	MBID        string `json:"artist_mbid"`
	Name        string `json:"artist_name"`
	ListenCount int64  `json:"listen_count"`
}

// Metadata
type ArtistMetaData struct {
	ArtistMBID string            `json:"artist_mbid"`
	MBID       string            `json:"mbid"`
	Gender     string            `json:"gender"`
	Name       string            `json:"name"`
	Rels       MetaDataArtistRel `json:"rels"`
	Tag        MetaDataTag       `json:"tag"`
	Type       string            `json:"type"`
	Area       string            `json:"area"`
	Beginyear  int               `json:"begin_year"`
}
type MetaDataArtistRel struct {
	OfficialHomepage    string `json:"official_homepage"`
	YouTube             string `json:"youtube"`
	PurchaseForDownload string `json:"purchase for download"`
	Wikidata            string `json:"wikidata"`
	FreeStreaming       string `json:"free streaming"`
	SocialNetwork       string `json:"social network"`
	Lyrics              string `json:"lyrics"`
}
type MetaDataTag struct {
	Recording []RecordingTag `json:"recording"`
	Artist    []ArtistTag    `json:"Artist"`
}
type RecordingTag struct {
	GenreMBID string `json:"genre_mbid"`
	Tag       string `json:"tag"`
	Count     int64  `json:"count"`
}
type ArtistTag struct {
	ArtistMBID string `json:"artist_mbid"`
	GenreMBID  string `json:"genre_mbid"`
	Tag        string `json:"tag"`
	Count      int64  `json:"count"`
}

type TopArtistsResult struct {
	TopArtists       []Artist
	TotalArtistCount int64
}

// Output
type Artist struct {
	IBID        string
	Name        string
	URL         string
	Genres      []string
	ListenCount int64
	Tags        []string
}
