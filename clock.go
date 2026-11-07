package main

import (
	"fmt"
	"os"
	"time"
)

// Tashkent — O'zbekiston vaqti (UTC+5, yozgi vaqt yo'q). Kompyuter qaysi
// mintaqada bo'lmasin, Ulov24 vaqtni shu bo'yicha o'qiydi va ko'rsatadi.
var Tashkent = time.FixedZone("Asia/Tashkent", 5*60*60)

// now hozirgi vaqtni qaytaradi. ULOV24_NOW berilgan bo'lsa (masalan
// 2026-03-20T06:00:00+05:00), o'sha vaqt olinadi: testlar va namoyish uchun.
func now() time.Time {
	value := os.Getenv("ULOV24_NOW")
	if value == "" {
		return time.Now().In(Tashkent)
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		fmt.Println("ULOV24_NOW noto'g'ri:", value)
		os.Exit(1)
	}
	return t.In(Tashkent)
}
