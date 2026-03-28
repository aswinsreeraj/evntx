package main

import (
	"log"
	"time"

	"github.com/aswinsreeraj/evntx/internal/infrastructure/database"
	infraRepo "github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	if err := db.AutoMigrate(
		&infraRepo.UserModel{},
		&infraRepo.WalletModel{},
	); err != nil {
		log.Fatal("failed to run wallet seeder migrations:", err)
	}

	var usersWithoutWallet []infraRepo.UserModel
	if err := db.Model(&infraRepo.UserModel{}).
		Joins("LEFT JOIN wallet_models ON wallet_models.user_id = user_models.id").
		Where("wallet_models.id IS NULL").
		Find(&usersWithoutWallet).Error; err != nil {
		log.Fatal("failed to fetch users without wallets:", err)
	}

	if len(usersWithoutWallet) == 0 {
		log.Println("No users without wallets found")
		return
	}

	now := time.Now()
	wallets := make([]infraRepo.WalletModel, 0, len(usersWithoutWallet))
	for _, user := range usersWithoutWallet {
		wallets = append(wallets, infraRepo.WalletModel{
			ID:               uuid.NewString(),
			UserID:           user.ID,
			AvailableBalance: 0,
			PendingBalance:   0,
			TotalCredited:    0,
			TotalDebited:     0,
			UpdatedAt:        now,
		})
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoNothing: true,
		}).Create(&wallets)
		if result.Error != nil {
			return result.Error
		}

		log.Printf("Created %d wallet(s) for existing users", result.RowsAffected)
		return nil
	}); err != nil {
		log.Fatal("failed to create wallets for existing users:", err)
	}
}
