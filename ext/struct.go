package ext

type FilmHostType struct {
	Host string
	UrlPlay string
	Tm string
	Sign string
	UserLink string
	Path string
	Port int
}

type FilmPlayer struct {
	Status int `json:"status"`
	PlayName string `json:"playname"`
	PlayerUrl string `json:"playerurl"`
}
