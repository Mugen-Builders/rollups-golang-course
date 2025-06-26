package sqlite

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	. "github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/custom_type"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/domain/entity"
)

type SQLiteRepository struct {
	Db *gorm.DB
}

func (r *SQLiteRepository) Close() error {
	sqlDB, err := r.Db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	return sqlDB.Close()
}

func NewSQLiteRepository(conn string) (*SQLiteRepository, error) {
	dbPath := strings.TrimPrefix(conn, "sqlite://")

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&entity.Campaign{},
		&entity.Order{},
		&entity.User{},
	)
	if err != nil {
		return nil, err
	}

	var adminAddress string
	if dbPath == ":memory:" {
		adminAddress = "0x976EA74026E726554dB657fA54763abd0C3a0aa9"
	} else {
		adminAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266" // TODO: change to admin address
	}

	adminUser := entity.User{
		Role:      entity.UserRoleAdmin,
		Address:   HexToAddress(adminAddress),
		CreatedAt: time.Now().Unix(),
	}

	if err := db.Create(&adminUser).Error; err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	return &SQLiteRepository{Db: db}, nil
}
