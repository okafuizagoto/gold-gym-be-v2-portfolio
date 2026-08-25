package graphql

import (
	"fmt"

	goldEntity "gold-gym-be/internal/entity/goldgym"

	"github.com/graphql-go/graphql"
)

// ─── QUERY RESOLVERS ──────────────────────────────────────────────────────────

// resolveGetGoldUsers → goldUsers: [GoldUser]
func resolveGetGoldUsers(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		users, err := svc.GetGoldUser(p.Context)
		if err != nil {
			return nil, err
		}

		result := make([]interface{}, 0, len(users))
		for _, u := range users {
			result = append(result, map[string]interface{}{
				"gold_id":                u.GoldId,
				"gold_email":             u.GoldEmail,
				"gold_password":          u.GoldPassword,
				"gold_nama":              u.GoldNama,
				"gold_nomorhp":           u.GoldNomorHp,
				"gold_nomorkartu":        u.GoldNomorKartu,
				"gold_cvv":               u.GoldCvv,
				"gold_expireddate":       u.GoldExpireddate,
				"gold_namapemegangkartu": u.GoldPemegangKartu,
			})
		}
		return result, nil
	}
}

// resolveGoldUserByEmail → goldUserByEmail(email): String
func resolveGoldUserByEmail(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		email, _ := p.Args["email"].(string)
		return svc.GetGoldUserByEmail(p.Context, email)
	}
}

// resolveGoldUserDataByEmail → goldUserDataByEmail(email): GoldUserDetail
func resolveGoldUserDataByEmail(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		email, _ := p.Args["email"].(string)
		u, err := svc.GetGoldUserDataByEmail(p.Context, email)
		if err != nil {
			return nil, err
		}
		token := ""
		if !u.GoldToken.IsZero() {
			token = u.GoldToken.String
		}
		return map[string]interface{}{
			"gold_id":                    u.GoldId,
			"gold_email":                 u.GoldEmail,
			"gold_password":              u.GoldPassword,
			"gold_nama":                  u.GoldNama,
			"gold_nomorhp":               u.GoldNomorHp,
			"gold_nomorkartu":            u.GoldNomorKartu,
			"gold_cvv":                   u.GoldCvv,
			"gold_expireddate":           u.GoldExpireddate,
			"gold_namapemegangkartu":     u.GoldPemegangKartu,
			"gold_validasiyn":            u.GoldValidasiYN,
			"gold_token":                 token,
			"gold_otp":                   u.GoldOtp,
			"gold_updated_by":            u.GoldUpdatedBy,
			"gold_updated_at":            u.GoldUpdatedAt,
			"gold_last_login":            u.GoldLastLogin,
			"gold_last_login_host":       u.GoldLastLoginHost,
			"gold_force_change_password": u.GoldForceChangePassword,
		}, nil
	}
}

// resolveAllSubscriptions → allSubscriptions: [Subscription]
func resolveAllSubscriptions(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		subs, err := svc.GetAllSubscription(p.Context)
		if err != nil {
			return nil, err
		}
		result := make([]interface{}, 0, len(subs))
		for _, s := range subs {
			result = append(result, map[string]interface{}{
				"gold_namapaket":       s.GoldNamaPaket,
				"gold_namalayanan":     s.GoldNamaLayanan,
				"gold_harga":           s.GoldHarga,
				"gold_jadwal":          s.GoldJadwal,
				"gold_listlatihan":     s.GoldListLatihan,
				"gold_jumlahpertemuan": s.GoldJumlahpertemuan,
				"gold_durasi":          s.GoldDurasi,
			})
		}
		return result, nil
	}
}

// resolveSubsWithUser → subsWithUser: [SubsWithUser]
func resolveSubsWithUser(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		list, err := svc.GetSubsWithUser(p.Context)
		if err != nil {
			return nil, err
		}
		result := make([]interface{}, 0, len(list))
		for _, u := range list {
			menuid := ""
			if !u.GoldMenuId.IsZero() {
				menuid = u.GoldMenuId.String
			}
			namapaket := ""
			if !u.GoldNamaPaket.IsZero() {
				namapaket = u.GoldNamaPaket.String
			}
			namalayanan := ""
			if !u.GoldNamaLayanan.IsZero() {
				namalayanan = u.GoldNamaLayanan.String
			}
			harga := float64(0)
			if !u.GoldHarga.IsZero() {
				harga = u.GoldHarga.Float64
			}
			listlatihan := ""
			if !u.GoldListLatihan.IsZero() {
				listlatihan = u.GoldListLatihan.String
			}
			jumlahpertemuan := int64(0)
			if !u.GoldJumlahpertemuan.IsZero() {
				jumlahpertemuan = u.GoldJumlahpertemuan.Int64
			}
			durasi := int64(0)
			if !u.GoldDurasi.IsZero() {
				durasi = u.GoldDurasi.Int64
			}
			statuslangganan := ""
			if !u.GoldStatuslangganan.IsZero() {
				statuslangganan = u.GoldStatuslangganan.String
			}
			result = append(result, map[string]interface{}{
				"gold_id":              u.GoldId,
				"gold_menuid":          menuid,
				"gold_email":           u.GoldEmail,
				"gold_nama":            u.GoldNama,
				"gold_nomorhp":         u.GoldNomorHp,
				"gold_expireddate":     u.GoldExpireddate,
				"gold_namapaket":       namapaket,
				"gold_namalayanan":     namalayanan,
				"gold_harga":           harga,
				"gold_listlatihan":     listlatihan,
				"gold_jumlahpertemuan": int(jumlahpertemuan),
				"gold_durasi":          int(durasi),
				"gold_statuslangganan": statuslangganan,
			})
		}
		return result, nil
	}
}

// resolveSubsTotalHarga → subscriptionTotalHarga(email): SubscriptionPayment
func resolveSubsTotalHarga(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		email, _ := p.Args["email"].(string)
		payment, err := svc.GetSubscriptionHeaderTotalHarga(p.Context, email)
		if err != nil {
			return nil, err
		}
		totalharga := float64(0)
		if !payment.GoldTotalharga.IsZero() {
			totalharga = payment.GoldTotalharga.Float64
		}
		return map[string]interface{}{
			"gold_totalharga":  totalharga,
			"gold_statusbayar": payment.GoldValidasiPayment,
		}, nil
	}
}

// ─── INSERT RESOLVERS ─────────────────────────────────────────────────────────

// resolveInsertGoldUser → insertGoldUser(input): MutationResult
func resolveInsertGoldUser(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		inp, _ := p.Args["input"].(map[string]interface{})
		user := goldEntity.GetGoldUsers{
			GoldEmail:         strVal(inp, "gold_email"),
			GoldPassword:      strVal(inp, "gold_password"),
			GoldNama:          strVal(inp, "gold_nama"),
			GoldNomorHp:       strVal(inp, "gold_nomorhp"),
			GoldNomorKartu:    strVal(inp, "gold_nomorkartu"),
			GoldCvv:           strVal(inp, "gold_cvv"),
			GoldExpireddate:   strVal(inp, "gold_expireddate"),
			GoldPemegangKartu: strVal(inp, "gold_namapemegangkartu"),
		}
		result, err := svc.InsertGoldUser(p.Context, user)
		if err != nil {
			return mutationFail(err.Error()), nil
		}
		msg := fmt.Sprintf("%v", result)
		success := msg == "Sukses"
		return map[string]interface{}{"success": success, "message": msg}, nil
	}
}

// resolveInsertSubscription → insertSubscription(input): MutationResult
func resolveInsertSubscription(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		inp, _ := p.Args["input"].(map[string]interface{})
		email := strVal(inp, "gold_email")
		rawIDs, _ := inp["menu_ids"].([]interface{})

		detailData := make([]goldEntity.SubscriptionDetail, 0, len(rawIDs))
		for _, id := range rawIDs {
			menuID, _ := id.(int)
			detailData = append(detailData, goldEntity.SubscriptionDetail{GoldMenuId: menuID})
		}
		subs := goldEntity.InsertSubsAll{
			HeaderData: goldEntity.SubscriptionAll{GoldEmail: email},
			DetailData: detailData,
		}
		result, err := svc.InsertSubscriptionUser(p.Context, subs)
		if err != nil {
			return mutationFail(err.Error()), nil
		}
		return map[string]interface{}{"success": result == "Berhasil", "message": result}, nil
	}
}

// resolveInsertSubsDetail → insertSubscriptionDetail(input): MutationResult
func resolveInsertSubsDetail(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		inp, _ := p.Args["input"].(map[string]interface{})
		detail := goldEntity.SubscriptionDetail{
			GoldId:     intVal(inp, "gold_id"),
			GoldMenuId: intVal(inp, "gold_menuid"),
		}
		result, err, _ := svc.InsertSubscriptionDetail(p.Context, detail)
		if err != nil {
			return mutationFail(err.Error()), nil
		}
		return map[string]interface{}{"success": result == "Berhasil", "message": result}, nil
	}
}

// ─── UPDATE RESOLVERS ─────────────────────────────────────────────────────────

// resolveUpdatePassword → updatePassword(input): MutationResult
func resolveUpdatePassword(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		inp, _ := p.Args["input"].(map[string]interface{})
		subs := goldEntity.UpdatePassword{
			GoldEmail:    strVal(inp, "gold_email"),
			GoldPassword: strVal(inp, "gold_password"),
			GoldOTP:      strVal(inp, "gold_otp"),
		}
		result, err := svc.UpdateDataPeserta(p.Context, subs)
		if err != nil {
			return mutationFail(err.Error()), nil
		}
		return map[string]interface{}{"success": result == "Berhasil", "message": result}, nil
	}
}

// resolveUpdateNama → updateNama(input): MutationResult
func resolveUpdateNama(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		inp, _ := p.Args["input"].(map[string]interface{})
		subs := goldEntity.UpdateNama{
			GoldEmail: strVal(inp, "gold_email"),
			GoldNama:  strVal(inp, "gold_nama"),
		}
		result, err := svc.UpdateNama(p.Context, subs)
		if err != nil {
			return mutationFail(err.Error()), nil
		}
		return map[string]interface{}{"success": result == "Berhasil", "message": result}, nil
	}
}

// resolveUpdateKartu → updateKartu(input): MutationResult
func resolveUpdateKartu(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		inp, _ := p.Args["input"].(map[string]interface{})
		subs := goldEntity.UpdateKartu{
			GoldEmail:      strVal(inp, "gold_email"),
			GoldNomorKartu: strVal(inp, "gold_nomorkartu"),
			GoldCvv:        strVal(inp, "gold_cvv"),
		}
		result, err := svc.UpdateKartu(p.Context, subs)
		if err != nil {
			return mutationFail(err.Error()), nil
		}
		return map[string]interface{}{"success": result == "Berhasil", "message": result}, nil
	}
}

// resolveUpdateSubsDetail → updateSubscriptionDetail(input): MutationResult
func resolveUpdateSubsDetail(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		inp, _ := p.Args["input"].(map[string]interface{})
		subs := goldEntity.UpdateSubs{
			GoldId:              intVal(inp, "gold_id"),
			GoldMenuId:          intVal(inp, "gold_menuid"),
			GoldJumlahpertemuan: intVal(inp, "gold_jumlahpertemuan"),
		}
		result, err := svc.UpdateSubscriptionDetail(p.Context, subs)
		if err != nil {
			return mutationFail(err.Error()), nil
		}
		return map[string]interface{}{"success": result == "Berhasil", "message": result}, nil
	}
}

// ─── OTP RESOLVERS ────────────────────────────────────────────────────────────

// resolveValidateOTP → validateOTP(otp, email): MutationResult
func resolveValidateOTP(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		otp, _ := p.Args["otp"].(string)
		email, _ := p.Args["email"].(string)
		result, err := svc.UpdateValidationOTP(p.Context, otp, email)
		if err != nil {
			return mutationFail(err.Error()), nil
		}
		return map[string]interface{}{"success": result == "Berhasil", "message": result}, nil
	}
}

// resolveSendOTP → sendOTP(email): MutationResult
func resolveSendOTP(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		email, _ := p.Args["email"].(string)
		result, err := svc.UpdateOTP(p.Context, email)
		if err != nil {
			return mutationFail(err.Error()), nil
		}
		return map[string]interface{}{"success": result == "Berhasil", "message": result}, nil
	}
}

// resolveSendOTPSubscription → sendOTPSubscription(email): MutationResult
func resolveSendOTPSubscription(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		email, _ := p.Args["email"].(string)
		result, err := svc.UpdateOTPSubscription(p.Context, email)
		if err != nil {
			return mutationFail(err.Error()), nil
		}
		return map[string]interface{}{"success": result == "Berhasil", "message": result}, nil
	}
}

// resolveProcessPayment → processPayment(otp, email): MutationResult
func resolveProcessPayment(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		otp, _ := p.Args["otp"].(string)
		email, _ := p.Args["email"].(string)
		result, err, _ := svc.UpdatePayment(p.Context, otp, email)
		if err != nil {
			return mutationFail(err.Error()), nil
		}
		return map[string]interface{}{"success": result == "OTP true", "message": result}, nil
	}
}

// ─── AUTH RESOLVERS ───────────────────────────────────────────────────────────

// resolveLoginUser → loginUser(email, password): LoginResult
func resolveLoginUser(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		email, _ := p.Args["email"].(string)
		password, _ := p.Args["password"].(string)

		token, metadata, _, err := svc.LoginUser(p.Context, email, password, "graphql", "")
		if err != nil {
			return nil, err
		}

		username := ""
		if u, ok := metadata["username"]; ok {
			username = fmt.Sprintf("%v", u)
		}

		return map[string]interface{}{
			"access_token":          token.AccessToken,
			"refresh_token":         token.RefreshToken,
			"expires_in":            int(token.ExpiresIn),
			"expires_at":            int(token.ExpiresAt),
			"token_type":            token.TokenType,
			"force_change_password": token.ForceChangePassword,
			"username":              username,
		}, nil
	}
}

// resolveLogout → logout(email): MutationResult
func resolveLogout(svc IgoldgymSvc) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		email, _ := p.Args["email"].(string)
		result, err := svc.Logout(p.Context, goldEntity.Logout{GoldEmail: email})
		if err != nil {
			return mutationFail(err.Error()), nil
		}
		return map[string]interface{}{"success": result == "Berhasil", "message": result}, nil
	}
}

// ─── HELPERS ──────────────────────────────────────────────────────────────────

func strVal(m map[string]interface{}, key string) string {
	v, _ := m[key].(string)
	return v
}

func intVal(m map[string]interface{}, key string) int {
	v, _ := m[key].(int)
	return v
}

func mutationFail(msg string) map[string]interface{} {
	return map[string]interface{}{"success": false, "message": msg}
}
