package main

import (
	"fmt"

	"github.com/Shaykhnazar/ulov24/internal/rides"
)

// routeLine — bitta yo'nalishning lentadagi qatori.
func routeLine(r rides.DriverRoute, label string) string {
	return fmt.Sprintf("#%d %s → %s · %s · %d joy · %s · %s",
		r.ID, rides.CityNames[r.FromLocationID], rides.CityNames[r.ToLocationID],
		r.DepartureAt.In(rides.Tashkent).Format("02.01 15:04"),
		r.SeatsRemaining, priceText(r.PricePerSeat), label)
}

// priceText: narx bo'lsa "150 000 so'm", bo'lmasa "narx kelishiladi".
func priceText(price *int) string {
	if price == nil {
		return "narx kelishiladi"
	}
	return rides.FormatSum(*price)
}
