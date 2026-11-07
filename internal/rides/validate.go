package rides

import (
	"fmt"
	"time"
)

// ValidateSeats: joylar (yoki yo'lovchilar) soni 1 dan MaxSeats gacha.
func ValidateSeats(n int) error {
	if n < 1 || n > MaxSeats {
		return fmt.Errorf("joylar soni 1 dan %d gacha bo'lishi kerak: %d", MaxSeats, n)
	}
	return nil
}

// ValidateTripMode jo'nash vaqti tanlangan rejimga mosligini tekshiradi.
func ValidateTripMode(mode TripMode, departure, now time.Time) error {
	switch mode {
	case TripModeNow:
		if departure.Before(now) || departure.After(now.Add(3*time.Hour)) {
			return fmt.Errorf(`"hozir" rejimida jo'nash 3 soat ichida bo'lishi kerak`)
		}
	case TripModeToday:
		if departure.Before(now) {
			return fmt.Errorf("jo'nash vaqti o'tib ketgan")
		}
		dy, dm, dd := departure.In(Tashkent).Date()
		ny, nm, nd := now.In(Tashkent).Date()
		if dy != ny || dm != nm || dd != nd {
			return fmt.Errorf(`"bugun" rejimida jo'nash bugun bo'lishi kerak`)
		}
	case TripModePlanned:
		if !departure.After(now) {
			return fmt.Errorf(`"rejali" rejimida jo'nash kelajakda bo'lishi kerak`)
		}
	default:
		return fmt.Errorf("noma'lum rejim: %s", mode)
	}
	return nil
}
