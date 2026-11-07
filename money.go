package main

import "strconv"

// FormatSum narxni o'qishga qulay ko'rinishga keltiradi: 150000 → "150 000 so'm".
func FormatSum(amount int) string {
	digits := strconv.Itoa(amount)
	result := ""
	for i, d := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			result += " "
		}
		result += string(d)
	}
	return result + " so'm"
}
