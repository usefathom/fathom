package api

import (
  "log"
  "time"
  "strconv"
  "net/http"
)

type Datapoint struct {
  Count int
  Label string
  Percentage float64 `json:",omitempty"`
}

const defaultPeriod = 7
const defaultLimit = 10

// log fatal errors
func checkError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func fillDatapoints(start int64, end int64, step time.Duration, points []Datapoint) []Datapoint {
  // be smart about received timestamps
  if start > end {
    tmp := end
    end = start
    start = tmp
  }

  startTime := time.Unix(start, 0)
  endTime := time.Unix(end, 0)
  newPoints := make([]Datapoint, 0)

  for startTime.Before(endTime) || startTime.Equal(endTime) {
    point := Datapoint{
      Count: 0,
      Label: startTime.Format("2006-01-02"),
    }

    for j, p := range points {
      if p.Label == point.Label || p.Label == startTime.Format("2006-01") {
        point.Count = p.Count
        points[j] = points[len(points)-1]
        break
      }
    }

    newPoints = append(newPoints, point)
    startTime = startTime.Add(step)
  }

  return newPoints
}

func getRequestedLimit(r *http.Request) int {
  limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
  if err != nil || limit == 0 {
    limit = 10
  }

  return limit
}

func getRequestedPeriods(r *http.Request) (int64, int64) {
  var before, after int64
  var err error

  before, err = strconv.ParseInt(r.URL.Query().Get("before"), 10, 64)
  if err != nil || before == 0 {
    before = time.Now().Unix()
  }

  after, err = strconv.ParseInt(r.URL.Query().Get("after"), 10, 64)
  if err != nil || before == 0 {
    after = time.Now().AddDate(0, 0, -7).Unix()
  }

  return before, after
}
