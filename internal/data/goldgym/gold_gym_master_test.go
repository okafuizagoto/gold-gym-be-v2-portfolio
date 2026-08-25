package goldgym

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	goldEntity "gold-gym-be/internal/entity/goldgym"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// setupMockDB creates a mock database for testing
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	dialector := mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	cleanup := func() {
		mockDB.Close()
	}

	return db, mock, cleanup
}

// setupMockSQLX creates a mock sqlx database for testing prepared statements
func setupMockSQLX(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	cleanup := func() {
		mockDB.Close()
	}

	return sqlxDB, mock, cleanup
}

// =============================================================================
// GetGoldUser Tests
// =============================================================================

func TestGetGoldUser_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &Data{db: db}

	// Mock query expectations
	rows := sqlmock.NewRows([]string{
		"gold_id", "gold_email", "gold_password", "gold_nama",
		"gold_nomorhp", "gold_nomorkartu", "gold_cvv",
		"gold_expireddate", "gold_namapemegangkartu",
	}).
		AddRow(1, "test@example.com", "hashedpass", "Test User",
			"08123456789", "1234567890123456", "123",
			"12/25", "TEST USER").
		AddRow(2, "user2@example.com", "hashedpass2", "User Two",
			"08198765432", "6543210987654321", "456",
			"11/26", "USER TWO")

	mock.ExpectQuery("SELECT \\* FROM `data_peserta`").
		WillReturnRows(rows)

	// Execute test
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	users, err := repo.GetGoldUser(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, 1, users[0].GoldId)
	assert.Equal(t, "test@example.com", users[0].GoldEmail)
	assert.Equal(t, "Test User", users[0].GoldNama)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetGoldUser_EmptyResult(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &Data{db: db}

	// Mock empty result
	rows := sqlmock.NewRows([]string{
		"gold_id", "gold_email", "gold_password", "gold_nama",
		"gold_nomorhp", "gold_nomorkartu", "gold_cvv",
		"gold_expireddate", "gold_namapemegangkartu",
	})

	mock.ExpectQuery("SELECT \\* FROM `data_peserta`").
		WillReturnRows(rows)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	users, err := repo.GetGoldUser(ctx)

	assert.NoError(t, err)
	assert.Empty(t, users)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetGoldUser_DatabaseError(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &Data{db: db}

	// Mock database error
	mock.ExpectQuery("SELECT \\* FROM `data_peserta`").
		WillReturnError(errors.New("connection refused"))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	users, err := repo.GetGoldUser(ctx)

	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Contains(t, err.Error(), "connection refused")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// =============================================================================
// GetGoldUserByEmail Tests
// =============================================================================

func TestGetGoldUserByEmail_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &Data{db: db}

	rows := sqlmock.NewRows([]string{
		"gold_id", "gold_email", "gold_password", "gold_nama",
		"gold_nomorhp", "gold_nomorkartu", "gold_cvv",
		"gold_expireddate", "gold_namapemegangkartu", "gold_validasiyn",
	}).
		AddRow(1, "test@example.com", "hashedpass", "Test User",
			"08123456789", "1234567890123456", "123",
			"12/25", "TEST USER", "Y")

	mock.ExpectQuery("SELECT \\* FROM `data_peserta` WHERE gold_email = \\? ORDER BY").
		WithArgs("test@example.com", 1).
		WillReturnRows(rows)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err := repo.GetGoldUserByEmail(ctx, "test@example.com")

	assert.NoError(t, err)
	assert.Equal(t, 1, user.GoldId)
	assert.Equal(t, "test@example.com", user.GoldEmail)
	assert.Equal(t, "Test User", user.GoldNama)
	assert.Equal(t, "Y", user.GoldValidasiYN)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetGoldUserByEmail_NotFound(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &Data{db: db}

	mock.ExpectQuery("SELECT \\* FROM `data_peserta` WHERE gold_email = \\? ORDER BY").
		WithArgs("notfound@example.com", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err := repo.GetGoldUserByEmail(ctx, "notfound@example.com")

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Equal(t, 0, user.GoldId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetGoldUserByEmail_EmptyEmail(t *testing.T) {
	db, _, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &Data{db: db}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err := repo.GetGoldUserByEmail(ctx, "")

	// Should handle empty email gracefully
	assert.Error(t, err)
	assert.Equal(t, 0, user.GoldId)
}

// =============================================================================
// GetGoldUserByID Tests
// =============================================================================

func TestGetGoldUserByID_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &Data{db: db}

	rows := sqlmock.NewRows([]string{
		"gold_id", "gold_email", "gold_password", "gold_nama",
		"gold_nomorhp", "gold_nomorkartu", "gold_cvv",
		"gold_expireddate", "gold_namapemegangkartu",
	}).
		AddRow(1, "test@example.com", "hashedpass", "Test User",
			"08123456789", "1234567890123456", "123",
			"12/25", "TEST USER")

	mock.ExpectQuery("SELECT \\* FROM `data_peserta` WHERE gold_id = \\? ORDER BY").
		WithArgs("1", 1).
		WillReturnRows(rows)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err := repo.GetGoldUserByID(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, 1, user.GoldId)
	assert.Equal(t, "test@example.com", user.GoldEmail)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// =============================================================================
// InsertGoldUser Tests
// =============================================================================

func TestInsertGoldUser_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &Data{db: db}

	user := goldEntity.GetGoldUsers{
		GoldEmail:         "newuser@example.com",
		GoldPassword:      "hashedpassword",
		GoldNama:          "New User",
		GoldNomorHp:       "08123456789",
		GoldNomorKartu:    "1234567890123456",
		GoldCvv:           "123",
		GoldExpireddate:   "12/25",
		GoldPemegangKartu: "NEW USER",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `data_peserta`").
		WithArgs(
			sqlmock.AnyArg(), // gold_id (0)
			user.GoldEmail,
			user.GoldPassword,
			user.GoldNama,
			user.GoldNomorHp,
			user.GoldNomorKartu,
			user.GoldCvv,
			user.GoldExpireddate,
			user.GoldPemegangKartu,
			user.GoldOTP, // "" (not set in test)
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userID, err := repo.InsertGoldUser(ctx, user)

	assert.NoError(t, err)
	assert.NotEmpty(t, userID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertGoldUser_DuplicateEmail(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &Data{db: db}

	user := goldEntity.GetGoldUsers{
		GoldEmail:    "duplicate@example.com",
		GoldPassword: "hashedpass",
		GoldNama:     "Duplicate User",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `data_peserta`").
		WillReturnError(errors.New("Error 1062: Duplicate entry"))
	mock.ExpectRollback()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userID, err := repo.InsertGoldUser(ctx, user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Duplicate entry")
	assert.Equal(t, "Gagal", userID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// =============================================================================
// UpdateDataPeserta Tests
// =============================================================================

func TestUpdateDataPeserta_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &Data{db: db}

	updateData := goldEntity.UpdatePassword{
		GoldEmail:    "test@example.com",
		GoldPassword: "newhashedpass",
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `data_peserta` SET").
		WithArgs(
			updateData.GoldPassword,
			updateData.GoldEmail,
			updateData.GoldOTP,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := repo.UpdateDataPeserta(ctx, updateData)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// =============================================================================
// UpdateNama Tests
// =============================================================================

func TestUpdateNama_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &Data{db: db}

	updateData := goldEntity.UpdateNama{
		GoldEmail: "test@example.com",
		GoldNama:  "Updated Name",
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `data_peserta` SET").
		WithArgs(
			updateData.GoldNama,
			updateData.GoldEmail,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := repo.UpdateNama(ctx, updateData)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// =============================================================================
// Logout Tests
// =============================================================================

func TestLogout_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &Data{db: db}

	logoutData := goldEntity.Logout{
		GoldEmail: "test@example.com",
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `data_peserta` SET").
		WithArgs(
			sqlmock.AnyArg(), // NULL value
			logoutData.GoldEmail,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := repo.Logout(ctx, logoutData)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// =============================================================================
// GetAllSubscription Tests
// =============================================================================

func TestGetAllSubscription_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &Data{db: db}

	rows := sqlmock.NewRows([]string{
		"gold_namapaket", "gold_namalayanan", "gold_harga", "gold_jadwal",
		"gold_listlatihan", "gold_jumlahpertemuan", "gold_durasi",
	}).
		AddRow("Monthly", "Personal Training", 500000.0, "Mon-Fri",
			"Cardio, Weight", 12, 30).
		AddRow("Yearly", "Group Class", 5000000.0, "Mon-Sun",
			"Full Access", 365, 365)

	mock.ExpectQuery("SELECT \\* FROM `subscription_product`").
		WillReturnRows(rows)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	subscriptions, err := repo.GetAllSubscription(ctx)

	assert.NoError(t, err)
	assert.Len(t, subscriptions, 2)
	assert.Equal(t, "Monthly", subscriptions[0].GoldNamaPaket)
	assert.Equal(t, "Personal Training", subscriptions[0].GoldNamaLayanan)
	assert.Equal(t, 500000.0, subscriptions[0].GoldHarga)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// =============================================================================
// Context Timeout Tests
// =============================================================================

func TestGetGoldUser_ContextTimeout(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &Data{db: db}

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Wait for context to timeout
	time.Sleep(2 * time.Millisecond)

	mock.ExpectQuery("SELECT \\* FROM `data_peserta`").
		WillDelayFor(100 * time.Millisecond)

	users, err := repo.GetGoldUser(ctx)

	// Should return context deadline exceeded error
	assert.Error(t, err)
	assert.Nil(t, users)
}

// =============================================================================
// Table-Driven Tests
// =============================================================================

func TestGetGoldUserByEmail_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		mockRows  *sqlmock.Rows
		mockError error
		wantErr   bool
		wantID    int
	}{
		{
			name:  "success - user found",
			email: "found@example.com",
			mockRows: sqlmock.NewRows([]string{
				"gold_id", "gold_email", "gold_nama",
			}).AddRow(1, "found@example.com", "Found User"),
			wantErr: false,
			wantID:  1,
		},
		{
			name:      "error - user not found",
			email:     "notfound@example.com",
			mockError: gorm.ErrRecordNotFound,
			wantErr:   true,
			wantID:    0,
		},
		{
			name:      "error - database error",
			email:     "error@example.com",
			mockError: sql.ErrConnDone,
			wantErr:   true,
			wantID:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			repo := &Data{db: db}

			expectation := mock.ExpectQuery("SELECT \\* FROM `data_peserta` WHERE gold_email = \\? ORDER BY").
				WithArgs(tt.email, 1)

			if tt.mockError != nil {
				expectation.WillReturnError(tt.mockError)
			} else {
				expectation.WillReturnRows(tt.mockRows)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			user, err := repo.GetGoldUserByEmail(ctx, tt.email)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, user.GoldId)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
