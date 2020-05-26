package holidaysService

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const averageNumberOfPublicHolidaysPerYear = 12

// 3-rd Party API's response contains json array of such objects (holidays)
type GetPublicHolidaysResponse struct {
	Name             string   `json:"name"`
	DateStringFormat string   `json:"date"`
	LocalName        string   `json:"localName"`
	CountryCode      string   `json:"countryCode"`
	Fixed            bool     `json:"fixed"`
	Global           bool     `json:"global"`
	Countries        []string `json:"countries"`
	LaunchYear       uint     `json:"launchYear"`
	Type             string   `json:"type"`
}

// GET request to 3-rd Party API to get json array of holidays
func GetPublicHolidays(serviceURL string) ([]GetPublicHolidaysResponse, error) {
	publicHolidays := make([]GetPublicHolidaysResponse, 0, averageNumberOfPublicHolidaysPerYear)
	publicHolidaysURL, err := url.Parse(serviceURL)
	if err != nil {
		return nil, err
	} else {
		resp, err := http.Get(publicHolidaysURL.String())
		if err != nil {
			return nil, err
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			err = resp.Body.Close()
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(body, &publicHolidays)
			if err != nil {
				return nil, err
			}
		}
	}
	return publicHolidays, nil
}
