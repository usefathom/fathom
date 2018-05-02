package models

type Page struct {
	ID       int64  `json:"-"`
	Scheme   string `json:"scheme"`
	Hostname string `json:"hostname"`
	Path     string `json:"path"`
	Title    string `json:"title"`
}
