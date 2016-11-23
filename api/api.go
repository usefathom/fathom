package api

import (
  "log"
  "time"
)

type Datapoint struct {
  Count int
  Label string
}

// log fatal errors
func checkError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func fillDatapoints(days int, points []Datapoint) []Datapoint {

  now := time.Now().AddDate(0, 0, 1)
  start := now.AddDate(0, 0, -days)
  newPoints := make([]Datapoint, days)

  for i := 0; i < days; i++ {
    newPoints[i] = Datapoint{
      Count: 0,
      Label: start.AddDate(0, 0, i).Format("2006-01-02"),
    }

    for j, p := range points {
      if p.Label == newPoints[i].Label {
        newPoints[i].Count = p.Count
        points[j] = points[len(points)-1]
        break
      }
    }
  }

  return newPoints
}
