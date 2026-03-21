package usecase

import (
	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	"github.com/google/uuid"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: r}
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

func (u *UserUsecase) GetProfile(userID string) (*domain.User, error) {
	return u.repo.FindByID(userID)
}

func (u *UserUsecase) UpdateProfile(userID, name, mobile, dob, gender, organizationName string, locations []string) error {
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
	user.OrganizationName = organizationName
	user.Locations = locations

	return u.repo.Update(user)
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
