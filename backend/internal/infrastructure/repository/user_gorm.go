package repository

import (
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"gorm.io/gorm"
)

type UserModel struct {
	ID               string   `gorm:"type:uuid;primaryKey"`
	Name             string
	Email            string   `gorm:"uniqueIndex"`
	Mobile           string
	Dob              string
	Gender           string
	ProfileImage     string
	OrganizationName string
	Locations        []string `gorm:"serializer:json"`
	IsActive         bool
	EmailVerified    bool
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
		ID:               user.ID,
		Name:             user.Name,
		Email:            user.Email,
		Mobile:           user.Mobile,
		Dob:              user.Dob,
		Gender:           user.Gender,
		ProfileImage:     user.ProfileImage,
		OrganizationName: user.OrganizationName,
		Locations:        user.Locations,
		IsActive:         user.IsActive,
		EmailVerified:    user.EmailVerified,
	}

	return r.db.Create(&model).Error
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
		Gender:           model.Gender,
		ProfileImage:     model.ProfileImage,
		OrganizationName: model.OrganizationName,
		Locations:        model.Locations,
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
		Gender:           model.Gender,
		ProfileImage:     model.ProfileImage,
		OrganizationName: model.OrganizationName,
		Locations:        model.Locations,
		IsActive:      model.IsActive,
		EmailVerified: model.EmailVerified,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}, nil
}

func (r *userGormRepository) Update(user *domain.User) error {
	return r.db.Model(&UserModel{}).
		Where("id = ?", user.ID).
		Select("name", "email", "mobile", "dob", "gender", "profile_image", "organization_name", "locations", "is_active", "email_verified", "updated_at").
		Updates(UserModel{
			Name:             user.Name,
			Email:            user.Email,
			Mobile:           user.Mobile,
			Dob:              user.Dob,
			Gender:           user.Gender,
			ProfileImage:     user.ProfileImage,
			OrganizationName: user.OrganizationName,
			Locations:        user.Locations,
			IsActive:         user.IsActive,
			EmailVerified:    user.EmailVerified,
			UpdatedAt:        time.Now(),
		}).Error
}

func (r *userGormRepository) Search(
	search string,
	status string,
	page int,
	limit int,
) ([]domain.User, int64, error) {

	var models []UserModel
	var total int64

	query := r.db.Model(&UserModel{})

	
	query = query.Where(
		"NOT EXISTS (SELECT 1 FROM user_role_models WHERE user_role_models.user_id::uuid = user_models.id AND user_role_models.role = ?)",
		domain.RoleAdmin,
	)

	if search != "" {
		query = query.Where(
			"name ILIKE ? OR email ILIKE ?",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	if status == "active" {
		query = query.Where("is_active = ?", true)
	} else if status == "suspended" || status == "inactive" {
		query = query.Where("is_active = ?", false)
	}

	query.Count(&total)

	offset := (page - 1) * limit

	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error

	if err != nil {
		return nil, 0, err
	}

	users := make([]domain.User, 0)
	for _, m := range models {
		users = append(users, domain.User{
			ID:               m.ID,
			Name:             m.Name,
			Email:            m.Email,
			Mobile:           m.Mobile,
			Dob:              m.Dob,
			Gender:           m.Gender,
			ProfileImage:     m.ProfileImage,
			OrganizationName: m.OrganizationName,
			Locations:        m.Locations,
			IsActive:         m.IsActive,
			EmailVerified:    m.EmailVerified,
			CreatedAt:        m.CreatedAt,
			UpdatedAt:        m.UpdatedAt,
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
		TotalEvents int64
	}
	var total int64

	query := r.db.Table("user_models").
		Select("user_models.*, COALESCE((SELECT COUNT(id) FROM event_models WHERE event_models.organizer_id = user_models.id), 0) as total_events").
		Joins("INNER JOIN user_role_models ON user_role_models.user_id::uuid = user_models.id AND user_role_models.role = ?", domain.RoleOrganizer)

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
				ID:               m.ID,
				Name:             m.Name,
				Email:            m.Email,
				Mobile:           m.Mobile,
				Dob:              m.Dob,
				Gender:           m.Gender,
				ProfileImage:     m.ProfileImage,
				OrganizationName: m.OrganizationName,
				Locations:        m.Locations,
				IsActive:         m.IsActive,
				EmailVerified:    m.EmailVerified,
				CreatedAt:        m.CreatedAt,
				UpdatedAt:        m.UpdatedAt,
			},
			TotalBookings: 0,
			TotalEvents:   m.TotalEvents,
			WalletBalance: 0,
			TotalRevenue:  0,
		})
	}

	return orgs, total, nil
}

func (r *userGormRepository) UpdateStatus(userID string, isActive bool) error {
	return r.db.Model(&UserModel{}).
		Where("id = ?", userID).
		Update("is_active", isActive).Error
}
