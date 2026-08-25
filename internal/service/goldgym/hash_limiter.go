package goldgym

import "github.com/raja/argon2pw"

// hashSemaphore membatasi jumlah operasi argon2 (hash & compare) yang boleh
// berjalan bersamaan per pod. argon2 sengaja CPU+memory-hard (parallelism=4,
// 32MB per operasi) -- tanpa batas ini, lonjakan request (misal badai
// pendaftaran) bisa menghabiskan CPU/RAM pod sampai readiness probe timeout
// dan pod restart (insiden 2026-08-18: 50 pendaftaran bersamaan -> kedua
// pod production restart). Kekuatan hash tidak berubah sama sekali, cuma
// antre kalau slot penuh.
var hashSemaphore = make(chan struct{}, 2)

func hashPassword(plain string) (string, error) {
	hashSemaphore <- struct{}{}
	defer func() { <-hashSemaphore }()
	return argon2pw.GenerateSaltedHash(plain)
}

func compareHash(hash, plain string) (bool, error) {
	hashSemaphore <- struct{}{}
	defer func() { <-hashSemaphore }()
	return argon2pw.CompareHashWithPassword(hash, plain)
}
