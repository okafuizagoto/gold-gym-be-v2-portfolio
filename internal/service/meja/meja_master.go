package goldgym

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	goldMejaEntity "gold-gym-be/internal/entity/meja"

	"github.com/pkg/errors"
)

var trailingDigitPattern = regexp.MustCompile(`^(.*?)(\d+)$`)

// GenerateSequentialNames menghasilkan `count` nama meja berurutan dari
// `start`, mempertahankan prefix huruf dan zero-padding angka di akhir.
// "A1" + 3      -> A1, A2, A3
// "MEJA010" + 3 -> MEJA010, MEJA011, MEJA012
// "VIP" (tanpa angka di akhir) -> error, mode bulk butuh nama diakhiri angka.
// Jalur utama generate nama ada di Flutter (untuk preview sebelum submit);
// fungsi ini disediakan untuk potensi reuse Go-side (mis. bulk-import CSV).
func GenerateSequentialNames(start string, count int) ([]string, error) {
	start = strings.TrimSpace(start)
	if start == "" {
		return nil, errors.New("nama awal wajib diisi")
	}
	if count < 1 {
		return nil, errors.New("jumlah meja minimal 1")
	}

	matches := trailingDigitPattern.FindStringSubmatch(start)
	if matches == nil {
		return nil, errors.New("nama awal untuk mode banyak meja harus diakhiri angka, contoh: A1, MEJA01")
	}

	prefix := matches[1]
	digits := matches[2]
	width := len(digits)

	startNum, err := strconv.Atoi(digits)
	if err != nil {
		return nil, errors.New("nama awal untuk mode banyak meja harus diakhiri angka, contoh: A1, MEJA01")
	}

	names := make([]string, 0, count)
	for i := 0; i < count; i++ {
		num := startNum + i
		names = append(names, fmt.Sprintf("%s%0*d", prefix, width, num))
	}
	return names, nil
}

func (s Service) InsertMejaBulk(ctx context.Context, goldid int, outcode string, rows []goldMejaEntity.InsertMeja) error {
	if outcode == "" {
		return errors.New("outlet wajib dipilih")
	}
	if len(rows) == 0 {
		return errors.New("data meja kosong")
	}

	existing, err := s.goldgymmeja.GetExistingMejaNames(ctx, outcode)
	if err != nil {
		return errors.Wrap(err, "[Service][InsertMejaBulk]")
	}

	seenInBatch := make(map[string]bool, len(rows))
	insertRows := make([]goldMejaEntity.Meja, 0, len(rows))

	for _, r := range rows {
		name := strings.TrimSpace(r.MejaName)
		if name == "" {
			return errors.New("nama/nomor meja wajib diisi")
		}
		if r.MejaCapacity < 1 {
			return errors.New("jumlah pelanggan (kapasitas) meja minimal 1")
		}
		if r.MejaAreaID <= 0 {
			return errors.New("area meja wajib dipilih")
		}
		if existing[name] {
			return errors.Errorf("nama meja %q sudah dipakai di outlet ini", name)
		}
		if seenInBatch[name] {
			return errors.Errorf("nama meja %q duplikat dalam permintaan ini", name)
		}
		seenInBatch[name] = true

		insertRows = append(insertRows, goldMejaEntity.Meja{
			MejaGoldID:   goldid,
			MejaOutcode:  outcode,
			MejaAreaID:   r.MejaAreaID,
			MejaName:     name,
			MejaCapacity: r.MejaCapacity,
			MejaStatus:   goldMejaEntity.MejaStatusKosong,
		})
	}

	if err := s.goldgymmeja.InsertMeja(ctx, insertRows); err != nil {
		return errors.Wrap(err, "[Service][InsertMejaBulk]")
	}
	return nil
}

func (s Service) GetMejaByOutlet(ctx context.Context, goldid int, outcode string) ([]goldMejaEntity.Meja, error) {
	rows, err := s.goldgymmeja.GetMejaByOutlet(ctx, goldid, outcode)
	if err != nil {
		return nil, errors.Wrap(err, "[Service][GetMejaByOutlet]")
	}
	return rows, nil
}

// ReserveMeja menandai meja jadi ISI (dipakai picker POS saat kasir
// finalisasi pilihan meja). Atomicity (semua-atau-tidak-sama-sekali)
// ditangani di data layer (lihat ReserveMeja data/meja) lewat SELECT ...
// FOR UPDATE -- kalau tidak semua meja yang diminta masih KOSONG, TIDAK
// ADA yang diubah sama sekali, jadi di sini cukup cek RowsAffected.
func (s Service) ReserveMeja(ctx context.Context, outcode string, mejaIDs []int) error {
	if len(mejaIDs) == 0 {
		return errors.New("meja belum dipilih")
	}

	affected, err := s.goldgymmeja.ReserveMeja(ctx, outcode, mejaIDs)
	if err != nil {
		return errors.Wrap(err, "[Service][ReserveMeja]")
	}

	if affected != int64(len(mejaIDs)) {
		return errors.New("beberapa meja yang dipilih baru saja terisi, silakan pilih ulang")
	}

	return nil
}

func (s Service) ReleaseMeja(ctx context.Context, outcode string, mejaIDs []int) error {
	if len(mejaIDs) == 0 {
		return nil
	}
	if _, err := s.goldgymmeja.UpdateMejaStatus(ctx, outcode, mejaIDs, goldMejaEntity.MejaStatusIsi, goldMejaEntity.MejaStatusKosong); err != nil {
		return errors.Wrap(err, "[Service][ReleaseMeja]")
	}
	return nil
}
