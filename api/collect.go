package api

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/dannyvankooten/ana/datastore"
	"github.com/dannyvankooten/ana/models"
	"github.com/mssola/user_agent"
)

func getRequestIp(r *http.Request) string {
	ipAddress := r.RemoteAddr

	headerForwardedFor := r.Header.Get("X-Forwarded-For")
	if headerForwardedFor != "" {
		ipAddress = headerForwardedFor
	}

	return ipAddress
}

func CollectHandler(w http.ResponseWriter, r *http.Request) {
	ua := user_agent.New(r.UserAgent())

	// abort if this is a bot.
	if ua.Bot() {
		return
	}

	q := r.URL.Query()

	// find or insert page
	page, err := datastore.GetPageByHostnameAndPath(q.Get("h"), q.Get("p"))
	if page.ID == 0 {
		page = &models.Page{
			Hostname: q.Get("h"),
			Path:     q.Get("p"),
			Title:    q.Get("t"),
		}
		err = datastore.SavePage(page)
	}
	checkError(err)

	// find or insert visitor.
	now := time.Now()
	ipAddress := getRequestIp(r)
	visitorKey := generateVisitorKey(now.Format("2006-01-02"), ipAddress, r.UserAgent())

	visitor, err := datastore.GetVisitorByKey(visitorKey)
	if visitor.ID == 0 {
		visitor = &models.Visitor{
			IpAddress:        ipAddress,
			BrowserLanguage:  q.Get("l"),
			ScreenResolution: q.Get("sr"),
			DeviceOS:         ua.OS(),
			Country:          "",
			Key:              visitorKey,
		}

		// add browser details
		visitor.BrowserName, visitor.BrowserVersion = ua.Browser()
		visitor.BrowserName = parseMajorMinor(visitor.BrowserName)
		err = datastore.SaveVisitor(visitor)
	}
	checkError(err)

	pageview := &models.Pageview{
		PageID:          page.ID,
		VisitorID:       visitor.ID,
		ReferrerUrl:     q.Get("ru"),
		ReferrerKeyword: q.Get("rk"),
		Timestamp:       now.Format("2006-01-02 15:04:05"),
	}

	// only store referrer URL if not coming from own site
	if strings.Contains(pageview.ReferrerUrl, page.Hostname) {
		pageview.ReferrerUrl = ""
	}

	err = datastore.SavePageview(pageview)
	checkError(err)

	// don't you cache this
	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Expires", "Mon, 01 Jan 1990 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(http.StatusOK)

	// 1x1 px transparent GIF
	b, _ := base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7")
	w.Write(b)
}

// generateVisitorKey generates the "unique" visitor key from date, user agent + screen resolution
func generateVisitorKey(date string, ipAddress string, userAgent string) string {
	byteKey := md5.Sum([]byte(date + ipAddress + userAgent))
	return hex.EncodeToString(byteKey[:])
}
