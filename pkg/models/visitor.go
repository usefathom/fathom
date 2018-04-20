package models

type Visitor struct {
	ID               int64
	Key              string
	BrowserName      string
	BrowserVersion   string
	BrowserLanguage  string
	Country          string
	DeviceOS         string
	IpAddress        string
	ScreenResolution string
}
