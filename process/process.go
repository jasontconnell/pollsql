package process

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/pkg/errors"
	"math"
	"time"
)

type PollResult struct {
	Count   uint64
	Time    time.Time
	Desired uint64
}

type PollStat struct {
	Diff uint64
	Time time.Duration
}

func Poll(connstr, query string, desired uint64, seconds int) {
	results := make(chan PollResult, 3)
	tick := time.Tick(time.Second * time.Duration(seconds))
	go calc(results, desired)
	for now := range tick {
		val, err := pollOne(connstr, query)
		if err != nil {
			fmt.Println("error", err)
		}
		pr := PollResult{Count: val, Time: now}
		results <- pr
	}
}

func calc(results chan PollResult, desired uint64) {
	var lastResult PollResult

	for {
		select {
		case pr := <-results:
			if lastResult.Count == 0 { // first result
				lastResult = pr
			}
			ps := PollStat{Diff: uint64(math.Abs(float64(lastResult.Count) - float64(pr.Count))), Time: pr.Time.Sub(lastResult.Time)}
			lastResult = pr

			var rps float64
			if ps.Time.Seconds() != 0 {
				rps = float64(ps.Diff) / ps.Time.Seconds()
			}

			var ttz float64
			if ps.Diff != 0 {
				chunksLeft := float64(pr.Count) / float64(ps.Diff)
				ttz = chunksLeft * rps
			}

			unit := "seconds"
			if ttz > 10000 {
				unit = "hours"
				ttz = ttz / float64(3600)
			} else if ttz > 100 {
				unit = "minutes"
				ttz = ttz / float64(60.0)
			}

			fmt.Printf("\r %v %v left\t\tRows per second: %v\tLeft: %v", uint64(ttz), unit, uint64(rps), pr.Count)
		}
	}
}

func pollOne(connstr, query string) (uint64, error) {
	db, err := sql.Open("mssql", connstr)
	if err != nil {
		return 0, errors.Wrapf(err, "couldn't open connection to sql server with %s", connstr)
	}
	defer db.Close()

	row := db.QueryRow(query)

	var val uint64
	serr := row.Scan(&val)
	if serr != nil {
		return 0, errors.Wrapf(serr, "scan error with %s", query)
	}

	return val, nil
}
