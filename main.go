// Ulov24 yo'lovchilarni o'sha tomonga ketayotgan haydovchilar bilan bog'laydi.
package main

import (
	"fmt"
	"os"
	"strconv"
)

// Version — dasturning qaysi bosqichda ekani (git tag bilan bir xil).
const Version = "c1-m5"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Ishlatish: ulov24 <buyruq>")
		fmt.Println("Buyruqlar: version, modes, price <son>, cities, routes, feed [qayerdan qayerga], match")
		return
	}

	switch os.Args[1] {
	case "routes":
		current := now()
		routes, _ := sampleData(current)
		for _, r := range routes {
			label := r.TripMode.Label()
			if isRouteExpired(r, current) {
				label = "eskirgan"
			}
			fmt.Println(routeLine(r, label))
		}
	case "cities":
		for id := CityTashkent; id <= CityNamangan; id++ {
			fmt.Println(id, CityNames[id])
		}
	case "feed":
		current := now()
		routes, _ := sampleData(current)
		if len(os.Args) >= 4 {
			from, ok := ParseCity(os.Args[2])
			if !ok {
				fmt.Println("Noma'lum shahar:", os.Args[2])
				os.Exit(1)
			}
			to, ok := ParseCity(os.Args[3])
			if !ok {
				fmt.Println("Noma'lum shahar:", os.Args[3])
				os.Exit(1)
			}
			routes = IndexByDirection(routes)[Direction{From: from, To: to}]
		}
		for _, r := range Feed(routes, current) {
			fmt.Println(routeLine(r, r.TripMode.Label()))
		}
	case "match":
		current := now()
		routes, requests := sampleData(current)
		for _, q := range requests {
			sameWay := directionFilter(q.FromLocationID, q.ToLocationID)
			found := ""
			for _, r := range routes {
				if !sameWay(r) || !compatible(r, q, current) {
					continue
				}
				if found != "" {
					found += ", "
				}
				found += fmt.Sprintf("#%d", r.ID)
			}
			if found == "" {
				fmt.Printf("So'rov #%d (%d kishi): mos yo'nalish yo'q\n", q.ID, q.PassengerCount)
			} else {
				fmt.Printf("So'rov #%d (%d kishi): yo'nalish %s\n", q.ID, q.PassengerCount, found)
			}
		}
	case "version":
		fmt.Println("ulov24", Version)
	case "modes":
		fmt.Println(TripModeNow, "—", TripModeNow.Label())
		fmt.Println(TripModeToday, "—", TripModeToday.Label())
		fmt.Println(TripModePlanned, "—", TripModePlanned.Label())
	case "price":
		if len(os.Args) < 3 {
			fmt.Println("Narxni yozing: ulov24 price 150000")
			os.Exit(1)
		}
		amount, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Narx butun son bo'lishi kerak:", os.Args[2])
			os.Exit(1)
		}
		if amount < 0 {
			fmt.Println("Narx manfiy bo'lishi mumkin emas:", amount)
			os.Exit(1)
		}
		fmt.Println(FormatSum(amount))
	default:
		fmt.Println("Noma'lum buyruq:", os.Args[1])
		os.Exit(1)
	}
}
