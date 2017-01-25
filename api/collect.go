package api

import (
	"database/sql"
	"encoding/base64"
	"log"
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
	page := models.Page{
		Path:     q.Get("p"),
		Title:    q.Get("t"),
		Hostname: q.Get("h"),
	}
	stmt, _ := datastore.DB.Prepare("SELECT p.id FROM pages p WHERE p.hostname = ? AND p.path = ? LIMIT 1")
	defer stmt.Close()
	err := stmt.QueryRow(page.Hostname, page.Path).Scan(&page.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			page.Save(datastore.DB)
		} else {
			log.Fatal(err)
		}
	}

	// find or insert visitor.
	now := time.Now()
	ip := getRequestIp(r)
	visitor := models.Visitor{
		IpAddress:        ip,
		BrowserLanguage:  q.Get("l"),
		ScreenResolution: q.Get("sr"),
		DeviceOS:         ua.OS(),
		Country:          "",
	}

	// add browser details
	visitor.BrowserName, visitor.BrowserVersion = ua.Browser()
	visitor.BrowserName = parseMajorMinor(visitor.BrowserName)

	// query by unique visitor key
	visitor.Key = visitor.GenerateKey(now.Format("2006-01-02"), visitor.IpAddress, r.UserAgent())

	stmt, _ = datastore.DB.Prepare("SELECT v.id FROM visitors v WHERE v.visitor_key = ? LIMIT 1")
	defer stmt.Close()
	err = stmt.QueryRow(visitor.Key).Scan(&visitor.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = visitor.Save(datastore.DB)
			checkError(err)
		} else {
			log.Fatal(err)
		}
	}

	pageview := models.Pageview{
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

	err = pageview.Save(datastore.DB)
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
