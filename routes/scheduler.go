package routes

import (
	"pinmarker/schedulers"
	"pinmarker/services"
	"time"

	"github.com/robfig/cron"
)

func SetUpScheduler(trackService services.TrackService) {
	// Initialize Scheduler
	houseKeepingScheduler := schedulers.NewHouseKeepingScheduler()
	auditScheduler := schedulers.NewAuditScheduler(trackService)
	cleanScheduler := schedulers.NewCleanScheduler(trackService)

	// Init Scheduler
	c := cron.New()
	Scheduler(c, houseKeepingScheduler, auditScheduler, cleanScheduler)
	c.Start()
	defer c.Stop()
}

func Scheduler(c *cron.Cron, houseKeepingScheduler *schedulers.HouseKeepingScheduler, auditScheduler *schedulers.AuditScheduler, cleanScheduler *schedulers.CleanScheduler) {
	// For Production
	c.AddFunc("0 5 2 * *", houseKeepingScheduler.SchedulerMonthlyLog)
	c.AddFunc("0 0 2 * * *", auditScheduler.SchedulerAuditAppsUserTotal)
	c.AddFunc("0 0 3 * * *", auditScheduler.SchedulerAuditAppsUserTotal)

	// For Development
	go func() {
		time.Sleep(5 * time.Second)

		// Audit Scheduler
		auditScheduler.SchedulerAuditAppsUserTotal()

		// Clean Scheduler
		cleanScheduler.SchedulerCleanAllTracksCreatedByDays()

		// House Keeping Scheduler
		houseKeepingScheduler.SchedulerMonthlyLog()
	}()
}
