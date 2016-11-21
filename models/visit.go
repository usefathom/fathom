package models

type Visit struct {
  ID int
  Path string
  IpAddress string
  ReferrerKeyword string
  ReferrerType string
  ReferrerUrl string
  DeviceBrand string
  DeviceModel string
  DeviceType string
  DeviceOS string
  BrowserName string
  BrowserVersion string
  BrowserLanguage string
  ScreenResolution string
  VisitorReturning bool
  Country string
}
