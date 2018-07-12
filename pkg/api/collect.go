package api

import (
	"encoding/base64"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mssola/user_agent"
	"github.com/usefathom/fathom/pkg/aggregator"
	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"
)

func shouldCollect(r *http.Request) bool {
	// abort if DNT header is set to "1" (these should have been filtered client-side already)
	if r.Header.Get("DNT") == "1" {
		return false
	}

	// don't track prerendered pages, see https://github.com/usefathom/fathom/issues/13
	if r.Header.Get("X-Moz") == "prefetch" || r.Header.Get("X-Purpose") == "preview" {
		return false
	}

	// abort if this is a bot.
	ua := user_agent.New(r.UserAgent())
	if ua.Bot() {
		return false
	}

	return true
}

func parsePathname(p string) string {
	return "/" + strings.TrimLeft(p, "/")
}

// TODO: Move this to aggregator, as we need this endpoint to be as fast as possible
func parseReferrer(r string) string {
	u, err := url.Parse(r)
	if err != nil {
		return ""
	}

	// remove AMP & UTM vars
	q := u.Query()
	keys := []string{"amp", "utm_campaign", "utm_medium", "utm_source"}
	for _, k := range keys {
		q.Del(k)
	}
	u.RawQuery = q.Encode()

	// remove /amp/
	if strings.HasSuffix(u.Path, "/amp/") {
		u.Path = u.Path[0:(len(u.Path) - 5)]
	}

	return u.String()
}

func parseHostname(r string) string {
	u, err := url.Parse(r)
	if err != nil {
		return ""
	}

	return u.Scheme + "://" + u.Host
}

func (api *API) NewCollectHandler() http.Handler {
	pageviews := make(chan *models.Pageview, 10)
	go aggregate(api.database)
	go collect(api.database, pageviews)

	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		if !shouldCollect(r) {
			return nil
		}

		q := r.URL.Query()
		now := time.Now()

		// get pageview details
		pageview := &models.Pageview{
			ID:           q.Get("id"),
			Hostname:     parseHostname(q.Get("h")),
			Pathname:     parsePathname(q.Get("p")),
			IsNewVisitor: q.Get("nv") == "1",
			IsNewSession: q.Get("ns") == "1",
			IsUnique:     q.Get("u") == "1",
			IsBounce:     q.Get("b") != "0",
			Referrer:     parseReferrer(q.Get("r")),
			Duration:     0,
			Timestamp:    now,
		}

		// find previous pageview by same visitor
		previousPageviewID := q.Get("pid")
		if !pageview.IsNewSession && previousPageviewID != "" {
			previousPageview, err := api.database.GetPageview(previousPageviewID)
			if err != nil && err != datastore.ErrNoResults {
				return err
			}

			// if we have a recent pageview that is less than 30 minutes old
			if previousPageview != nil && previousPageview.Timestamp.After(now.Add(-30*time.Minute)) {
				previousPageview.Duration = (now.Unix() - previousPageview.Timestamp.Unix())
				previousPageview.IsBounce = false

				// push onto channel to be updated (in batch) later
				pageviews <- previousPageview
			}
		}

		// push pageview onto channel to be inserted (in batch) later
		pageviews <- pageview

		// indicate that we're not tracking user data, see https://github.com/usefathom/fathom/issues/65
		w.Header().Set("Tk", "N")

		// headers to prevent caching
		w.Header().Set("Content-Type", "image/gif")
		w.Header().Set("Expires", "Mon, 01 Jan 1990 00:00:00 GMT")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")

		// response
		w.WriteHeader(http.StatusOK)

		// 1x1 px transparent GIF
		b, _ := base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7")
		w.Write(b)
		return nil
	})
}

// runs the aggregate func every minute
func aggregate(db datastore.Datastore) {
	agg := aggregator.New(db)
	agg.Run()

	timeout := 1 * time.Minute

	for {
		select {
		case <-time.After(timeout):
			agg.Run()
		}
	}
}

func collect(db datastore.Datastore, pageviews chan *models.Pageview) {
	var buffer []*models.Pageview
	var size = 250
	var timeout = 500 * time.Millisecond

	for {
		select {
		case pageview := <-pageviews:
			buffer = append(buffer, pageview)
			if len(buffer) >= size {
				persist(db, buffer)
				buffer = buffer[:0]
			}
		case <-time.After(timeout):
			if len(buffer) > 0 {
				persist(db, buffer)
				buffer = buffer[:0]
			}
		}
	}
}

func persist(db datastore.Datastore, pageviews []*models.Pageview) {
	n := len(pageviews)
	updates := make([]*models.Pageview, 0, n)
	inserts := make([]*models.Pageview, 0, n)

	for _, p := range pageviews {
		if !p.IsBounce {
			updates = append(updates, p)
		} else {
			inserts = append(inserts, p)
		}
	}

	log.Debugf("persisting %d pageviews (%d inserts, %d updates)", len(pageviews), len(inserts), len(updates))

	var err error
	err = db.InsertPageviews(inserts)
	if err != nil {
		log.Error(err)
	}

	err = db.UpdatePageviews(updates)
	if err != nil {
		log.Error(err)
	}
}
