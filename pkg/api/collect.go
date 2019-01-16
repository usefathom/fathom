package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/mssola/user_agent"
	"github.com/usefathom/fathom/pkg/aggregator"
	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"
)

type Collector struct {
	Store         datastore.Datastore
	Pageviews     chan *models.Pageview
	EventTriggers chan *models.EventTrigger

	// pageview buffer vars
	pageUpdates []*models.Pageview
	pageInserts []*models.Pageview

	sizePageUpdates int
	sizePageInserts int

	// event-trigger buffer vars
	eventInserts []*models.EventTrigger
	eventUpdates []*models.EventTrigger

	sizeEventUpdates int
	sizeEventInserts int
}

func NewCollector(store datastore.Datastore) *Collector {
	bufferCap := 100                         // persist every 100 pageviews, see https://github.com/usefathom/fathom/issues/132
	bufferTimeout := 1000 * time.Millisecond // or every 1000 ms, whichever comes first

	c := &Collector{
		Store:            store,
		Pageviews:        make(chan *models.Pageview),
		EventTriggers:    make(chan *models.EventTrigger),
		pageUpdates:      make([]*models.Pageview, bufferCap),
		pageInserts:      make([]*models.Pageview, bufferCap),
		eventInserts:     make([]*models.EventTrigger, bufferCap),
		eventUpdates:     make([]*models.EventTrigger, bufferCap),
		sizePageUpdates:  0,
		sizePageInserts:  0,
		sizeEventUpdates: 0,
		sizeEventInserts: 0,
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

	if q.Get("n") != "" && q.Get("c") != "" {
		event := &models.EventTrigger{
			ID:             q.Get("id"),
			SiteTrackingID: q.Get("sid"),
			Hostname:       parseHostname(q.Get("h")),
			Pathname:       parsePathname(q.Get("p")),
			IsNewVisitor:   q.Get("nv") == "1",
			IsNewSession:   q.Get("ns") == "1",
			IsUnique:       q.Get("u") == "1",
			Referrer:       q.Get("r"),
			IsFinished:     false,
			IsBounce:       true,
			Timestamp:      now,
			EventName:      q.Get("n"),
			EventContent:   q.Get("c"),
		}

		// push event onto channel to be inserted (in batch) later
		c.EventTriggers <- event
	} else {
		pageview := &models.Pageview{
			ID:             q.Get("id"),
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

	// indicate that we're not tracking user data, see https://github.com/usefathom/fathom/issues/65
	w.Header().Set("Tk", "N")

	// headers to prevent caching
	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Expires", "Mon, 01 Jan 1990 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")

	// response, 1x1 px transparent GIF
	w.WriteHeader(http.StatusOK)
	b, _ := base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7")
	w.Write(b)
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
			size = c.pageviewBuffer(p)
			if size >= cap {
				err := c.persistPageviews(c.sizePageInserts, c.sizePageUpdates, c.pageInserts, c.pageUpdates)
				if err == nil {
					c.sizeEventInserts = 0
					c.sizeEventUpdates = 0
				}
			}

		// persist event-triggers in buffer when buffer at capacity
		case e := <-c.EventTriggers:
			size = c.eventTriggerBuffer(e)
			if size >= cap {
				err := c.persistEventTriggers(c.sizeEventInserts, c.sizeEventUpdates, c.eventInserts, c.eventUpdates)
				if err == nil {
					c.sizeEventInserts = 0
					c.sizeEventUpdates = 0
				}
			}

		// or after timeout passed
		case <-time.After(timeout):
			err := c.persistPageviews(c.sizePageInserts, c.sizePageUpdates, c.pageInserts, c.pageUpdates)
			if err == nil {
				c.sizeEventInserts = 0
				c.sizeEventUpdates = 0
			}

		case <-time.After(timeout):
			err := c.persistEventTriggers(c.sizeEventInserts, c.sizeEventUpdates, c.eventInserts, c.eventUpdates)
			if err == nil {
				c.sizeEventInserts = 0
				c.sizeEventUpdates = 0
			}
	}
}

func (c *Collector) pageviewBuffer(p *models.Pageview) int {
	if !p.IsFinished {
		c.pageInserts[c.sizePageInserts] = p
		c.sizePageInserts++
	} else {
		c.pageUpdates[c.sizePageUpdates] = p
		c.sizePageUpdates++
	}

	return (c.sizePageUpdates + c.sizePageInserts)
}

func (c *Collector) eventTriggerBuffer(e *models.EventTrigger) int {
	if !e.IsFinished {
		c.eventInserts[c.sizeEventInserts] = e
		c.sizeEventInserts++
	} else {
		c.eventUpdates[c.sizeEventUpdates] = e
		c.sizeEventUpdates++
	}

	return (c.sizeEventUpdates + c.sizeEventInserts)
}

func (c *Collector) persist(sizei, sizeu int) {
	

	log.Debugf("persisting %d pageviews (%d inserts, %d updates)", (sizeu + sizei), sizei, sizeu)

	// reset buffer
	return 0, 0, nil
}

func (c *Collector) persistPageviews(sizei, sizeu int, inserts, updates []*models.Pageview) error {
	if (sizeu + sizei) == 0 {
		return fmt.Errorf("No need to reset as `%d + %d == 0`", sizeu, sizei)
	}

	if err := c.Store.InsertPageviews(inserts[0:sizei]); err != nil {
		log.Errorf("error inserting pageviews: %s", err)
	}

	if err := c.Store.UpdatePageviews(updates[0:sizeu]); err != nil {
		log.Errorf("error updating pageviews: %s", err)
	}

	return  nil
}

func (c *Collector) persistEventTriggers(sizei, sizeu int, inserts, updates []*models.EventTrigger) (int, int, error) {
	si, su, err := c.persist(sizei, sizeu)
	if err != nil {
		return si, su, err
	}

	if err := c.Store.InsertEventTriggers(inserts[0:sizei]); err != nil {
		log.Errorf("error inserting event-trriggers: %s", err)
	}

	if err := c.Store.UpdateEventTriggers(updates[0:sizeu]); err != nil {
		log.Errorf("error updating event-triggers: %s", err)
	}

	return si, su, nil
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
	requiredQueryVars := []string{"id", "h", "p"}
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
