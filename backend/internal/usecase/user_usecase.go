package usecase

import (
	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserUsecase struct {
	repo		repository.UserRepository
	roleRepo	repository.UserRoleRepository
	walletRepo	repository.WalletRepository
	emailSender	repository.EmailSender
}

func NewUserUsecase(
	r repository.UserRepository,
	roleRepo repository.UserRoleRepository,
	walletRepo repository.WalletRepository,
	emailSender repository.EmailSender,
) *UserUsecase {
	return &UserUsecase{repo: r, roleRepo: roleRepo, walletRepo: walletRepo, emailSender: emailSender}
}

func (u *UserUsecase) Register(email string) (*domain.User, error) {
	user := &domain.User{
		ID:		uuid.NewString(),
		Email:		email,
		IsActive:	true,
		EmailVerified:	false,
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

func (u *UserUsecase) GetWallet(userID string) (*domain.Wallet, error) {
	return u.walletRepo.GetWalletByUserID(userID)
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
				UserID:			userID,
				OrganizationName:	organizationName,
				Address:		address,
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
) ([]domain.AdminUserDetails, int64, error) {
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

func (u *UserUsecase) AdminApproveOrganizer(userID string) error {
	detail, err := u.repo.GetOrganizerDetails(userID)
	if err != nil {
		return err
	}
	if detail.ApprovalStatus == "approved" {
		return apiErrors.New(409, apiErrors.DuplicateResource, "Organizer is already approved")
	}

	if err := u.repo.UpdateOrganizerApprovalStatus(userID, "approved"); err != nil {
		return err
	}
	if err := u.roleRepo.AddRole(userID, domain.RoleOrganizer); err != nil {
		return err
	}

	if u.emailSender != nil {
		user, err := u.repo.FindByID(userID)
		if err == nil {
			if mailErr := u.emailSender.SendOrganizerApproval(user.Email, user.Name); mailErr != nil {
				logger.Log.Warn().Err(mailErr).Str("user_id", userID).Msg("failed to send organizer approval email")
			}
		}
	}
	return nil
}

func (u *UserUsecase) AdminRejectOrganizer(userID string) error {
	detail, err := u.repo.GetOrganizerDetails(userID)
	if err != nil {
		return err
	}
	if detail.ApprovalStatus == "rejected" {
		return apiErrors.New(409, apiErrors.DuplicateResource, "Organizer is already rejected")
	}

	return u.repo.UpdateOrganizerApprovalStatus(userID, "rejected")
}

func (u *UserUsecase) ListAdminUsers() ([]domain.User, error) {
	return u.repo.FindUsersByRole(domain.RoleAdmin)
}

func (u *UserUsecase) AddAdmin(name, email string) (*domain.User, error) {
	user := &domain.User{
		ID:		uuid.NewString(),
		Name:		name,
		Email:		email,
		IsActive:	true,
		EmailVerified:	true,
	}

	if err := u.repo.Create(user); err != nil {
		return nil, err
	}

	if err := u.roleRepo.AddRole(user.ID, domain.RoleAdmin); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) DeleteAdmin(adminID string) error {
	return u.repo.Delete(adminID)
}
