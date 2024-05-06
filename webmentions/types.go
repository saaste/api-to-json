package webmentions

type Response struct {
	Status      string       `json:"status"`
	Webmentions []Webmention `json:"json"`
}

type Webmention struct {
	Author    Author `json:"author"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	Published string `json:"published"`
	URL       string `json:"url"`
	Type      string `json:"type"`
	Source    string `json:"source"`
	Target    string `json:"target"`
}

type Author struct {
	Name       string `json:"name"`
	PictureURL string `json:"picture"`
}
