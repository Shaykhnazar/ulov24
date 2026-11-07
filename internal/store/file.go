// Package store Ulov24 ma'lumotlarini saqlaydi. Hozircha bitta JSON fayl;
// 2-kursda uning o'rniga PostgreSQL keladi.
package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Shaykhnazar/ulov24/internal/ridematching"
	"github.com/Shaykhnazar/ulov24/internal/rides"
)

// ErrCorruptData — fayl bor, lekin uni JSON sifatida o'qib bo'lmadi.
var ErrCorruptData = errors.New("data file is corrupt")

// Snapshot — fayldagi hamma narsa.
type Snapshot struct {
	Routes   []rides.DriverRoute      `json:"driver_routes"`
	Requests []rides.PassengerRequest `json:"passenger_requests"`
	Matches  []ridematching.RideMatch `json:"ride_matches"`
}

// LoadFile faylni o'qiydi. Fayl hali yo'q bo'lsa — bo'sh Snapshot, xato emas.
func LoadFile(path string) (Snapshot, error) {
	var s Snapshot
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return s, fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("%w: %s: %v", ErrCorruptData, path, err)
	}
	return s, nil
}

// SaveFile avval vaqtinchalik faylga yozadi, keyin uni asl nomga o'tkazadi.
// Yozish yarmida to'xtasa ham asl fayl buzilmaydi.
func SaveFile(path string, s Snapshot) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("replace %s: %w", path, err)
	}
	return nil
}
