package rides

// Page — uzun ro'yxatning bir bo'lagi. T istalgan tur bo'lishi mumkin:
// yo'nalishlar, so'rovlar, mosliklar — bitta kod hammasiga ishlaydi.
type Page[T any] struct {
	Items  []T
	Total  int
	Offset int
	Limit  int
}

// Paginate items'dan offset'dan boshlab ko'pi bilan limit ta elementni oladi.
func Paginate[T any](items []T, offset, limit int) Page[T] {
	offset = max(offset, 0)
	limit = max(limit, 0)
	start := min(offset, len(items))
	end := min(start+limit, len(items))
	return Page[T]{Items: items[start:end], Total: len(items), Offset: offset, Limit: limit}
}
