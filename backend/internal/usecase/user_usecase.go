package usecase

import (
	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserUsecase struct {
	repo     repository.UserRepository
	roleRepo repository.UserRoleRepository
}

func NewUserUsecase(r repository.UserRepository, roleRepo repository.UserRoleRepository) *UserUsecase {
	return &UserUsecase{repo: r, roleRepo: roleRepo}
}

func (u *UserUsecase) Register(email string) (*domain.User, error) {
	user := &domain.User{
		ID:            uuid.NewString(),
		Email:         email,
		IsActive:      true,
		EmailVerified: false,
	}

	err := u.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) GetUserProfile(userID string) (*domain.User, error) {
	return u.repo.FindByID(userID)
}

func (u *UserUsecase) GetProfile(userID string) (*domain.User, *domain.OrganizerDetail, []domain.UserRole, error) {
	user, err := u.repo.FindByID(userID)
	if err != nil {
		return nil, nil, nil, err
	}

	roles, err := u.roleRepo.GetRolesByUserID(userID)
	if err == nil {
		for _, role := range roles {
			if role == domain.RoleOrganizer {
				detail, detailErr := u.repo.GetOrganizerDetails(userID)
				if detailErr == nil {
					return user, detail, roles, nil
				} else if detailErr != gorm.ErrRecordNotFound {
					return nil, nil, nil, detailErr
				}
				break
			}
		}
		return user, nil, roles, nil
	}

	return user, nil, []domain.UserRole{}, nil
}

func (u *UserUsecase) UpdateProfile(userID, name, mobile, dob, gender, organizationName, address string, locations []string) error {
	user, err := u.repo.FindByID(userID)
	if err != nil {
		return err
	}

	if name != "" {
		user.Name = name
	}
	user.Mobile = mobile
	user.Dob = dob
	user.Gender = gender
	user.Locations = locations

	if err := u.repo.Update(user); err != nil {
		return err
	}

	roles, err := u.roleRepo.GetRolesByUserID(userID)
	if err != nil {
		return err
	}

	for _, role := range roles {
		if role == domain.RoleOrganizer {
			return u.repo.UpsertOrganizerDetails(&domain.OrganizerDetail{
				UserID:           userID,
				OrganizationName: organizationName,
				Address:          address,
			})
		}
	}

	return nil
}

func (u *UserUsecase) AdminSearchUsers(
	search string,
	status string,
	page int,
	limit int,
) ([]domain.User, int64, error) {
	return u.repo.Search(search, status, page, limit)
}

func (u *UserUsecase) AdminSearchOrganizers(
	search string,
	status string,
	page int,
	limit int,
) ([]domain.OrganizerDetails, int64, error) {
	return u.repo.SearchOrganizers(search, status, page, limit)
}

func (u *UserUsecase) AdminUpdateUserStatus(userID string, isActive bool) error {
	return u.repo.UpdateStatus(userID, isActive)
}

func (u *UserUsecase) UploadProfileImage(userID string, imageURL string) error {
	user, err := u.repo.FindByID(userID)
	if err != nil {
		return err
	}

	user.ProfileImage = imageURL
	return u.repo.Update(user)
}
