package holidaysService

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const HolidaysURL = "https://date.nager.at/api/v2/publicholidays/2020/UA"
const CountriesCodeURL = "https://restcountries.eu/rest/v2/alpha?codes=col;no;ee"
const CountriesCodeURL2 = "https://api.printful.com/countries"

type CountryCode string
type Obj []interface{}

type HolidaysCalendar struct {
	TempCache          *Obj // caches
	StructType         *structType
	CountriesCodeCache *Obj
	HolidaysCache      *Obj
}

type structType struct {
	Code   int               `json:"code"`
	Result []structInnerType `json:"result"`

	/*"code": 200,
	  "result": [{
	      "code": "AF",
	      "name": "Afghanistan",
	      "states": null
	  },*/
}

type structInnerType struct {
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

func (calendar *HolidaysCalendar) getRequest(serviceURL string) (err error) {
	url, err := url.Parse(serviceURL)
	if err != nil {
		return nil
	}
	resp, err := http.Get(url.String())
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	calendar.TempCache = &Obj{}
	if err = json.Unmarshal(body, calendar.TempCache); err != nil {
		calendar.StructType = &structType{}
		if err = json.Unmarshal(body, calendar.StructType); err != nil {
			return
		} else {
			return nil
		}
	} else {
		return nil
	}
}

// Method support year since 2000 to current year + 40 (2060). Country code format: "UA", "AT", "BE", "GB" (2 chars)
func (calendar *HolidaysCalendar) GetHolidaysForThisYear(year int, code CountryCode) (err error) {
	if /*calendar.HolidaysCache == nil && */ year >= 2000 && year <= time.Now().Year()+40 {
		if err = calendar.getRequest(HolidaysURL); err != nil {
			return
		} else {
			if err = calendar.GetCountriesCodes(); err != nil {
				return
			}
		}
	} else {
		return fmt.Errorf("year isn't valid, method support year since 2000 to (current + 40)")
	}
}

func (calendar *HolidaysCalendar) GetCountriesCodes() (err error) {
	//if calendar.HolidaysCache == nil &&
	calendar.getRequest(CountriesCodeURL2)

	return
}

//func getCountriesCodesFromString()
