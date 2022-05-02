package api

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/mssola/user_agent"
	"github.com/usefathom/fathom/pkg/aggregator"
	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"
)

type Collector struct {
	Store     datastore.Datastore
	Pageviews chan *models.Pageview

	// buffer vars
	updates []*models.Pageview
	inserts []*models.Pageview
	sizeu   int
	sizei   int
}

func NewCollector(store datastore.Datastore) *Collector {
	bufferCap := 100                         // persist every 100 pageviews, see https://github.com/usefathom/fathom/issues/132
	bufferTimeout := 1000 * time.Millisecond // or every 1000 ms, whichever comes first

	c := &Collector{
		Store:     store,
		Pageviews: make(chan *models.Pageview),
		updates:   make([]*models.Pageview, bufferCap),
		inserts:   make([]*models.Pageview, bufferCap),
		sizeu:     0,
		sizei:     0,
	}
	go c.aggregate()
	go c.worker(bufferCap, bufferTimeout)
	return c
}

func (c *Collector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !shouldCollect(r) {
		return
	}

	q := r.URL.Query()
	now := time.Now()

	pageview := &models.Pageview{
		ID:             uuid.NewString()[0:30],
		SiteTrackingID: q.Get("sid"),
		Hostname:       parseHostname(q.Get("h")),
		Pathname:       parsePathname(q.Get("p")),
		IsNewVisitor:   q.Get("nv") == "1",
		IsNewSession:   q.Get("ns") == "1",
		IsUnique:       q.Get("u") == "1",
		Referrer:       q.Get("r"),
		IsFinished:     false,
		IsBounce:       true,
		Duration:       0,
		Timestamp:      now,
	}

	// push pageview onto channel to be inserted (in batch) later
	c.Pageviews <- pageview

	// indicate that we're not tracking user data, see https://github.com/usefathom/fathom/issues/65
	w.Header().Set("Tk", "N")

	// headers to prevent caching
	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Expires", "Mon, 01 Jan 1990 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	// response, 1x1 px transparent GIF
	w.WriteHeader(http.StatusOK)
	b, _ := base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7")
	w.Write(b)

	// find previous pageview by same visitor
	previousPageviewID := q.Get("pid")
	if !pageview.IsNewSession && previousPageviewID != "" {
		previousPageview, err := c.Store.GetPageview(previousPageviewID)
		if err != nil && err != datastore.ErrNoResults {
			log.Errorf("error getting previous pageview: %s", err)
			return
		}

		// if we have a recent pageview that is less than 30 minutes old
		if previousPageview != nil && previousPageview.Timestamp.After(now.Add(-30*time.Minute)) {
			previousPageview.Duration = (now.Unix() - previousPageview.Timestamp.Unix())
			previousPageview.IsBounce = false
			previousPageview.IsFinished = true

			// push onto channel to be updated (in batch) later
			c.Pageviews <- previousPageview
		}
	}
}

func (c *Collector) aggregate() {
	var report aggregator.Report

	agg := aggregator.New(c.Store)
	timeout := 1 * time.Minute
	agg.Run()

	for {
		select {
		case <-time.After(timeout):
			// run aggregator at least once
			report = agg.Run()

			// if pool is not empty yet, keep running
			for !report.PoolEmpty {
				report = agg.Run()
			}
		}
	}
}

func (c *Collector) worker(cap int, timeout time.Duration) {
	var size int

	for {
		select {
		// persist pageviews in buffer when buffer at capacity
		case p := <-c.Pageviews:
			size = c.buffer(p)
			if size >= cap {
				c.persist()
			}

		// or after timeout passed
		case <-time.After(timeout):
			c.persist()
		}
	}
}

func (c *Collector) buffer(p *models.Pageview) int {
	if !p.IsFinished {
		c.inserts[c.sizei] = p
		c.sizei++
	} else {
		c.updates[c.sizeu] = p
		c.sizeu++
	}

	return (c.sizeu + c.sizei)
}

func (c *Collector) persist() {
	if (c.sizeu + c.sizei) == 0 {
		return
	}

	log.Debugf("persisting %d pageviews (%d inserts, %d updates)", (c.sizeu + c.sizei), c.sizei, c.sizeu)

	if err := c.Store.InsertPageviews(c.inserts[0:c.sizei]); err != nil {
		log.Errorf("error inserting pageviews: %s", err)
	}

	if err := c.Store.UpdatePageviews(c.updates[0:c.sizeu]); err != nil {
		log.Errorf("error updating pageviews: %s", err)
	}

	// reset buffer
	c.sizei = 0
	c.sizeu = 0
}

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

	// discard if required query vars are missing
	requiredQueryVars := []string{"h", "p"}
	q := r.URL.Query()
	for _, k := range requiredQueryVars {
		if q.Get(k) == "" {
			return false
		}
	}

	return true
}

func parsePathname(p string) string {
	return "/" + strings.TrimLeft(p, "/")
}

func parseHostname(r string) string {
	u, err := url.Parse(r)
	if err != nil {
		return ""
	}

	return u.Scheme + "://" + u.Host
}
