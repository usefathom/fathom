package api

import (
	"encoding/base64"
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
	// abort if this is a bot.
	ua := user_agent.New(r.UserAgent())
	if ua.Bot() {
		return false
	}

	// abort if DNT header is set to "1" (these should have been filtered client-side already)
	if r.Header.Get("DNT") == "1" {
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
	go aggregate(api.database)

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
				err := api.database.UpdatePageview(previousPageview)
				if err != nil {
					return err
				}
			}
		}

		// save new pageview
		err := api.database.SavePageview(pageview)
		if err != nil {
			return err
		}

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
