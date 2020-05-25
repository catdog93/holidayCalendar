package main

import (
	hol "github.com/catdog93/test-task/holidaysService"
	"log"
)

func main() {
	// input
	calendar := &hol.HolidaysCalendar{}
	//calendar.GetHolidaysForThisYear(2019, "UA")*/
	if PublicHolidaysResponse, err := calendar.GetPublicHolidays(hol.PublicHolidaysURL); err != nil {
		log.Fatal(err)
	} else {
		calendar.PublicHolidaysResponse = PublicHolidaysResponse

		if _, err := calendar.GetNearHolidaysInfo(); err != nil {
			log.Fatal(err)
		}

		//fmt.Println(calendar.PublicHolidaysResponse)
	}
}
