package db

import (
	"mail-provider/config"
	"time"
)

func CronMaintenance() {
	for {
		QueryAndCopyStrategys()
		AlarmVerificationsFailed()
		time.Sleep(time.Duration(config.Config().CronStep) * time.Second)
	}
}
