package api

import "net/http"

// GET /health
func (api *API) Health(w http.ResponseWriter, _ *http.Request) error {
	if err := api.database.Health(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
