package main

import (
	"fmt"
	ph "github.com/catdog93/test-task/publicHolidays"
	"log"
)

//App prints info about whether it’s a holiday today (and the name of it) or not.

//If today isn’t a holiday, the application should print the next closest holiday.
//For example, "The next holiday is International Workers' Day, May 1, and the weekend will last 3 days: May 1 - May 3".
func main() {
	calendar := &ph.Calendar{}
	err := calendar.InitHolidaysCalendar()
	if err != nil {
		log.Fatal(err)
	}
	if len(calendar.Holidays) == 0 {
		fmt.Println("No more holidays this year !")
		return
	}
	holiday, err := calendar.IsHolidayToday()
	if err != nil {
		log.Fatal(err)
	}
	if holiday != nil {
		fmt.Printf("Holiday: %s is today !", holiday.Name)
		return
	}
	if info, err := calendar.GetNearHolidaysInfo(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(info)
	}
}
