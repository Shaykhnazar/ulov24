// Ulov24 yo'lovchilarni o'sha tomonga ketayotgan haydovchilar bilan bog'laydi.
package main

import (
	"fmt"
	"os"
)

// Version — dasturning qaysi bosqichda ekani (git tag bilan bir xil).
const Version = "c1-m1"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Ishlatish: ulov24 <buyruq>")
		fmt.Println("Buyruqlar: version")
		return
	}

	switch os.Args[1] {
	case "version":
		fmt.Println("ulov24", Version)
	default:
		fmt.Println("Noma'lum buyruq:", os.Args[1])
		os.Exit(1)
	}
}
