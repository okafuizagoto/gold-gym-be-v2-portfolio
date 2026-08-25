package goldgym

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/storage"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	jaegerLog "gold-gym-be/pkg/log"
)

type (
	// Data ...
	Data struct {
		db   *gorm.DB
		dbr  *sqlx.DB
		fsdb *firestore.Client
		s    *storage.Client
		rdb  *redis.Client
		stmt *map[string]*sqlx.Stmt

		tracer opentracing.Tracer
		logger jaegerLog.Factory
	}

	// statement ...
	statement struct {
		key   string
		query string
	}
)

// Tambahkan query di dalam const
// getAllUser = "GetAllUser"
// qGetAllUser = "SELECT * FROM users"
const (
	getOneStockProduct  = "GetOneStockProduct"
	qGetOneStockProduct = `SELECT stock_id, stock_name,stock_qty, stock_last_update, stock_update_by FROM stock WHERE stock_id = ? AND (? = "" OR stock_name LIKE ?) and (? = "" OR stock_id = ?)`

	insertStockSales  = "InsertStockSales"
	qInsertStockSales = `INSERT INTO stock (stock_id, stock_name, stock_qty, stock_last_update, stock_update_by) VALUES (?,?,?,NOW(),?)`

	getLastStock  = "GetLastStock"
	qGetLastStock = `SELECT stock_id, stock_name, stock_qty, stock_last_update, stock_update_by FROM stock ORDER BY stock_last_update DESC LIMIT 1`

	addStockByDate  = "AddStockByDate"
	qAddStockByDate = `INSERT INTO td_stock (stock_id, stock_name, stock_qty, stock_last_update, stock_update_by) VALUES (?,?,?,NOW(),?)`

	getStockByID  = "GetStockByID"
	qGetStockByID = `SELECT stock_id, stock_name, stock_qty, stock_last_update, stock_update_by FROM stock WHERE stock_id = ?`

	updateStockQty  = "UpdateStockQty"
	qUpdateStockQty = `UPDATE stock SET stock_qty = ?, stock_qty_update = NOW() WHERE stock_id = ?`

	getAllStockHeader  = "GetAllStockHeader"
	qGetAllStockHeader = `SELECT stock_id, stock_name, stock_qty, stock_last_update, stock_update_by
FROM stock order by stock_id asc`

	// -------------------------------------------------------------------------------------------------------------------------------

	insertStockNew  = "InsertStockNew"
	qInsertStockNew = `INSERT INTO stock
(stock_id, stock_gold_id, stock_outcode, stock_item_id, stock_name, stock_qty, stock_created_at, stock_last_update, stock_qty_update, stock_update_by)
VALUES %s`

	insertStockHistory  = "InsertStockHistory"
	qInsertStockHistory = `INSERT INTO u868654674_gold_gym_bez.stock_history
(sh_id, sh_stock_id, sh_gold_id, sh_outcode, sh_item_id, sh_status, sh_qty_before, sh_qty_change, sh_qty_after, sh_note, sh_created_at, sh_created_by)
VALUES %s`

// // getJadwal  = "GetJadwal"
// // qGetJadwal = "SELECT * FROM m_jadwal"

// getGoldUser  = "GetGoldUser"
// qGetGoldUser = `SELECT gold_email,gold_password,gold_nama,gold_nomorhp,gold_nomorkartu,gold_cvv,gold_expireddate,gold_namapemegangkartu FROM data_peserta`

// getPasswordByUser  = "GetPasswordByUser"
// qGetPasswordByUser = `select gold_password from data_peserta where gold_email = ?`

// updateLastLogin  = "UpdateLastLogin"
// qUpdateLastLogin = `UPDATE data_peserta SET gold_last_login = NOW(), gold_last_login_host = ? WHERE gold_email = ?`

// insertGoldUser  = "InsertGoldUser"
// qInsertGoldUser = `INSERT INTO data_peserta (gold_id, gold_email,gold_password,gold_nama,gold_nomorhp,gold_nomorkartu,gold_cvv,gold_expireddate,gold_namapemegangkartu, gold_otp) VALUES (?, ?,?,?,?,?,?,?,?,?)`

// getGoldUserByEmail  = "GetGoldUserByEmail"
// qGetGoldUserByEmail = `SELECT gold_id, gold_email,gold_password,gold_nama,gold_nomorhp,gold_nomorkartu,gold_cvv,gold_expireddate,gold_namapemegangkartu,gold_validasiyn, gold_token, gold_otp, gold_updated_by, gold_updated_at, gold_last_login, gold_last_login_host, gold_force_change_password FROM data_peserta WHERE gold_email = ?`

// getGoldUserByEmailLogin  = "GetGoldUserByEmailLogin"
// qGetGoldUserByEmailLogin = `SELECT gold_nama,gold_nomorhp,gold_nomorkartu,gold_cvv,gold_expireddate,gold_namapemegangkartu FROM data_peserta WHERE gold_email = ? and gold_password = ?`

// getGoldToken  = "GetGoldToken"
// qGetGoldToken = `SELECT gold_token FROM data_token`

// updateGoldToken  = "UpdateGoldToken"
// qUpdateGoldToken = `UPDATE data_peserta SET gold_token = ? WHERE gold_email = ?`

// getAllSubscription  = "GetAllSubscription"
// qGetAllSubpsription = `SELECT gold_namapaket, gold_namalayanan, gold_harga, gold_jadwal, gold_listlatihan, gold_jumlahpertemuan, gold_durasi FROM subscription_product`

// getOneSubscription  = "GetOneSubscription"
// qGetOneSubscription = `SELECT gold_namapaket, gold_namalayanan, gold_harga, gold_jadwal, gold_listlatihan, gold_jumlahpertemuan, gold_durasi FROM subscription_product where gold_menuid = ?`

// insertSubscription = "InsertSubscription"
// // qInsertSubscription = `INSERT INTO subscription (gold_id, gold_menuid, gold_namapaket, gold_namalayanan, gold_harga) VALUES (?, ?, ?, ?, ?)`
// qInsertSubscription = `INSERT INTO subscription (gold_id, gold_totalharga, gold_lastupdate) VALUES (?,?,NOW())`

// insertSubscriptionDetail = "InsertSubscriptionDetail"
// // qInsertSubscriptionDetail = `INSERT INTO subscription_detail (gold_id, gold_menuid, gold_jadwal, gold_listlatihan, gold_jumlahpertemuan, gold_durasi, gold_statuslangganan) VALUES (?, ?, ?, ?, ?, ?, ?)`
// qInsertSubscriptionDetail = `INSERT INTO subscription_detail (gold_id, gold_menuid, gold_namapaket, gold_namalayanan, gold_harga, gold_jadwal, gold_listlatihan, gold_jumlahpertemuan, gold_durasi, gold_statuslangganan) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

// // deleteSubscriptionHeader = "DeleteSubscriptionHeader"
// // // qDeleteSubscriptionHeader = `DELETE FROM subscription WHERE gold_id = ? AND gold_menuid = ?`
// // qDeleteSubscriptionHeader = `DELETE FROM subscription WHERE gold_id = ?`

// deleteSubscriptionDetail  = "DeleteSubscriptionDetail"
// qDeleteSubscriptionDetail = `DELETE FROM subscription_detail WHERE gold_id = ? AND gold_menuid = ?`

// updateSubscriptionDetail  = "UpdateSubscriptionDetail"
// qUpdateSubscriptionDetail = `UPDATE subscription_detail SET gold_jumlahpertemuan = ? WHERE gold_id = ? AND gold_menuid = ?`

// updateDataPeserta  = "UpdateDataPeserta"
// qUpdateDataPeserta = `UPDATE data_peserta SET gold_password = ? WHERE gold_email = ? and gold_otp = ?`

// updateNama  = "UpdateNama"
// qUpdateNama = `UPDATE data_peserta SET gold_nama = ? WHERE gold_email = ?`

// updateKartu  = "UpdateKartu"
// qUpdateKartu = `UPDATE data_peserta SET gold_nomorkartu = ?, gold_cvv = ? WHERE gold_email = ?`

// logout  = "Logout"
// qLogout = `UPDATE data_peserta SET gold_token = NULL WHERE gold_email = ?`

// getSubsWithUser = "GetSubsWithUser"
// // qGetSubsWithUser = `SELECT a.gold_id, b.gold_menuid, a.gold_email, a.gold_nama, a.gold_nomorhp, a.gold_expireddate,
// // b.gold_namapaket, b.gold_namalayanan, b.gold_harga, c.gold_listlatihan, c.gold_jumlahpertemuan, c.gold_durasi, c.gold_statuslangganan
// // FROM data_peserta a
// // LEFT JOIN subscription b
// // ON a.gold_id = b.gold_id
// // LEFT JOIN subscription_detail c
// // ON b.gold_id = c.gold_id AND b.gold_menuid = c.gold_menuid`

// qGetSubsWithUser = `SELECT a.gold_id, c.gold_menuid, a.gold_email, a.gold_nama, a.gold_nomorhp, a.gold_expireddate,
// c.gold_namapaket, c.gold_namalayanan, c.gold_harga, c.gold_listlatihan, c.gold_jumlahpertemuan, c.gold_durasi, c.gold_statuslangganan
// FROM data_peserta a
// LEFT JOIN subscription b
// ON a.gold_id = b.gold_id
// LEFT JOIN subscription_detail c
// ON b.gold_id = c.gold_id
// ORDER BY gold_id`

// getValidationGoldOTP  = "GetValidationGoldOTP"
// qGetValidationGoldOTP = `SELECT gold_otp FROM data_peserta WHERE gold_otp = ?`

// updateValidationOTP  = "UpdateValidationOTP"
// qUpdateValidationOTP = `UPDATE data_peserta SET gold_validasiyn = "Y" WHERE gold_email = ? AND gold_otp IS NOT NULL `

// updateOtpIsNull  = "UpdateOtpIsNull"
// qUpdateOtpIsNull = `UPDATE data_peserta SET gold_otp = NULL WHERE gold_email = ?`

// updateOTP  = "UpdateOTP"
// qUpdateOTP = `UPDATE data_peserta SET gold_otp = ? WHERE gold_email = ?`

// updateOTPSubscription  = "UpdateOTPSubscription"
// qUpdateOTPSubscription = `UPDATE subscription SET gold_otp = ?, gold_lastupdate = NOW() WHERE gold_id = ?`

// // getGoldUserByEmail = "GetGoldUserByEmail"
// // qGetGoldUserByID = `SELECT gold_email,gold_password,gold_nama,gold_nomorhp,gold_nomorkartu,gold_cvv,gold_expireddate,gold_namapemegangkartu FROM data_peserta where gold`

// getSubscriptionHeader  = "GetSubscriptionHeader"
// qGetSubscriptionHeader = `SELECT gold_id, gold_totalharga, gold_validasipayment, gold_otp, gold_lastupdate FROM subscription WHERE gold_id = ?`

// updateValidasiPaymentHeader  = "UpdateValidasiPaymentHeader"
// qUpdateValidasiPaymentHeader = `UPDATE subscription SET gold_validasipayment = "Y", gold_lastupdate = NOW() WHERE gold_id = ?`

// updateValidasiPaymentDetail  = "UpdateValidasiPaymentDetail"
// qUpdateValidasiPaymentDetail = `UPDATE subscription_detail SET gold_startdate = NOW(), gold_enddate = DATE_ADD(NOW(),INTERVAL 30 DAY), gold_statuslangganan = "Berlangganan" WHERE gold_id = ?`

// getSubscriptionHeaderTotalHarga  = "GetSubscriptionHeaderTotalHarga"
// qGetSubscriptionHeaderTotalHarga = `SELECT gold_totalharga, gold_validasipayment FROM subscription WHERE gold_id = ?`
)

var (
	readStmt = []statement{
		{getOneStockProduct, qGetOneStockProduct},
		{getLastStock, qGetLastStock},
		{getStockByID, qGetStockByID},
		{getAllStockHeader, qGetAllStockHeader},
		// {getGoldUser, qGetGoldUser},
		// {getGoldUserByEmail, qGetGoldUserByEmail},
		// {getGoldUserByEmailLogin, qGetGoldUserByEmailLogin},
		// {getGoldToken, qGetGoldToken},
		// {getAllSubscription, qGetAllSubpsription},
		// {getSubsWithUser, qGetSubsWithUser},
		// {getValidationGoldOTP, qGetValidationGoldOTP},
		// {getOneSubscription, qGetOneSubscription},
		// {getSubscriptionHeader, qGetSubscriptionHeader},
		// {getSubscriptionHeaderTotalHarga, qGetSubscriptionHeaderTotalHarga},
		// {getPasswordByUser, qGetPasswordByUser},
	}
	insertStmt = []statement{
		{insertStockSales, qInsertStockSales},
		{addStockByDate, qAddStockByDate},
		// {insertGoldUser, qInsertGoldUser},
		// {insertSubscription, qInsertSubscription},
		// {insertSubscriptionDetail, qInsertSubscriptionDetail},
	}
	updateStmt = []statement{
		{updateStockQty, qUpdateStockQty},
		// {updateGoldToken, qUpdateGoldToken},
		// {updateSubscriptionDetail, qUpdateSubscriptionDetail},
		// {updateDataPeserta, qUpdateDataPeserta},
		// {updateNama, qUpdateNama},
		// {updateKartu, qUpdateKartu},
		// {logout, qLogout},
		// {updateValidationOTP, qUpdateValidationOTP},
		// {updateOtpIsNull, qUpdateOtpIsNull},
		// {updateOTP, qUpdateOTP},
		// {updateOTPSubscription, qUpdateOTPSubscription},
		// {updateValidasiPaymentHeader, qUpdateValidasiPaymentHeader},
		// {updateValidasiPaymentDetail, qUpdateValidasiPaymentDetail},
		// {updateLastLogin, qUpdateLastLogin},
	}
	deleteStmt = []statement{
		// {deleteSubscriptionHeader, qDeleteSubscriptionHeader},
		// {deleteSubscriptionDetail, qDeleteSubscriptionDetail},
	}
)

// New ...
// func New(db *gorm.DB, fsdb *db.Client, fs *storage.Client, rdb *redis.Client, tracer opentracing.Tracer, logger jaegerLog.Factory) *Data {
func New(db *gorm.DB, dbr *sqlx.DB, fsdb *firestore.Client, fs *storage.Client, rdb *redis.Client, tracer opentracing.Tracer, logger jaegerLog.Factory) *Data {
	var (
		stmts = make(map[string]*sqlx.Stmt)
	)
	d := &Data{
		db:     db,
		dbr:    dbr,
		fsdb:   fsdb,
		s:      fs,
		rdb:    rdb,
		tracer: tracer,
		logger: logger,
		stmt:   &stmts,
	}

	d.InitStmt()
	return d
}

func (d *Data) InitStmt() {
	var (
		err   error
		stmts = make(map[string]*sqlx.Stmt)
	)

	for _, v := range readStmt {
		stmts[v.key], err = d.dbr.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize select statement key %v, err : %v", v.key, err)
		}
	}

	for _, v := range insertStmt {
		stmts[v.key], err = d.dbr.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize insert statement key %v, err : %v", v.key, err)
		}
	}

	for _, v := range updateStmt {
		stmts[v.key], err = d.dbr.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize update statement key %v, err : %v", v.key, err)
		}
	}

	for _, v := range deleteStmt {
		stmts[v.key], err = d.dbr.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize delete statement key %v, err : %v", v.key, err)
		}
	}

	*d.stmt = stmts
}

// Close will cleanup prepared statements
func (d *Data) Close() {
	if d.stmt == nil {
		return
	}

	for k, stmt := range *d.stmt {
		if stmt != nil {
			if err := stmt.Close(); err != nil {
				log.Printf("[DB] failed to close stmt %s: %v", k, err)
			}
		}
	}
}
