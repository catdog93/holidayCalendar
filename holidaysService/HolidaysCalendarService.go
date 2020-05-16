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

type CalendarService interface {
	GetHolidaysForThisYear(year int, code CountryCode) (err error)
	GetCountriesCodes() (err error)
}

/*
func (service *CalendarService) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	//response.Header().Set("Content-type", "application/json")

}*/

type HolidaysCalendar struct {
	TempCache          *Obj // caches
	CountriesCodeCache *Obj
	HolidaysCache      *Obj
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
	//calendar.TempCache = &Obj{}
	if err := json.Unmarshal(body, calendar.TempCache); err != nil {
		return err
	} else {
		return nil
	}
}

// Method support year since 2000 to current year + 40 (2060). Country code format: "UA", "AT", "BE", "GB" (2 chars)
func (calendar *HolidaysCalendar) GetHolidaysForThisYear(year int, code CountryCode) (err error) {
	if /*calendar.HolidaysCache == nil && */ year >= 2000 && year <= time.Now().Year()+40 {
		err = calendar.getRequest(HolidaysURL)
		return
	} else {
		return fmt.Errorf("year isn't valid, method support year since 2000 to (current + 40)")
	}
}

func (calendar *HolidaysCalendar) GetCountriesCodes() (err error) {
	//if calendar.HolidaysCache == nil &&
	err = calendar.getRequest(CountriesCodeURL2)
	return
}
