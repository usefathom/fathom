package api

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/mssola/user_agent"
	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"

	log "github.com/sirupsen/logrus"
)

var buffer []*models.Pageview
var bufferSize = 250
var timeout = 100 * time.Millisecond

func persistPageviews() {
	if len(buffer) > 0 {
		err := datastore.SavePageviews(buffer)
		if err != nil {
			log.Errorf("error saving pageviews: %s", err)
		}

		// clear buffer regardless of error... this means data loss, but better than filling the buffer for now
		buffer = buffer[:0]
	}
}

func processBuffer(pv chan *models.Pageview) {
	for {
		select {
		case pageview := <-pv:
			buffer = append(buffer, pageview)
			if len(buffer) >= bufferSize {
				persistPageviews()
			}
		case <-time.After(timeout):
			persistPageviews()
		}
	}
}

/* middleware */
func NewCollectHandler() http.Handler {
	pageviews := make(chan *models.Pageview, 100)
	go processBuffer(pageviews)

	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {

		// abort if this is a bot.
		ua := user_agent.New(r.UserAgent())
		if ua.Bot() {
			return nil
		}

		q := r.URL.Query()

		// find or insert page
		page, err := datastore.GetPageByHostnameAndPath(q.Get("h"), q.Get("p"))
		if err != nil {
			page = &models.Page{
				Hostname: q.Get("h"),
				Path:     q.Get("p"),
				Title:    q.Get("t"),
			}

			err = datastore.SavePage(page)
			if err != nil {
				return err
			}
		}

		// find or insert visitor.
		now := time.Now()
		ipAddress := getRequestIp(r)
		visitorKey := generateVisitorKey(now.Format("2006-01-02"), ipAddress, r.UserAgent())

		visitor, err := datastore.GetVisitorByKey(visitorKey)
		if err != nil {
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
			if err != nil {
				return err
			}
		}

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

		// push onto channel
		pageviews <- pageview

		// don't you cache this
		w.Header().Set("Content-Type", "image/gif")
		w.Header().Set("Expires", "Mon, 01 Jan 1990 00:00:00 GMT")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.WriteHeader(http.StatusOK)

		// 1x1 px transparent GIF
		b, _ := base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7")
		w.Write(b)
		return nil
	})
}

// generateVisitorKey generates the "unique" visitor key from date, user agent + screen resolution
func generateVisitorKey(date string, ipAddress string, userAgent string) string {
	byteKey := md5.Sum([]byte(date + ipAddress + userAgent))
	return hex.EncodeToString(byteKey[:])
}
