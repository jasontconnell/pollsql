package main

import (
	"flag"
	"fmt"
	"github.com/jasontconnell/pollsql/conf"
	"github.com/jasontconnell/pollsql/process"
	"os"
	"strings"
)

func main() {
	query := flag.String("q", "select 1", "the query to poll")
	desired := flag.Uint64("d", 0, "desired value. default 0")
	seconds := flag.Int("s", 5, "poll time in seconds")
	flag.Parse()

	if !strings.HasPrefix(*query, "select count(*)") {
		fmt.Println("I only support select count(*) queries at the moment")
		os.Exit(1)
	}

	cfg := conf.LoadConfig("config.json")

	process.Poll(cfg.ConnectionString, *query, *desired, *seconds)
}
