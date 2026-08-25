package graphql

import (
	"github.com/graphql-go/graphql"
)

// ─── OUTPUT TYPES ────────────────────────────────────────────────────────────

var goldUserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "GoldUser",
	Fields: graphql.Fields{
		"gold_id":                {Type: graphql.Int},
		"gold_email":             {Type: graphql.String},
		"gold_password":          {Type: graphql.String},
		"gold_nama":              {Type: graphql.String},
		"gold_nomorhp":           {Type: graphql.String},
		"gold_nomorkartu":        {Type: graphql.String},
		"gold_cvv":               {Type: graphql.String},
		"gold_expireddate":       {Type: graphql.String},
		"gold_namapemegangkartu": {Type: graphql.String},
	},
})

var goldUserDetailType = graphql.NewObject(graphql.ObjectConfig{
	Name: "GoldUserDetail",
	Fields: graphql.Fields{
		"gold_id":                      {Type: graphql.Int},
		"gold_email":                   {Type: graphql.String},
		"gold_password":                {Type: graphql.String},
		"gold_nama":                    {Type: graphql.String},
		"gold_nomorhp":                 {Type: graphql.String},
		"gold_nomorkartu":              {Type: graphql.String},
		"gold_cvv":                     {Type: graphql.String},
		"gold_expireddate":             {Type: graphql.String},
		"gold_namapemegangkartu":       {Type: graphql.String},
		"gold_validasiyn":              {Type: graphql.String},
		"gold_token":                   {Type: graphql.String},
		"gold_otp":                     {Type: graphql.String},
		"gold_updated_by":              {Type: graphql.String},
		"gold_updated_at":              {Type: graphql.String},
		"gold_last_login":              {Type: graphql.String},
		"gold_last_login_host":         {Type: graphql.String},
		"gold_force_change_password":   {Type: graphql.Int},
	},
})

var subscriptionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Subscription",
	Fields: graphql.Fields{
		"gold_namapaket":       {Type: graphql.String},
		"gold_namalayanan":     {Type: graphql.String},
		"gold_harga":           {Type: graphql.Float},
		"gold_jadwal":          {Type: graphql.String},
		"gold_listlatihan":     {Type: graphql.String},
		"gold_jumlahpertemuan": {Type: graphql.Int},
		"gold_durasi":          {Type: graphql.Int},
	},
})

var subsWithUserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SubsWithUser",
	Fields: graphql.Fields{
		"gold_id":              {Type: graphql.Int},
		"gold_menuid":          {Type: graphql.String},
		"gold_email":           {Type: graphql.String},
		"gold_nama":            {Type: graphql.String},
		"gold_nomorhp":         {Type: graphql.String},
		"gold_expireddate":     {Type: graphql.String},
		"gold_namapaket":       {Type: graphql.String},
		"gold_namalayanan":     {Type: graphql.String},
		"gold_harga":           {Type: graphql.Float},
		"gold_listlatihan":     {Type: graphql.String},
		"gold_jumlahpertemuan": {Type: graphql.Int},
		"gold_durasi":          {Type: graphql.Int},
		"gold_statuslangganan": {Type: graphql.String},
	},
})

var subscriptionPaymentType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SubscriptionPayment",
	Fields: graphql.Fields{
		"gold_totalharga":      {Type: graphql.Float},
		"gold_statusbayar":     {Type: graphql.String},
	},
})

var loginResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "LoginResult",
	Fields: graphql.Fields{
		"access_token":          {Type: graphql.String},
		"refresh_token":         {Type: graphql.String},
		"expires_in":            {Type: graphql.Int},
		"expires_at":            {Type: graphql.Int},
		"token_type":            {Type: graphql.String},
		"force_change_password": {Type: graphql.Int},
		"username":              {Type: graphql.String},
	},
})

var mutationResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "MutationResult",
	Fields: graphql.Fields{
		"success": {Type: graphql.Boolean},
		"message": {Type: graphql.String},
	},
})

// ─── INPUT TYPES ──────────────────────────────────────────────────────────────

var insertGoldUserInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "InsertGoldUserInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"gold_email":              {Type: graphql.NewNonNull(graphql.String)},
		"gold_password":           {Type: graphql.NewNonNull(graphql.String)},
		"gold_nama":               {Type: graphql.NewNonNull(graphql.String)},
		"gold_nomorhp":            {Type: graphql.NewNonNull(graphql.String)},
		"gold_nomorkartu":         {Type: graphql.NewNonNull(graphql.String)},
		"gold_cvv":                {Type: graphql.NewNonNull(graphql.String)},
		"gold_expireddate":        {Type: graphql.NewNonNull(graphql.String)},
		"gold_namapemegangkartu":  {Type: graphql.NewNonNull(graphql.String)},
	},
})

var updatePasswordInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UpdatePasswordInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"gold_email":    {Type: graphql.NewNonNull(graphql.String)},
		"gold_password": {Type: graphql.NewNonNull(graphql.String)},
		"gold_otp":      {Type: graphql.NewNonNull(graphql.String)},
	},
})

var updateNamaInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UpdateNamaInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"gold_email": {Type: graphql.NewNonNull(graphql.String)},
		"gold_nama":  {Type: graphql.NewNonNull(graphql.String)},
	},
})

var updateKartuInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UpdateKartuInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"gold_email":      {Type: graphql.NewNonNull(graphql.String)},
		"gold_nomorkartu": {Type: graphql.NewNonNull(graphql.String)},
		"gold_cvv":        {Type: graphql.NewNonNull(graphql.String)},
	},
})

var insertSubscriptionInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "InsertSubscriptionInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"gold_email": {Type: graphql.NewNonNull(graphql.String)},
		"menu_ids":   {Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))},
	},
})

var insertSubsDetailInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "InsertSubsDetailInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"gold_id":     {Type: graphql.NewNonNull(graphql.Int)},
		"gold_menuid": {Type: graphql.NewNonNull(graphql.Int)},
	},
})

var updateSubsDetailInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UpdateSubsDetailInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"gold_id":              {Type: graphql.NewNonNull(graphql.Int)},
		"gold_menuid":          {Type: graphql.NewNonNull(graphql.Int)},
		"gold_jumlahpertemuan": {Type: graphql.NewNonNull(graphql.Int)},
	},
})

// ─── QUERY TYPE ───────────────────────────────────────────────────────────────

func buildQueryType(svc IgoldgymSvc) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"goldUsers": {
				Type:        graphql.NewList(goldUserType),
				Description: "Get semua gold gym users",
				Resolve:     resolveGetGoldUsers(svc),
			},
			"goldUserByEmail": {
				Type:        graphql.String,
				Description: "Cek status registrasi email (TERDAFTAR / TIDAK TERDAFTAR / BELUM TERVALIDASI)",
				Args: graphql.FieldConfigArgument{
					"email": {Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolveGoldUserByEmail(svc),
			},
			"goldUserDataByEmail": {
				Type:        goldUserDetailType,
				Description: "Get data lengkap user berdasarkan email",
				Args: graphql.FieldConfigArgument{
					"email": {Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolveGoldUserDataByEmail(svc),
			},
			"allSubscriptions": {
				Type:        graphql.NewList(subscriptionType),
				Description: "Get semua paket langganan yang tersedia",
				Resolve:     resolveAllSubscriptions(svc),
			},
			"subsWithUser": {
				Type:        graphql.NewList(subsWithUserType),
				Description: "Get langganan beserta informasi user",
				Resolve:     resolveSubsWithUser(svc),
			},
			"subscriptionTotalHarga": {
				Type:        subscriptionPaymentType,
				Description: "Get total harga langganan berdasarkan email user",
				Args: graphql.FieldConfigArgument{
					"email": {Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolveSubsTotalHarga(svc),
			},
		},
	})
}

// ─── MUTATION TYPE ────────────────────────────────────────────────────────────

func buildMutationType(svc IgoldgymSvc) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			// INSERT
			"insertGoldUser": {
				Type:        mutationResultType,
				Description: "Insert user baru ke sistem (otomatis kirim OTP ke email)",
				Args: graphql.FieldConfigArgument{
					"input": {Type: graphql.NewNonNull(insertGoldUserInputType)},
				},
				Resolve: resolveInsertGoldUser(svc),
			},
			"insertSubscription": {
				Type:        mutationResultType,
				Description: "Insert langganan (header + detail sekaligus) berdasarkan email dan menu_ids",
				Args: graphql.FieldConfigArgument{
					"input": {Type: graphql.NewNonNull(insertSubscriptionInputType)},
				},
				Resolve: resolveInsertSubscription(svc),
			},
			"insertSubscriptionDetail": {
				Type:        mutationResultType,
				Description: "Insert satu detail langganan untuk user",
				Args: graphql.FieldConfigArgument{
					"input": {Type: graphql.NewNonNull(insertSubsDetailInputType)},
				},
				Resolve: resolveInsertSubsDetail(svc),
			},

			// UPDATE
			"updatePassword": {
				Type:        mutationResultType,
				Description: "Update password user (wajib validasi OTP terlebih dahulu via sendOTP)",
				Args: graphql.FieldConfigArgument{
					"input": {Type: graphql.NewNonNull(updatePasswordInputType)},
				},
				Resolve: resolveUpdatePassword(svc),
			},
			"updateNama": {
				Type:        mutationResultType,
				Description: "Update nama member",
				Args: graphql.FieldConfigArgument{
					"input": {Type: graphql.NewNonNull(updateNamaInputType)},
				},
				Resolve: resolveUpdateNama(svc),
			},
			"updateKartu": {
				Type:        mutationResultType,
				Description: "Update informasi kartu (nomor kartu dan CVV di-hash argon2 otomatis)",
				Args: graphql.FieldConfigArgument{
					"input": {Type: graphql.NewNonNull(updateKartuInputType)},
				},
				Resolve: resolveUpdateKartu(svc),
			},
			"updateSubscriptionDetail": {
				Type:        mutationResultType,
				Description: "Update detail langganan (jumlah pertemuan, dll)",
				Args: graphql.FieldConfigArgument{
					"input": {Type: graphql.NewNonNull(updateSubsDetailInputType)},
				},
				Resolve: resolveUpdateSubsDetail(svc),
			},

			// OTP
			"validateOTP": {
				Type:        mutationResultType,
				Description: "Validasi OTP untuk aktivasi akun setelah registrasi",
				Args: graphql.FieldConfigArgument{
					"otp":   {Type: graphql.NewNonNull(graphql.String)},
					"email": {Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolveValidateOTP(svc),
			},
			"sendOTP": {
				Type:        mutationResultType,
				Description: "Kirim OTP ke email untuk keperluan reset password atau aktivasi",
				Args: graphql.FieldConfigArgument{
					"email": {Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolveSendOTP(svc),
			},
			"sendOTPSubscription": {
				Type:        mutationResultType,
				Description: "Kirim OTP ke email untuk proses pembayaran langganan",
				Args: graphql.FieldConfigArgument{
					"email": {Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolveSendOTPSubscription(svc),
			},
			"processPayment": {
				Type:        mutationResultType,
				Description: "Proses pembayaran langganan dengan verifikasi OTP (OTP expired setelah 5 menit)",
				Args: graphql.FieldConfigArgument{
					"otp":   {Type: graphql.NewNonNull(graphql.String)},
					"email": {Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolveProcessPayment(svc),
			},

			// AUTH
			"loginUser": {
				Type:        loginResultType,
				Description: "Login user — mengembalikan JWT access token",
				Args: graphql.FieldConfigArgument{
					"email":    {Type: graphql.NewNonNull(graphql.String)},
					"password": {Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolveLoginUser(svc),
			},
			"logout": {
				Type:        mutationResultType,
				Description: "Logout user (hapus token)",
				Args: graphql.FieldConfigArgument{
					"email": {Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: resolveLogout(svc),
			},
		},
	})
}
