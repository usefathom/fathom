package count

import(
  "testing"
  "time"
)

func TestFill(t *testing.T) {
  start := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.Local)
  end := time.Date(2016, time.January, 5, 0, 0, 0, 0, time.Local)

  points := []Point{
    Point{
      Label: start.Format("2006-01-01"),
      Value: 1,
    },
    Point {
      Label: end.Format("2006-01-01"),
      Value: 1,
    },
  }

  filled := fill(start.Unix(), end.Unix(), points)
  if len(filled) != 5 {
    t.Error("Length of filled points slice does not match expected length")
  }

}
