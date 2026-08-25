# Gold Gym POS — Backend (Portfolio)

> ⚠️ **Catatan:** Repo ini adalah versi **portofolio/showcase** dari proyek
> Gold Gym POS. Kode di sini dipublikasikan untuk keperluan demo dan tidak
> dipakai untuk deployment produksi. Semua kredensial (password DB, Redis,
> webhook, SMTP, dsb.) sudah diganti dengan placeholder `CHANGE_ME` — repo
> asli yang dipakai untuk deployment nyata disimpan terpisah secara privat.

Backend Go untuk sistem POS Gold Gym: manajemen item/stok multi-outlet,
penjualan (sales), booking, diskon, dan sinkronisasi data via Kafka/Debezium.

## Tech stack

- Go 1.24, [Gin](https://gin-gonic.com/) (HTTP), [GORM](https://gorm.io/) + SQLX (DB), gRPC
- MySQL/MariaDB, Redis, Kafka + Debezium (CDC), Elasticsearch (opsional)
- Deploy: Docker, Helm charts (`charts/`), Jenkins CI/CD

## Struktur

```
cmd/http/main.go          entrypoint HTTP + gRPC server
internal/entity/          model/struct domain
internal/data/            akses data (query DB, Redis, dst.)
internal/service/         business logic
internal/delivery/http/   handler HTTP (Gin)
internal/boot/            wiring/bootstrap service
files/etc/gold-gym-be/    file konfigurasi per environment
files/migrations/         migrasi SQL
proto/                    definisi gRPC
charts/                   Helm chart untuk deploy ke Kubernetes
```

## Menjalankan secara lokal

1. **Siapkan dependency lokal** (MySQL, Kafka, Zookeeper, Debezium, Elasticsearch)
   lewat Docker Compose:
   ```bash
   docker compose up -d
   ```

2. **Lengkapi konfigurasi.** File di `files/etc/gold-gym-be/gold-gym-be.*.yaml`
   berisi placeholder `CHANGE_ME` untuk semua kredensial (DB, Redis, webhook
   Discord, dst.) — isi dengan nilai environment kamu sendiri sebelum
   menjalankan service. Untuk pengembangan lokal, pakai
   `gold-gym-be.development.yaml` sebagai acuan.

3. **Jalankan migrasi** SQL dari `files/migrations/` ke database yang sudah
   disiapkan.

4. **Jalankan service:**
   ```bash
   go run cmd/http/main.go
   ```
   Default: HTTP di port `8080`/`8085` (tergantung environment), gRPC di
   port `50051`. Dokumentasi API tersedia di `docs/swagger.yaml`.

5. **Testing:**
   ```bash
   go test ./... -v
   ```

## Environment

Pilih file config sesuai environment lewat variabel `CONFIG_PATH`
(lihat `charts/values.yaml` / `charts/values-production.yaml` untuk contoh
di Kubernetes):

| Environment | File config |
|---|---|
| Local/dev | `files/etc/gold-gym-be/gold-gym-be.development.yaml` |
| Staging | `files/etc/gold-gym-be/gold-gym-be.staging.yaml` |
| Production | `files/etc/gold-gym-be/gold-gym-be.production.yaml` |

Semua field kredensial di ketiga file di atas berisi `CHANGE_ME` — **wajib**
diganti dengan nilai asli sebelum dipakai untuk menjalankan service apa pun.
