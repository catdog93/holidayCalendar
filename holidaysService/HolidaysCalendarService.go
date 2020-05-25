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

//const CountriesCodeURL = "https://restcountries.eu/rest/v2/alpha?codes=col;no;ee"
const CountriesCodeURL = "https://api.printful.com/countries"

const ErrStringGetCountriesCodes = "error occurred in holidaysService.GetPublicHolidays(): "

type CountryCode string
type Obj []interface{}

const CustomTimeFormat = "2006-01-02"

type HolidaysCalendar struct {
	//TempCache          *Obj // caches
	PublicHolidaysResponse *[]GetPublicHolidaysResponse `json:"countriesCodes"`
	HolidaysCash           map[float64]time.Time        `json:"holidaysCash"`
	/*CountriesCodeCache *Obj
	HolidaysCache      *Obj*/
}

type CountriesCodes struct {
	Code   int                   `json:"code"`
	Result []CountriesCodeResult `json:"result"`

	/*"code": 200,
	  "result": [{
	      "code": "AF",
	      "name": "Afghanistan",
	      "states": null
	  },*/
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

type CountriesCodeResult struct {
	Code   string      `json:"code"`
	Name   string      `json:"name"`
	States interface{} `json:"states"`
}

/*
func (service *CalendarService) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	//response.Header().Set("Content-type", "application/json")

}*/

type CalendarService interface {
	GetHolidaysForThisYear(year int, code CountryCode) (err error)
	GetCountriesCodes() (err error)
}

// Method support year since 2000 to current year + 40 (2060). Country code format: "UA", "AT", "BE", "GB" (2 chars)
/*func (calendar *HolidaysCalendar) GetHolidaysForThisYear(year int, code CountryCode, serviceURL string) error {
	errorString := "error holidaysService.GetHolidaysForThisYear(): "
	if url, err := url.Parse(serviceURL); err != nil {
		errorString += err.Error()
	} else {
		if resp, err := http.Get(url.String()); err != nil { // when there is an error, the response will be nil and an error will be returned.
			errorString += err.Error()
		}
		if resp != nil { // But, when you get a redirection error, response will not be nil but there’ll be an error.
			defer resp.Body.Close()/////
			body, err := ioutil.ReadAll(resp.Body)
				calendar.TempCache = &Obj{}
				if err = json.Unmarshal(body, calendar.TempCache); err != nil {

				}
} else {


		}
	}

	if year >= 2000 && year <= time.Now().Year()+40 {
		if err = calendar.GetRequest(HolidaysURL); err != nil {
			return
		} else {
			if err = calendar.GetPublicHolidays(); err != nil {
				return
			}
		}
	} else {
		return fmt.Errorf("year isn't valid, method support year since 2000 to (current + 40)")
	}
	return nil
}*/
/*calendar.HolidaysCache == nil && */

func (calendar *HolidaysCalendar) GetPublicHolidays(serviceURL string) (*[]GetPublicHolidaysResponse, error) { //publicHolidays...GetPublicHolidaysResponse) error {
	errorString := "error occurred in holidaysService.GetPublicHolidays(): "
	publicHolidays := make([]GetPublicHolidaysResponse, AverageNumberOfPublicHolidaysPerYear)
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
	todayDate := time.Now()

	calendar.HolidaysCash = make(map[float64]time.Time)

	for _, value := range *calendar.PublicHolidaysResponse {
		if holidayDate, err := time.Parse(CustomTimeFormat, value.DateStringFormat); err != nil {
			return "", err
		} else {
			duration := math.Abs(todayDate.Sub(holidayDate).Hours())
			calendar.HolidaysCash[duration] = holidayDate
		}
	}
	fmt.Println(len(calendar.HolidaysCash))
	/*todayStringFormat = today.Format(customDateFormat)
	fmt.Println(todayStringFormat)*/
	return
}

// time zone
func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

func getMin(values ...float64) float64 {
	//min := values[0]
	return 0
}

//func getCountriesCodesFromString()
