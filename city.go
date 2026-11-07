package main

// Shaharlar raqami. Keyinroq ular ma'lumotlar bazasidagi locations jadvalining
// id'lari bo'ladi.
const (
	CityTashkent  int64 = 1
	CitySamarkand int64 = 2
	CityBukhara   int64 = 3
	CityFergana   int64 = 4
	CityAndijan   int64 = 5
	CityNamangan  int64 = 6
)

// cityName shahar raqamini nomga aylantiradi.
func cityName(id int64) string {
	switch id {
	case CityTashkent:
		return "Toshkent"
	case CitySamarkand:
		return "Samarqand"
	case CityBukhara:
		return "Buxoro"
	case CityFergana:
		return "Farg'ona"
	case CityAndijan:
		return "Andijon"
	case CityNamangan:
		return "Namangan"
	default:
		return "?"
	}
}
