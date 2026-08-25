// package registry

// import (
// 	"context"
// 	"encoding/json"
// 	"log"

// 	"gold-gym-be/internal/entity/goldgym"
// 	"gold-gym-be/internal/resources"
// )

// // HandlerFunc adalah tipe fungsi handler CDC untuk setiap tabel
// type HandlerFunc func(ctx context.Context, op string, after, before map[string]interface{}) error

// // Registry menyimpan daftar handler untuk setiap tabel
// type Registry struct {
// 	handlers map[string]HandlerFunc
// }

// // New membuat Registry baru dengan resource yang dibutuhkan
// func New(res *resources.BootResources) *Registry {
// 	return &Registry{
// 		handlers: GetRegistry(res),
// 	}
// }

// // GetHandler mencari handler berdasarkan nama tabel
// func (r *Registry) GetHandler(table string) (HandlerFunc, bool) {
// 	h, ok := r.handlers[table]
// 	return h, ok
// }

// // GetRegistry daftar semua table CDC → handler function
// func GetRegistry(res *resources.BootResources) map[string]HandlerFunc {
// 	return map[string]HandlerFunc{
// 		"users":   handleUsers(res),
// 		"goldgym": handleGoldGym(res),
// 	}
// }

// // contoh handler untuk tabel users
// func handleUsers(res *resources.BootResources) HandlerFunc {
// 	return func(ctx context.Context, op string, after, before map[string]interface{}) error {
// 		switch op {
// 		case "c":
// 			log.Println("[CDC] users insert event", after)
// 			// TODO: mapping after → entity, panggil res.GoldSvcLocal atau service lain
// 			//---------------------------------------------------------------------
// 			// 1. mapping dari map[string]interface{} → struct entity
// 			var user goldgym.GetGoldUsers
// 			raw, _ := json.Marshal(after)
// 			if err := json.Unmarshal(raw, &user); err != nil {
// 				return err
// 			}

// 			// 2. panggil service yang biasa dipakai di HTTP handler
// 			_, err := res.GoldSvcProd.InsertGoldUser(ctx, user)
// 			return err
// 			//---------------------------------------------------------------------
// 		case "u":
// 			log.Println("[CDC] users update event", after)
// 		case "d":
// 			log.Println("[CDC] users delete event", before)
// 		default:
// 			log.Printf("[CDC] users unknown op=%s", op)
// 		}
// 		return nil
// 	}
// }

// // handler untuk tabel goldgym
// func handleGoldGym(res *resources.BootResources) HandlerFunc {
// 	return func(ctx context.Context, op string, after, before map[string]interface{}) error {
// 		switch op {
// 		case "c":
// 			log.Println("[CDC] goldgym insert event", after)

// 			// mapping CDC → entity yg dipakai di service
// 			var user goldgym.GetGoldUsers
// 			raw, _ := json.Marshal(after)
// 			if err := json.Unmarshal(raw, &user); err != nil {
// 				return err
// 			}

// 			// panggil service yang sama dengan HTTP InsertGoldUser
// 			_, err := res.GoldSvcLocal.InsertGoldUser(ctx, user)
// 			return err

// 		case "u":
// 			log.Println("[CDC] goldgym update event", after)
// 			// TODO: mapping ke entity dan panggil Update service
// 			return nil

// 		case "d":
// 			log.Println("[CDC] goldgym delete event", before)
// 			// TODO: mapping ke entity dan panggil Delete service
// 			return nil

// 		default:
// 			log.Printf("[CDC] goldgym unknown op=%s", op)
// 			return nil
// 		}
// 	}
// }

package registry

import (
	"context"
	"log"

	"gold-gym-be/internal/resources"
)

// HandlerFunc adalah tipe fungsi handler CDC untuk setiap tabel
type HandlerFunc func(ctx context.Context, op string, after, before map[string]interface{}) error

// Registry menyimpan daftar handler untuk setiap tabel
type Registry struct {
	handlers map[string]HandlerFunc
}

// New membuat Registry baru dengan resource yang dibutuhkan
func New(res *resources.BootResources) *Registry {
	return &Registry{
		handlers: GetRegistry(res),
	}
}

// GetHandler mencari handler berdasarkan nama tabel
func (r *Registry) GetHandler(table string) (HandlerFunc, bool) {
	h, ok := r.handlers[table]
	return h, ok
}

// GetRegistry daftar semua table CDC → handler function
func GetRegistry(res *resources.BootResources) map[string]HandlerFunc {
	return map[string]HandlerFunc{
		"users":        handleUsers(res),
		"goldgym":      handleGoldGym(res),
		"data_peserta": handlePeserta(res), // contoh tambahan
	}
}

// ===================== HANDLERS =====================

// handler untuk tabel users
func handleUsers(res *resources.BootResources) HandlerFunc {
	return func(ctx context.Context, op string, after, before map[string]interface{}) error {
		switch op {
		case "c":
			var err error
			// log.Println("[CDC] users insert event", after)

			// var user goldgym.GetGoldUsers
			// raw, _ := json.Marshal(after)
			// if err := json.Unmarshal(raw, &user); err != nil {
			// 	return err
			// }

			// // pake service production
			// _, err := res.GoldSvcProd.InsertGoldUser(ctx, user)
			return err

		case "u":
			log.Println("[CDC] users update event", after)
			// var user goldgym.GetGoldUsers
			// raw, _ := json.Marshal(after)
			// if err := json.Unmarshal(raw, &user); err != nil {
			// 	return err
			// }
			// _, err := res.GoldSvcProd.UpdateGoldUser(ctx, user)
			var err error
			return err

		case "d":
			log.Println("[CDC] users delete event", before)
			// var user goldgym.GetGoldUsers
			// raw, _ := json.Marshal(before)
			// if err := json.Unmarshal(raw, &user); err != nil {
			// 	return err
			// }
			// return res.GoldSvcProd.DeleteGoldUser(ctx, user.ID)
			var err error
			return err

		default:
			log.Printf("[CDC] users unknown op=%s", op)
			return nil
		}
	}
}

// handler untuk tabel goldgym
func handleGoldGym(res *resources.BootResources) HandlerFunc {
	return func(ctx context.Context, op string, after, before map[string]interface{}) error {
		switch op {
		case "c":
			var err error
			// log.Println("[CDC] goldgym insert event", after)

			// var g goldgym.GetGoldUsers
			// raw, _ := json.Marshal(after)
			// if err := json.Unmarshal(raw, &g); err != nil {
			// 	return err
			// }

			// _, err := res.GoldSvcLocal.InsertGoldUser(ctx, g)
			return err

		case "u":
			log.Println("[CDC] goldgym update event", after)
			// sama pola nya
			return nil

		case "d":
			log.Println("[CDC] goldgym delete event", before)
			return nil

		default:
			log.Printf("[CDC] goldgym unknown op=%s", op)
			return nil
		}
	}
}

// handler untuk tabel data_peserta
func handlePeserta(res *resources.BootResources) HandlerFunc {
	return func(ctx context.Context, op string, after, before map[string]interface{}) error {
		var (
		// user goldgym.GetGoldUsers
		// keys []string
		)

		switch op {
		case "c":
			// log.Println("[CDC] peserta insert event", after)
			// log.Println("[CDC] peserta insert event Two", after["gold_cvv"])
			// // for key := range after {
			// // 	keys = append(keys, key)
			// // }
			// // fmt.Println("test-keys", keys)

			// if v, ok := after["gold_expireddate"]; ok {
			// 	switch t := v.(type) {
			// 	case float64:
			// 		// Konversi milisecond → detik → time.Time
			// 		tm := time.UnixMilli(int64(t))
			// 		after["gold_expireddate"] = tm.Format("20060102") // hasil: "20230928"
			// 	}
			// }

			// // 1. Marshal map ke JSON
			// jsonData, err := json.Marshal(after)
			// if err != nil {
			// 	log.Println("Error marshal map:", err)
			// 	return nil
			// }
			// // 2. Unmarshal ke struct
			// err = json.Unmarshal(jsonData, &user)
			// if err != nil {
			// 	log.Println("Error unmarshal ke struct:", err)
			// 	return nil
			// }

			// _, err = res.GoldSvcProd.InsertGoldUser(ctx, user)
			// // TODO: mapping ke entity peserta → panggil service peserta
			return nil

		case "u":
			log.Println("[CDC] peserta update event", after)
			return nil

		case "d":
			log.Println("[CDC] peserta delete event", before)
			return nil

		default:
			log.Printf("[CDC] peserta unknown op=%s", op)
			return nil
		}
	}
}
