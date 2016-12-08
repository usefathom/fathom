package models

type Visit struct {
  ID int64
  PageID int64
  IpAddress string
  ReferrerKeyword string
  ReferrerUrl string
  BrowserName string
  BrowserVersion string
  BrowserLanguage string
  ScreenResolution string
  Country string
  Timestamp string
}
