package main

import (
	"fmt"
	hol "github.com/catdog93/test-task/holidaysService"
	"log"
)

func main() {
	// input
	calendar := &hol.HolidaysCalendar{}
	//calendar.GetHolidaysForThisYear(2019, "UA")*/

	//prints if it’s a holiday today (and the name of it). If today isn’t a holiday, the application should print the next closest holiday.

	if sameTime, holidayName, err := calendar.IsHolidayToday(); err != nil {
		log.Fatal(err)
	} else {
		if sameTime {
			fmt.Printf("Holiday: %s is today !", holidayName)
		} else {
			if info, err := calendar.GetNearHolidaysInfo(); err != nil {
				log.Fatal(err)
			} else {
				fmt.Println(info)
				/*	fmt.Println(calendar.HolidaysNamesList)
					fmt.Println(calendar.HolidaysDatesList)
					fmt.Println(calendar.HolidaysDifferenceSinceNowList)*/
			}
		}
	}

	//The next holiday is International Workers' Day, May 1, and the weekend will last 3 days: May 1 - May 3
}
