package holidaysService

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"time"
)

const AverageNumberOfPublicHolidaysPerYear = 12

const PublicHolidaysURL = "https://date.nager.at/api/v2/publicholidays/2020/UA"

const ErrStringGetCountriesCodes = "error occurred in holidaysService.getPublicHolidays(): "

type CountryCode string
type Obj []interface{}

const CustomTimeFormat = "2006-01-02"

type HolidaysCalendar struct {
	PublicHolidaysResponse         *[]GetPublicHolidaysResponse `json:"countriesCodes"`
	HolidaysDifferenceSinceNowList *[]float64                   `json:"holidays"`
	*PublicHolidays
}

/*"date": "2020-01-07",
"localName": "Різдво",
"name": "(Julian) Christmas",
"countryCode": "UA",
"fixed": true,
"global": true,
"counties": null,
"launchYear": null,
"type": "Public"*/

type GetPublicHolidaysResponse struct {
	Name             string      `json:"name"`
	DateStringFormat string      `json:"date"`
	LocalName        string      `json:"localName"`
	CountryCode      CountryCode `json:"countryCode"`
	Fixed            bool        `json:"fixed"`
	Global           bool        `json:"global"`
	Countries        []string    `json:"countries"`
	LaunchYear       uint        `json:"launchYear"`
	Type             string      `json:"type"`
}

type PublicHolidays struct {
	HolidaysDatesList []time.Time `json:"holidaysDates"`
	HolidaysNamesList []string    `json:"holidaysNames"`
}

func (calendar *HolidaysCalendar) initCaches() error {
	if PublicHolidaysResponse, err := calendar.getPublicHolidays(PublicHolidaysURL); err != nil {
		return err
	} else {
		calendar.PublicHolidaysResponse = PublicHolidaysResponse
		todayDate := time.Now()
		todayDate = time.Date(todayDate.Year(), 6, 20, 0, 0, 0, 0, time.UTC)
		calendar.HolidaysDifferenceSinceNowList = &[]float64{}
		calendar.PublicHolidays = &PublicHolidays{}

		*calendar.HolidaysDifferenceSinceNowList = make([]float64, 0, AverageNumberOfPublicHolidaysPerYear)
		calendar.HolidaysDatesList = make([]time.Time, 0, AverageNumberOfPublicHolidaysPerYear)
		calendar.HolidaysNamesList = make([]string, 0, AverageNumberOfPublicHolidaysPerYear)

		for _, value := range *calendar.PublicHolidaysResponse {
			if holidayDate, err := time.Parse(CustomTimeFormat, value.DateStringFormat); err != nil {
				return err
			} else {
				if duration := todayDate.Sub(holidayDate).Hours(); duration <= 0 {
					*calendar.HolidaysDifferenceSinceNowList = append(*calendar.HolidaysDifferenceSinceNowList, math.Abs(duration))
					calendar.HolidaysDatesList = append(calendar.HolidaysDatesList, holidayDate)
					calendar.HolidaysNamesList = append(calendar.HolidaysNamesList, value.Name)
				}
			}
		}
	}
	return nil
}

func (calendar *HolidaysCalendar) IsHolidayToday() (bool, string, error) {
	if err := calendar.initCaches(); err != nil {
		return false, "", err
	}

	today := time.Now()
	today = time.Date(today.Year(), 6, 20, 0, 0, 0, 0, time.UTC)
	//today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	sameTime := false
	holidayName := ""
	// compare time with time.Equal()
	for index, _ := range calendar.HolidaysDatesList {
		sameTime = today.Equal(calendar.HolidaysDatesList[index])
		if sameTime {
			holidayName = calendar.HolidaysNamesList[index]
			break
		}
	}
	return sameTime, holidayName, nil
}

func (calendar *HolidaysCalendar) getPublicHolidays(serviceURL string) (*[]GetPublicHolidaysResponse, error) { //publicHolidays...GetPublicHolidaysResponse) error {
	errorString := "error occurred in holidaysService.getPublicHolidays(): "
	publicHolidays := make([]GetPublicHolidaysResponse, 0, AverageNumberOfPublicHolidaysPerYear)
	if url, err := url.Parse(serviceURL); err != nil {
		errorString += err.Error()
	} else {
		resp, err := http.Get(url.String())
		if err != nil { // when there is an error, the response will be nil and an error will be returned.
			errorString += err.Error()
		} else {
			if resp != nil { // But, when you get a redirection error, response will not be nil but there’ll be an error.
				body, err := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()

				if err = json.Unmarshal(body, &publicHolidays); err != nil { //array
					errorString += err.Error()
				}
			} else {
				errorString += "third party service returned nil response"
			}
		}
	}
	if errorString != ErrStringGetCountriesCodes { // expression looks like if errorString != ""
		return nil, fmt.Errorf(errorString)
	}
	return &publicHolidays, nil
}

func (calendar *HolidaysCalendar) GetNearHolidaysInfo() (nearHolidaysInfo string, error error) {
	nearHolidaysInfo = fmt.Sprintf("The next holiday is %s, %v %v", calendar.HolidaysNamesList[0], calendar.HolidaysDatesList[0].Month(), calendar.HolidaysDatesList[0].Day())
	holidayDate := calendar.HolidaysDatesList[0]
	weekday := calendar.HolidaysDatesList[0].Weekday()
	longWeekendsString := ", and the weekend will last 3 days: "
	switch weekday {
	case 1:
		firstWeekend := holidayDate.Add(-24 * 2 * time.Hour)
		temp := fmt.Sprintf("%v %v - %v %v.\n", firstWeekend.Month(), firstWeekend.Day(), holidayDate.Month(), holidayDate.Day())
		nearHolidaysInfo += longWeekendsString + temp
	case 5, 6:
		lastWeekend := holidayDate.Add(24 * 2 * time.Hour)
		temp := fmt.Sprintf("%v %v - %v %v.\n", holidayDate.Month(), holidayDate.Day(), lastWeekend.Month(), lastWeekend.Day())
		nearHolidaysInfo += longWeekendsString + temp
	case 0:
		firstWeekend := holidayDate.Add(-24 * time.Hour)
		lastWeekend := holidayDate.Add(24 * time.Hour)
		temp := fmt.Sprintf("%v %v - %v %v.\n", firstWeekend.Month(), firstWeekend.Day(), lastWeekend.Month(), lastWeekend.Day())
		nearHolidaysInfo += longWeekendsString + temp
	}
	return nearHolidaysInfo, nil
}

// Method support year since 2000 to current year + 40 (2060). Country code format: "UA", "AT", "BE", "GB" (2 chars)
/*
	if year >= 2000 && year <= time.Now().Year()+40 {
		if err = calendar.GetRequest(HolidaysURL); err != nil {
			return
		} else {
			if err = calendar.getPublicHolidays(); err != nil {
				return
			}
		}
	} else {
		return fmt.Errorf("year isn't valid, method support year since 2000 to (current + 40)")
	}
	return nil
}*/
/*calendar.HolidaysCache == nil && */
