package vk

type Response struct {
	Response struct {
		Count int    `json:"count"`
		Items []Post `json:"items"`
	} `json:"response"`
}

type Post struct {
	ID      int     `json:"id"`
	Text    string  `json:"text"`
	Date    int64   `json:"date"`
	Media   []Media `json:"attachments"`
	OwnerID int     `json:"owner_id"`
}

type Media struct {
	Type  string `json:"type"`
	Photo *Photo `json:"photo,omitempty"`
	Video *Video `json:"video,omitempty"`
	Audio *Audio `json:"audio,omitempty"`
}

type Photo struct {
	Sizes []Size `json:"sizes"`
}

type Size struct {
	URL string `json:"url"`
}

type Video struct {
	Description string       `json:"description"`
	Image       []VideoImage `json:"image"`
}

type VideoImage struct {
	URL string `json:"url"`
}

type Audio struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	URL    string `json:"url"`
}
