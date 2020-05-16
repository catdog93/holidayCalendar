package main

import (
	"fmt"
	"github.com/catdog93/test-task/holidaysService"
)

func main() {
	calendar := holidaysService.HolidaysCalendar{}
	//calendar.GetHolidaysForThisYear(2019, "UA")
	fmt.Println(calendar.GetCountriesCodes())
	fmt.Println(calendar.TempCache)
}
