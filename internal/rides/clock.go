package rides

import (
	"fmt"
	"os"
	"time"
)

// Tashkent — O'zbekiston vaqti (UTC+5, yozgi vaqt yo'q). Kompyuter qaysi
// mintaqada bo'lmasin, Ulov24 vaqtni shu bo'yicha o'qiydi va ko'rsatadi.
var Tashkent = time.FixedZone("Asia/Tashkent", 5*60*60)

// LocalTimeLayout — foydalanuvchi vaqtni shu ko'rinishda yozadi.
const LocalTimeLayout = "2006-01-02 15:04"

// Now hozirgi vaqtni qaytaradi. ULOV24_NOW berilgan bo'lsa (masalan
// 2026-03-20T06:00:00+05:00), o'sha vaqt olinadi: testlar va namoyish uchun.
func Now() (time.Time, error) {
	value := os.Getenv("ULOV24_NOW")
	if value == "" {
		return time.Now().In(Tashkent), nil
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("ULOV24_NOW: %w", err)
	}
	return t.In(Tashkent), nil
}

// ParseLocalTime "2026-03-20 08:00" ko'rinishidagi Toshkent vaqtini o'qiydi.
func ParseLocalTime(s string) (time.Time, error) {
	return time.ParseInLocation(LocalTimeLayout, s, Tashkent)
}
