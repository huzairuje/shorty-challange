package tiny_url

type Data struct {
	ShortCode     string `json:"shortcode"`
	Url           string `json:"url"`
	StartDate     string `json:"startDate"`
	LastSeenDate  string `json:"lastSeenDate"`
	RedirectCount int    `json:"redirectCount"`
}
