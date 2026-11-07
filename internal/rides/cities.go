package rides

import "strings"

// CityNames — shahar raqami bo'yicha uning nomi.
var CityNames = map[int64]string{
	CityTashkent:  "Toshkent",
	CitySamarkand: "Samarqand",
	CityBukhara:   "Buxoro",
	CityFergana:   "Farg'ona",
	CityAndijan:   "Andijon",
	CityNamangan:  "Namangan",
}

// cityAliases — foydalanuvchi yozishi mumkin bo'lgan har bir shakl → shahar
// raqami. Kalitlar kichik harfda va oddiy apostrof bilan.
var cityAliases = map[string]int64{}

func init() {
	for id, name := range CityNames {
		key := strings.ToLower(name)
		cityAliases[key] = id
		cityAliases[strings.ReplaceAll(key, "'", "")] = id // fargona
	}
}

// apostrophes — klaviaturalarda uchraydigan apostrof turlari: ʻ (U+02BB),
// ’ (U+2019) va `. Hammasi oddiy ' ga aylantiriladi.
var apostrophes = strings.NewReplacer("ʻ", "'", "’", "'", "`", "'")

// ParseCity shahar nomini raqamga aylantiradi. Katta-kichik harf, chetdagi
// bo'shliqlar va apostrof turi ahamiyatsiz.
func ParseCity(name string) (int64, bool) {
	key := apostrophes.Replace(strings.ToLower(strings.TrimSpace(name)))
	id, ok := cityAliases[key]
	return id, ok
}
