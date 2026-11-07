// Ulov24 yo'lovchilarni o'sha tomonga ketayotgan haydovchilar bilan bog'laydi.
package main

import (
	"fmt"
	"os"
	"strconv"
)

// Version — dasturning qaysi bosqichda ekani (git tag bilan bir xil).
const Version = "c1-m2"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Ishlatish: ulov24 <buyruq>")
		fmt.Println("Buyruqlar: version, modes, price <son>")
		return
	}

	switch os.Args[1] {
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
