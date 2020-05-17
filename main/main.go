package main

import (
	"github.com/catdog93/test-task/holidaysService"
)

func main() {
	// input
	calendar := holidaysService.HolidaysCalendar{}
	calendar.GetHolidaysForThisYear(2019, "UA")
}
