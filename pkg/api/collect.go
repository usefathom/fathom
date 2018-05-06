package api

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/mssola/user_agent"
	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"

	log "github.com/sirupsen/logrus"
)

var buffer []*models.RawPageview
var bufferSize = 50
var timeout = 200 * time.Millisecond

func persistPageviews() {
	if len(buffer) > 0 {
		err := datastore.SaveRawPageviews(buffer)
		if err != nil {
			log.Errorf("error saving pageviews: %s", err)
		}

		// clear buffer regardless of error... this means data loss, but better than filling the buffer for now
		buffer = buffer[:0]
	}
}

func processBuffer(pv chan *models.RawPageview) {
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
	pageviews := make(chan *models.RawPageview, bufferSize)
	go processBuffer(pageviews)

	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		// abort if this is a bot.
		userAgent := r.UserAgent()
		ua := user_agent.New(userAgent)
		if ua.Bot() {
			return nil
		}

		q := r.URL.Query()
		now := time.Now()

		// get pageview details
		pageview := &models.RawPageview{
			SessionID:    q.Get("sid"),
			Pathname:     q.Get("p"),
			IsNewVisitor: q.Get("n") == "1",
			IsUnique:     q.Get("u") == "1",
			IsBounce:     q.Get("b") != "0",
			Referrer:     q.Get("r"),
			Duration:     0,
			Timestamp:    now,
		}

		err := datastore.SaveRawPageview(pageview)
		if err != nil {
			return err
		}
		// push onto channel
		//pageviews <- pageview

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
