package hub

import (
	"fmt"
	"github.com/robfig/cron"
)

// https://godoc.org/github.com/robfig/cron
//https://www.cnblogs.com/zuxingyu/p/6023919.html
func CronTest() {
	c := cron.New()
	c.AddFunc("0 30 * * * *", func() { fmt.Println("Every hour on the half hour") })
	c.AddFunc("@hourly", func() { fmt.Println("Every hour") })
	c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty") })
	c.Start()

	// Funcs are invoked in their own goroutine, asynchronously.

	// Funcs may also be added to a running Cron
	c.AddFunc("@daily", func() { fmt.Println("Every day") })

	// Inspect the cron job entries' next and previous run times.
	//do something

	c.Stop() // Stop the scheduler (does not stop any jobs already running).
	select {}
}
