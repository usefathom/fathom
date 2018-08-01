package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"time"
)

var statsCmd = cli.Command{
	Name:   "stats",
	Usage:  "view stats",
	Action: stats,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "start-date",
			Usage: "start date, expects a date in format 2006-01-02",
		},
		cli.StringFlag{
			Name:  "end-date",
			Usage: "end date, expects a date in format 2006-01-02",
		},
		cli.BoolFlag{
			Name:  "json",
			Usage: "get a json response",
		},
	},
}

func stats(c *cli.Context) error {
	start, _ := time.Parse("2006-01-02", c.String("start-date"))
	if start.IsZero() {
		return errors.New("Invalid argument: supply a valid --start-date")
	}

	end, _ := time.Parse("2006-01-02", c.String("end-date"))
	if end.IsZero() {
		return errors.New("Invalid argument: supply a valid --end-date")
	}

	result, err := app.database.GetAggregatedSiteStats(start, end)
	if err != nil {
		return err
	}

	if c.Bool("json") {
		if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
			return err
		}

		return nil
	}

	fmt.Printf("%s - %s\n", start.Format("Jan 01, 2006"), end.Format("Jan 01, 2006"))
	fmt.Printf("===========================\n")
	fmt.Printf("Visitors: \t%d\n", result.Visitors)
	fmt.Printf("Pageviews: \t%d\n", result.Pageviews)
	fmt.Printf("Sessions: \t%d\n", result.Sessions)
	fmt.Printf("Avg duration: \t%s\n", result.FormattedDuration())
	fmt.Printf("Bounce rate: \t%.0f%%\n", result.BounceRate*100.00)
	return nil
}
