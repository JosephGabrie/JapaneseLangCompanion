package utils

import (
    "time"
	"fmt"
    "japlearning/models"
)

func CalculateUserMastery(progress *models.Progress, userTypedAnswer bool ) int {
        
	masterylevel := progress.MasteryLevel
	fmt.Println(masterylevel)
	if userTypedAnswer {
		masterylevel++
	} else if !userTypedAnswer{
		masterylevel--
	}

	progress.MasteryLevel = masterylevel
	fmt.Println("mastery level in calculate user mastery ", masterylevel)
	 if masterylevel > 0 && masterylevel < 4{
		return 1
	} else if masterylevel >= 4 && masterylevel < 6{
		return 2
	} else if masterylevel >= 6 && masterylevel < 7{
		return 3
	} else if masterylevel >= 7{
		return 4
	}
	return 0
}

func SetNextTime(userAnswer int) time.Time {
	currentTime := time.Now()
	var nextTime time.Time

	switch userAnswer {
	case 1:
		nextTime = currentTime.Add(24 * time.Hour) // 24 hours later
	case 2:
		nextTime = currentTime.Add(72 * time.Hour) // 72 hours (3 days) later
	case 3:
		nextTime = currentTime.AddDate(0, 0, 7) // 1 week later
	case 4:
		nextTime = currentTime.AddDate(0, 1, 14) // 1 month and 2 weeks later
	default:
		nextTime = currentTime // If userAnswer doesn't match any case, keep it as current time
	}

	return nextTime
}
