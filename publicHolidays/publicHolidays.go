package publicHolidays

import (
	"fmt"
	hol "github.com/catdog93/test-task/holidaysService"
	"time"
)

const (
	nextHolidayString  = "The next holiday is %s, %v %v"
	longWeekendsString = ", and the weekend will last 3 days: "
)

const DateFormat = "2006-01-02"

type Calendar struct {
	Holidays []Holiday `json:"holidays"`
}

type Holiday struct {
	Date time.Time `json:"date"`
	Name string    `json:"name"`
}

// This method inits instance of Holiday. Holiday store "sorted slices" of upcoming holidays since now
func (calendar *Calendar) InitHolidaysCalendar() error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	publicHolidaysResponse, err := hol.GetPublicHolidays(today.Year())
	if err != nil {
		return err
	}

	holidays := make([]Holiday, 0, cap(publicHolidaysResponse))

	// Holiday'll store only upcoming holidays since now
	for _, holiday := range publicHolidaysResponse {

		holidayDate, err := time.Parse(DateFormat, holiday.Date)
		if err != nil {
			return err
		}

		timeDifference := holidayDate.Sub(today).Hours()

		if timeDifference >= 0 {
			holidays = append(holidays, Holiday{Name: holiday.Name, Date: holidayDate})
		}
	}
	calendar.Holidays = holidays
	return nil
}

// Gives info whether are holidays today
func (calendar *Calendar) IsHolidaysToday() (todayHolidays []Holiday, found bool) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// compare dates with time.Equal()
	for _, holiday := range calendar.Holidays {
		if today.Equal(holiday.Date) {
			todayHolidays = append(todayHolidays, holiday)
			found = true
		} else {
			return
		}
	}
	return
}

// Method returns info about the next closest holiday. Detects long weekends, if holidays are adjacent to it.
func (calendar *Calendar) GetNearestHolidays() (holidays []Holiday, found bool) {
	var nearestHoliday Holiday
	for index, holiday := range calendar.Holidays {
		if index == 0 {
			found = true
			nearestHoliday = holiday
			holidays = append(holidays, holiday)
			continue
		}

		if nearestHoliday.Date.Equal(holiday.Date) {
			holidays = append(holidays, holiday)
		} else {
			return
		}
	}
	return
}

// Method returns info about the next closest holiday. Detects long weekends, if holiday is adjacent to it.
func (calendar *Calendar) GetNearHolidaysInfo() (nearHolidaysInfo string) {
	holiday := calendar.Holidays[0] // the nearest holiday since today
	nearHolidaysInfo = fmt.Sprintf(nextHolidayString, holiday.Name, holiday.Date.Month(), holiday.Date.Day())
	weekday := holiday.Date.Weekday()
	weekendsRangeInfo := ""
	switch weekday {
	case time.Friday, time.Saturday:
		weekendsRangeInfo += getWeekendsRangeInfo(holiday.Date)
	case time.Sunday:
		weekendsRangeInfo += getWeekendsRangeInfo(holiday.Date.Add(-24 * time.Hour))
	case time.Monday:
		weekendsRangeInfo += getWeekendsRangeInfo(holiday.Date.Add(-24 * 2 * time.Hour))
	}
	nearHolidaysInfo += weekendsRangeInfo
	return nearHolidaysInfo
}

// Function calculates range of long weekends, when holiday is adjacent to it.
func getWeekendsRangeInfo(firstWeekend time.Time) string {
	lastWeekend := firstWeekend.Add(24 * 2 * time.Hour)
	weekendsRange := fmt.Sprintf("%v %v - %v %v.\n", firstWeekend.Month(), firstWeekend.Day(), lastWeekend.Month(), lastWeekend.Day())
	return longWeekendsString + weekendsRange
}
