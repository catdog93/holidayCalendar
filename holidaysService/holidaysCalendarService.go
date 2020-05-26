package holidaysService

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	holidaysURL          = "https://date.nager.at/api/v2/publicholidays/%s/%s"
	UACountyCode         = "UA"
	publicHolidaysNumber = 12
)

// 3-rd Party API's response contains json array of such objects (holidays)
type GetPublicHolidaysResponse struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

// GET request to 3-rd Party API to get json array of holidays
func GetPublicHolidays(year int) ([]GetPublicHolidaysResponse, error) {
	// Marshall URL to 3-rd Party using current year and UA country code. GET /PublicHolidays/{Year}/{CountryCode}
	publicHolidaysURL := fmt.Sprintf(holidaysURL, strconv.Itoa(year), UACountyCode)
	publicHolidays := make([]GetPublicHolidaysResponse, 0, publicHolidaysNumber)
	resp, err := http.Get(publicHolidaysURL)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &publicHolidays)
	return publicHolidays, err
}
