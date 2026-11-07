# Ulov24

Ulov24 yo'lovchilarni shu tomonga allaqachon ketayotgan haydovchilar bilan
bog'laydi: haydovchi yo'nalish va bo'sh joylarini e'lon qiladi, yo'lovchi
so'rov qoldiradi, tizim ularni moslaydi. Ulov24 tashuvchi ham, taksi xizmati
ham emas — u faqat bir tomonga ketayotgan odamlarni tanishtiradi.

Bu repozitoriy Maqsadly'dagi Go kurslari davomida noldan quriladi. Har bir
modul oxirida bitta git tag qo'yiladi; istalgan bosqichni ochib, ishga tushirib
ko'rishingiz mumkin.

## Ishga tushirish

Go 1.26 yoki yangiroq kerak.

```bash
go run ./cmd/ulov24 version
```

Testlar:

```bash
go test ./...
```

## Bosqichlar

| Tag | Kurs, modul | Nima qo'shildi |
|---|---|---|
| `c1-m1` | Go asoslari, 1-modul | Repozitoriy, birinchi dastur: `ulov24 version` |
| `c1-m2` | Go asoslari, 2-modul | Safar rejimlari, shaharlar, narx: `ulov24 modes`, `ulov24 price 150000` |
| `c1-m3` | Go asoslari, 3-modul | Namuna safarlarni moslash: `ulov24 routes`, `ulov24 match` |
| `c1-m4` | Go asoslari, 4-modul | Moslash qoidalari funksiyalarga, tekshiruvlar `error` qaytaradi, closure |
| `c1-m5` | Go asoslari, 5-modul | Shaharlar map'i, tomon bo'yicha indeks, lenta: `ulov24 cities`, `ulov24 feed toshkent samarqand` |
| `c1-m6` | Go asoslari, 6-modul | Metodlar, pointerlar, ikki tomonlama rozilik va joy band qilish: `ulov24 accept-demo` |
| `c1-m7` | Go asoslari, 7-modul | Paketlar (`internal/rides`, `internal/ridematching`, `internal/store`), xatolar `errors.Is` bilan, JSON fayl: `route add`, `request add`, `match`, `accept` |

Bosqichni ochish: `git checkout c1-m1`.

## Litsenziya

MIT — [LICENSE](LICENSE).
