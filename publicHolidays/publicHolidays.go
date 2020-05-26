package publicHolidays

import (
	"fmt"
	hol "github.com/catdog93/test-task/holidaysService"
	"strconv"
	"time"
)

const holidaysURLWithoutParameters = "https://date.nager.at/api/v2/publicholidays/"
const holidaysUAParameter = "/UA"

const averageNumberOfPublicHolidaysPerYear = 12
const longWeekendsString = ", and the weekend will last 3 days: "

const customTimeFormat = "2006-01-02"

type HolidaysCalendar struct {
	HolidaysDatesList []time.Time `json:"holidaysDates"`
	HolidaysNamesList []string    `json:"holidaysNames"`
}

// This method inits instance of HolidaysCalendar. HolidaysCalendar store "sorted slices" of upcoming holidays since now
func (calendar *HolidaysCalendar) initHolidaysCalendar() error {
	todayDate := time.Now()
	// Marshall URL to 3-rd Party using current year and UA country code. GET /PublicHolidays/{Year}/{CountryCode}
	publicHolidaysURL := holidaysURLWithoutParameters + strconv.Itoa(todayDate.Year()) + holidaysUAParameter
	publicHolidaysResponse, err := hol.GetPublicHolidays(publicHolidaysURL)
	if err != nil {
		return err
	}
	calendar.HolidaysDatesList = make([]time.Time, 0, averageNumberOfPublicHolidaysPerYear)
	calendar.HolidaysNamesList = make([]string, 0, averageNumberOfPublicHolidaysPerYear)
	// HolidaysCalendar'll store only upcoming holidays since now
	for _, value := range publicHolidaysResponse {
		holidayDate, err := time.Parse(customTimeFormat, value.DateStringFormat)
		if err != nil {
			return err
		}
		timeDifference := holidayDate.Sub(todayDate).Hours()
		if timeDifference >= 0 {
			calendar.HolidaysDatesList = append(calendar.HolidaysDatesList, holidayDate)
			calendar.HolidaysNamesList = append(calendar.HolidaysNamesList, value.Name)
		}
	}
	return nil
}

// Simple method, gives info whether is holiday today
func (calendar *HolidaysCalendar) IsHolidayToday() (sameTime bool, holidayName string, err error) {
	// init HolidaysCalendar instance
	err = calendar.initHolidaysCalendar()
	if err != nil {
		return
	}
	today := time.Now()
	// compare dates with time.Equal()
	for index, value := range calendar.HolidaysDatesList {
		sameTime = today.Equal(value)
		if sameTime {
			holidayName = calendar.HolidaysNamesList[index]
			break
		}
	}
	return
}

// Method returns info about next closest holiday. Detect long weekends, if holiday is adjacent to it.
func (calendar *HolidaysCalendar) GetNearHolidaysInfo() (nearHolidaysInfo string, error error) {
	holidayName := calendar.HolidaysNamesList[0] // the nearest holiday since today
	holidayDate := calendar.HolidaysDatesList[0]
	nearHolidaysInfo = fmt.Sprintf("The next holiday is %s, %v %v", holidayName, holidayDate.Month(), holidayDate.Day())
	weekday := calendar.HolidaysDatesList[0].Weekday()
	switch weekday {
	case time.Friday, time.Saturday:
		lastWeekend := holidayDate.Add(24 * 2 * time.Hour)
		weekendsRange := fmt.Sprintf("%v %v - %v %v.\n", holidayDate.Month(), holidayDate.Day(), lastWeekend.Month(), lastWeekend.Day())
		nearHolidaysInfo += longWeekendsString + weekendsRange
	case time.Sunday:
		firstWeekend := holidayDate.Add(-24 * time.Hour)
		lastWeekend := holidayDate.Add(24 * time.Hour)
		weekendsRange := fmt.Sprintf("%v %v - %v %v.\n", firstWeekend.Month(), firstWeekend.Day(), lastWeekend.Month(), lastWeekend.Day())
		nearHolidaysInfo += longWeekendsString + weekendsRange
	case time.Monday:
		firstWeekend := holidayDate.Add(-24 * 2 * time.Hour)
		weekendsRange := fmt.Sprintf("%v %v - %v %v.\n", firstWeekend.Month(), firstWeekend.Day(), holidayDate.Month(), holidayDate.Day())
		nearHolidaysInfo += longWeekendsString + weekendsRange
	}
	return nearHolidaysInfo, nil
}
