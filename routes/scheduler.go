package routes

import (
	"pinmarker/schedulers"
	"time"

	"github.com/robfig/cron"
)

func SetUpScheduler() {

	// Initialize Scheduler
	houseKeepingScheduler := schedulers.NewHouseKeepingScheduler()

	// Init Scheduler
	c := cron.New()
	Scheduler(c, houseKeepingScheduler)
	c.Start()
	defer c.Stop()
}

func Scheduler(c *cron.Cron, houseKeepingScheduler *schedulers.HouseKeepingScheduler) {
	// For Production
	c.AddFunc("0 5 2 * *", houseKeepingScheduler.SchedulerMonthlyLog)

	// For Development
	go func() {
		time.Sleep(5 * time.Second)

		// House Keeping Scheduler
		houseKeepingScheduler.SchedulerMonthlyLog()
	}()
}
