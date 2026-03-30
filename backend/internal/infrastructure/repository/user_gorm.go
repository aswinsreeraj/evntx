package repository

import (
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID            string `gorm:"type:uuid;primaryKey"`
	Name          string
	Email         string `gorm:"uniqueIndex"`
	Mobile        string
	Dob           string
	Gender        string
	ProfileImage  string
	Locations     []string `gorm:"serializer:json"`
	IsActive      bool
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type OrganizerDetailModel struct {
	UserID           string `gorm:"type:uuid;primaryKey"`
	OrganizationName string
	Address          string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type userGormRepository struct {
	db *gorm.DB
}

func NewUserGormRepository(db *gorm.DB) *userGormRepository {
	return &userGormRepository{db: db}
}

func (r *userGormRepository) Create(user *domain.User) error {
	model := UserModel{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		Mobile:        user.Mobile,
		Dob:           user.Dob,
		Gender:        user.Gender,
		ProfileImage:  user.ProfileImage,
		Locations:     user.Locations,
		IsActive:      user.IsActive,
		EmailVerified: user.EmailVerified,
	}

	walletModel := WalletModel{
		ID:               uuid.NewString(),
		UserID:           user.ID,
		AvailableBalance: 0,
		PendingBalance:   0,
		TotalCredited:    0,
		TotalDebited:     0,
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&model).Error; err != nil {
			return err
		}

		return tx.Create(&walletModel).Error
	})
}

func (r *userGormRepository) FindByEmail(email string) (*domain.User, error) {
	var model UserModel

	err := r.db.Where("email = ?", email).First(&model).Error
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:            model.ID,
		Name:          model.Name,
		Email:         model.Email,
		Mobile:        model.Mobile,
		Dob:           model.Dob,
		Gender:        model.Gender,
		ProfileImage:  model.ProfileImage,
		Locations:     model.Locations,
		IsActive:      model.IsActive,
		EmailVerified: model.EmailVerified,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}, nil
}

func (r *userGormRepository) FindByID(id string) (*domain.User, error) {
	var model UserModel

	err := r.db.First(&model, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:            model.ID,
		Name:          model.Name,
		Email:         model.Email,
		Mobile:        model.Mobile,
		Dob:           model.Dob,
		Gender:        model.Gender,
		ProfileImage:  model.ProfileImage,
		Locations:     model.Locations,
		IsActive:      model.IsActive,
		EmailVerified: model.EmailVerified,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}, nil
}

func (r *userGormRepository) Update(user *domain.User) error {
	return r.db.Model(&UserModel{}).
		Where("id = ?", user.ID).
		Select("name", "email", "mobile", "dob", "gender", "profile_image", "locations", "is_active", "email_verified", "updated_at").
		Updates(UserModel{
			Name:          user.Name,
			Email:         user.Email,
			Mobile:        user.Mobile,
			Dob:           user.Dob,
			Gender:        user.Gender,
			ProfileImage:  user.ProfileImage,
			Locations:     user.Locations,
			IsActive:      user.IsActive,
			EmailVerified: user.EmailVerified,
			UpdatedAt:     time.Now(),
		}).Error
}

func (r *userGormRepository) GetOrganizerDetails(userID string) (*domain.OrganizerDetail, error) {
	var model OrganizerDetailModel

	if err := r.db.Where("user_id = ?", userID).First(&model).Error; err != nil {
		return nil, err
	}

	return &domain.OrganizerDetail{
		UserID:           model.UserID,
		OrganizationName: model.OrganizationName,
		Address:          model.Address,
	}, nil
}

func (r *userGormRepository) UpsertOrganizerDetails(detail *domain.OrganizerDetail) error {
	model := OrganizerDetailModel{
		UserID:           detail.UserID,
		OrganizationName: detail.OrganizationName,
		Address:          detail.Address,
		UpdatedAt:        time.Now(),
	}

	return r.db.Where("user_id = ?", detail.UserID).Assign(model).FirstOrCreate(&model).Error
}

func (r *userGormRepository) Search(
	search string,
	status string,
	page int,
	limit int,
) ([]domain.AdminUserDetails, int64, error) {

	var results []struct {
		UserModel
		TotalBookings int64
		WalletBalance float64
	}
	var total int64

	query := r.db.Table("user_models").
		Select(`
			user_models.*,
			COALESCE((
				SELECT COUNT(b.id) FROM booking_models b
				WHERE b.user_id = user_models.id::text AND b.status IN ('paid', 'confirmed')
			), 0) AS total_bookings,
			COALESCE((
				SELECT available_balance FROM wallet_models w
				WHERE w.user_id = user_models.id
			), 0) AS wallet_balance
		`)

	query = query.Where(
		"NOT EXISTS (SELECT 1 FROM user_role_models WHERE user_role_models.user_id::uuid = user_models.id)",
	)

	if search != "" {
		query = query.Where(
			"user_models.name ILIKE ? OR user_models.email ILIKE ?",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	if status == "active" {
		query = query.Where("user_models.is_active = ?", true)
	} else if status == "suspended" || status == "inactive" {
		query = query.Where("user_models.is_active = ?", false)
	}

	query.Count(&total)

	offset := (page - 1) * limit

	err := query.
		Order("user_models.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&results).Error

	if err != nil {
		return nil, 0, err
	}

	users := make([]domain.AdminUserDetails, 0, len(results))
	for _, res := range results {
		users = append(users, domain.AdminUserDetails{
			User: domain.User{
				ID:            res.ID,
				Name:          res.Name,
				Email:         res.Email,
				Mobile:        res.Mobile,
				Dob:           res.Dob,
				Gender:        res.Gender,
				ProfileImage:  res.ProfileImage,
				Locations:     res.Locations,
				IsActive:      res.IsActive,
				EmailVerified: res.EmailVerified,
				CreatedAt:     res.CreatedAt,
				UpdatedAt:     res.UpdatedAt,
			},
			TotalBookings: res.TotalBookings,
			WalletBalance: res.WalletBalance,
		})
	}

	return users, total, nil
}

func (r *userGormRepository) SearchOrganizers(
	search string,
	status string,
	page int,
	limit int,
) ([]domain.OrganizerDetails, int64, error) {

	var models []struct {
		UserModel
		OrganizationName string
		Address          string
		TotalEvents      int64
		TotalBookings    int64
		TotalRevenue     float64
		WalletBalance    float64
	}
	var total int64

	query := r.db.Table("user_models").
		Select(`
			user_models.*,
			organizer_detail_models.organization_name,
			organizer_detail_models.address,
			COALESCE((
				SELECT COUNT(e.id) FROM event_models e WHERE e.organizer_id = user_models.id::text
			), 0) AS total_events,
			COALESCE((
				SELECT COUNT(b.id) FROM booking_models b
				JOIN event_models e ON e.id = b.event_id
				WHERE e.organizer_id = user_models.id::text AND b.status IN ('paid', 'confirmed')
			), 0) AS total_bookings,
			COALESCE((
				SELECT SUM(b.total_amount) FROM booking_models b
				JOIN event_models e ON e.id = b.event_id
				WHERE e.organizer_id = user_models.id::text AND b.status IN ('paid', 'confirmed')
			), 0) AS total_revenue,
			COALESCE((
				SELECT available_balance FROM wallet_models w
				WHERE w.user_id = user_models.id
			), 0) AS wallet_balance
		`).
		Joins("INNER JOIN user_role_models ON user_role_models.user_id::uuid = user_models.id AND user_role_models.role = ?", domain.RoleOrganizer).
		Joins("LEFT JOIN organizer_detail_models ON organizer_detail_models.user_id::uuid = user_models.id")

	if search != "" {
		query = query.Where(
			"user_models.name ILIKE ? OR user_models.email ILIKE ?",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	if status == "active" {
		query = query.Where("user_models.is_active = ?", true)
	} else if status == "suspended" || status == "inactive" {
		query = query.Where("user_models.is_active = ?", false)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.
		Order("user_models.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error

	if err != nil {
		return nil, 0, err
	}

	orgs := make([]domain.OrganizerDetails, 0, len(models))
	for _, m := range models {
		orgs = append(orgs, domain.OrganizerDetails{
			User: domain.User{
				ID:        m.ID,
				Name:      m.Name,
				Email:     m.Email,
				IsActive:  m.IsActive,
				CreatedAt: m.CreatedAt,
			},
			OrganizerDetail: domain.OrganizerDetail{
				OrganizationName: m.OrganizationName,
				Address:          m.Address,
			},
			TotalBookings: m.TotalBookings,
			TotalEvents:   m.TotalEvents,
			WalletBalance: m.WalletBalance,
			TotalRevenue:  m.TotalRevenue,
		})
	}

	return orgs, total, nil
}

func (r *userGormRepository) UpdateStatus(userID string, isActive bool) error {
	return r.db.Model(&UserModel{}).
		Where("id = ?", userID).
		Update("is_active", isActive).Error
}
