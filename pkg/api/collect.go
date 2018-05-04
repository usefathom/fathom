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
var timeout = 200 * time.Millisecond

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
	pageviews := make(chan *models.Pageview, bufferSize)
	go processBuffer(pageviews)

	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {

		// abort if this is a bot.
		userAgent := r.UserAgent()
		ua := user_agent.New(userAgent)
		if ua.Bot() {
			return nil
		}

		q := r.URL.Query()

		// find page
		page, err := datastore.GetPageByHostnameAndPath(q.Get("h"), q.Get("p"))
		if err != nil && err != datastore.ErrNoResults {
			return err
		}

		// page does not exist yet, get details & save it
		if page == nil {
			page = &models.Page{
				Scheme:   "http",
				Hostname: q.Get("h"),
				Path:     q.Get("p"),
				Title:    q.Get("t"),
			}

			if scheme := q.Get("scheme"); scheme != "" {
				page.Scheme = scheme
			}

			err = datastore.SavePage(page)
			if err != nil {
				return err
			}
		}

		// find visitor by anonymized key from query params
		now := time.Now()
		visitorKey := q.Get("vk")
		visitorKey = enhanceVisitorKey(visitorKey, now.Format("2006-01-02"), userAgent, q.Get("l"), q.Get("sr"))
		visitor, err := datastore.GetVisitorByKey(visitorKey)
		if err != nil && err != datastore.ErrNoResults {
			return err
		}

		// visitor is new, save it
		if visitor == nil {
			visitor = &models.Visitor{
				BrowserLanguage:  q.Get("l"),
				ScreenResolution: q.Get("sr"),
				DeviceOS:         ua.OS(),
				Country:          "",
				Key:              visitorKey,
			}

			// add browser details
			visitor.BrowserName, visitor.BrowserVersion = ua.Browser()

			// get rid of exact browser versions
			visitor.BrowserVersion = parseMajorMinor(visitor.BrowserVersion)
			err = datastore.SaveVisitor(visitor)
			if err != nil {
				return err
			}
		} else {
			lastPageview, err := datastore.GetLastPageviewForVisitor(visitor.ID)
			if err != nil && err != datastore.ErrNoResults {
				return err
			}

			if lastPageview != nil && lastPageview.Timestamp.After(now.Add(-30*time.Minute)) {
				lastPageview.Bounced = false
				lastPageview.TimeOnPage = now.Unix() - lastPageview.Timestamp.Unix()

				// TODO: Delay storage until in buffer?
				err := datastore.UpdatePageview(lastPageview)
				if err != nil {
					return err
				}
			}
		}

		// get pageview details
		pageview := &models.Pageview{
			PageID:          page.ID,
			VisitorID:       visitor.ID,
			ReferrerUrl:     q.Get("ru"),
			ReferrerKeyword: q.Get("rk"),
			TimeOnPage:      0,
			Bounced:         true, // TODO: Only mark as bounced if no other pageviews in this session
			Timestamp:       now,
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
func enhanceVisitorKey(key string, date string, userAgent string, lang string, screenRes string) string {
	byteKey := md5.Sum([]byte(date + userAgent + lang + screenRes))
	return hex.EncodeToString(byteKey[:])
}
