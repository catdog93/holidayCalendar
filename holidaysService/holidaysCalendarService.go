package holidaysService

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	publicHolidaysURL = "https://date.nager.at/api/v2/publicholidays/%s/%s"
	CountyCode        = "UA"

	defaultTimeout = 15
)

// 3-rd Party API's response contains json array of such objects (holidays)
type publicHolidayResponse struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

func GetPublicHolidaysMock() ([]publicHolidayResponse, error) {
	jsonData := []byte("[{\"date\":\"2021-03-28\",\"localName\":\"День праці\",\"name\":\"International Workers' Day\",\"countryCode\":\"UA\",\"fixed\":true,\"global\":true,\"counties\":null,\"launchYear\":null,\"type\":\"Public\"}," +
		"{\"date\":\"2021-03-28\",\"localName\":\"Великдень\",\"name\":\"(Julian) Easter Sunday\",\"countryCode\":\"UA\",\"fixed\":false,\"global\":true,\"counties\":null,\"launchYear\":null,\"type\":\"Public\"}]")

	var publicHolidays []publicHolidayResponse
	err := json.Unmarshal(jsonData, &publicHolidays)
	return publicHolidays, err
}

// GET request to 3-rd Party API to get json array of holidays
func GetPublicHolidays(year int) ([]publicHolidayResponse, error) {
	// Marshall URL to 3-rd Party using current year and UA country code. GET /PublicHolidays/{Year}/{CountryCode}
	publicHolidaysURL := fmt.Sprintf(publicHolidaysURL, strconv.Itoa(year), CountyCode)

	var publicHolidays []publicHolidayResponse

	client := http.Client{Timeout: defaultTimeout * time.Second}

	resp, err := client.Get(publicHolidaysURL)
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
