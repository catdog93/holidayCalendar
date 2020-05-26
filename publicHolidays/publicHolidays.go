package publicHolidays

import (
	"fmt"
	hol "github.com/catdog93/test-task/holidaysService"
	"time"
)

const longWeekendsString = ", and the weekend will last 3 days: "

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
	today := time.Now()
	publicHolidaysResponse, err := hol.GetPublicHolidays(today.Year())
	if err != nil {
		return err
	}
	holidays := make([]Holiday, 0, cap(publicHolidaysResponse))
	// Holiday'll store only upcoming holidays since now
	for _, value := range publicHolidaysResponse {
		holidayDate, err := time.Parse(DateFormat, value.Date)
		if err != nil {
			return err
		}
		timeDifference := holidayDate.Sub(today).Hours()
		if timeDifference >= 0 {
			holidays = append(holidays, Holiday{Name: value.Name, Date: holidayDate})
		}
	}
	calendar.Holidays = holidays
	return nil
}

// Simple method, gives info whether is holiday today
func (calendar *Calendar) IsHolidayToday() (*Holiday, error) {
	today := time.Now()
	// compare dates with time.Equal()
	for _, holiday := range calendar.Holidays {
		if today.Equal(holiday.Date) {
			return &holiday, nil
		}
	}
	return nil, nil
}

// Method returns info about the next closest holiday. Detect long weekends, if holiday is adjacent to it.
func (calendar *Calendar) GetNearHolidaysInfo() (nearHolidaysInfo string, error error) {
	holiday := calendar.Holidays[0] // the nearest holiday since today
	nearHolidaysInfo = fmt.Sprintf("The next holiday is %s, %v %v", holiday.Name, holiday.Date.Month(), holiday.Date.Day())
	weekday := holiday.Date.Weekday()
	switch weekday {
	case time.Friday, time.Saturday:
		lastWeekend := holiday.Date.Add(24 * 2 * time.Hour)
		weekendsRange := fmt.Sprintf("%v %v - %v %v.\n", holiday.Date.Month(), holiday.Date.Day(), lastWeekend.Month(), lastWeekend.Day())
		nearHolidaysInfo += longWeekendsString + weekendsRange
	case time.Sunday:
		firstWeekend := holiday.Date.Add(-24 * time.Hour)
		lastWeekend := holiday.Date.Add(24 * time.Hour)
		weekendsRange := fmt.Sprintf("%v %v - %v %v.\n", firstWeekend.Month(), firstWeekend.Day(), lastWeekend.Month(), lastWeekend.Day())
		nearHolidaysInfo += longWeekendsString + weekendsRange
	case time.Monday:
		firstWeekend := holiday.Date.Add(-24 * 2 * time.Hour)
		weekendsRange := fmt.Sprintf("%v %v - %v %v.\n", firstWeekend.Month(), firstWeekend.Day(), holiday.Date.Month(), holiday.Date.Day())
		nearHolidaysInfo += longWeekendsString + weekendsRange
	}
	return nearHolidaysInfo, nil
}
