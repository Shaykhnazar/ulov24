package main

// TripMode — safar qachonligi. Qiymatlar matn: ma'lumotlar bazasida ham aynan
// shu so'zlar saqlanadi, shuning uchun iota emas.
type TripMode string

const (
	TripModePlanned TripMode = "planned" // oldindan rejalashtirilgan
	TripModeToday   TripMode = "today"   // bugun
	TripModeNow     TripMode = "now"     // hozir, yaqin 3 soat ichida
)

// Label — foydalanuvchiga ko'rsatiladigan nom.
func (m TripMode) Label() string {
	switch m {
	case TripModeNow:
		return "hozir"
	case TripModeToday:
		return "bugun"
	case TripModePlanned:
		return "rejali"
	default:
		return string(m)
	}
}

// Priority — lentada qaysi rejim yuqorida turadi: kattasi oldinda.
func (m TripMode) Priority() int {
	switch m {
	case TripModeNow:
		return 3
	case TripModeToday:
		return 2
	case TripModePlanned:
		return 1
	default:
		return 0
	}
}

// MaxSeats — bitta mashinada ko'pi bilan nechta joy.
const MaxSeats = 8
