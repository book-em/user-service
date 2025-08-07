package api

import (
	"bookem-user-service/domain"
	repo "bookem-user-service/repo"
	util "bookem-user-service/util"
	"fmt"
	"log"
	"strings"
)

type Service interface {
	Register(input *domain.UserDTO) (*domain.User, error)
	Login(dto domain.LoginDTO) (string, error)
	Update(callerID uint, dto domain.UserUpdateDTO) (*domain.User, error)
	ChangePassword(callerID uint, dto domain.PasswordUpdateDTO) (*domain.User, error)
}

type service struct {
	repo repo.Repository
}

func NewService(r repo.Repository) Service {
	return &service{r}
}

func (s *service) Register(dto *domain.UserDTO) (*domain.User, error) {
	hashed, err := util.HashPassword(dto.Password)
	if err != nil {
		return nil, domain.ErrHashingPassword
	}

	user := &domain.User{
		Username: dto.Username,
		Password: hashed,
		Email:    strings.ToLower(dto.Email),
		Name:     dto.Name,
		Surname:  dto.Surname,
		Role:     domain.UserRole(dto.Role),
		Address:  dto.Address,
	}

	existing, _ := s.repo.FindByUsernameOrEmail(dto.Username, dto.Email)
	if existing != nil {
		if existing.Username == dto.Username {
			return nil, domain.ErrUsernameExists
		}
		if existing.Email == dto.Email {
			return nil, domain.ErrEmailExists
		}
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrDBInternal, err)
	}

	return user, nil
}

// Login signs the user in.
// It can accept both an email or a username.
// On success, it returns a JWT string.
// On error, it returns an empty string.
func (s *service) Login(dto domain.LoginDTO) (string, error) {
	user, _ := s.repo.FindByUsernameOrEmail(dto.UsernameOrEmail, dto.UsernameOrEmail)

	if user == nil {
		log.Printf("User %s not found", dto.UsernameOrEmail)
		return "", domain.ErrLoginFailed
	}

	err := util.VerifyPassword(user.Password, dto.Password)
	if err != nil {
		log.Print(err)
		return "", domain.ErrLoginFailed
	}

	jwt, err := util.CreateJWT(int(user.ID), user.Username, user.Role)
	if err != nil {
		log.Print(err)
		return "", domain.ErrLoginFailed
	}

	return jwt, nil
}

// Update updates the user (specified by his ID in the dto) with the new values
// in the DTO. Fields with null values are skipped.
func (s *service) Update(callerID uint, dto domain.UserUpdateDTO) (*domain.User, error) {
	log.Printf("User %d wants to update user %d", callerID, dto.Id)

	// Users can only update themselves.

	if callerID != dto.Id {
		return nil, domain.ErrUnauthorized
	}

	// Search for the user.

	user, err := s.repo.FindById(dto.Id)
	if err != nil {
		return nil, err
	}

	// Check if the username or email is already taken by someone else.

	if dto.Username != nil || dto.Email != nil {
		usernameSafe := ""
		if dto.Username != nil {
			usernameSafe = *dto.Username
		}
		emailSafe := ""
		if dto.Email != nil {
			emailSafe = *dto.Email
		}

		existing, _ := s.repo.FindByUsernameOrEmail(usernameSafe, emailSafe)

		if existing != nil && existing.ID != dto.Id {
			if dto.Username != nil {
				return nil, domain.ErrUsernameExists
			}
			if dto.Email != nil {
				return nil, domain.ErrEmailExists
			}
		}
	}

	// Update non-null fields.

	if dto.Username != nil {
		user.Username = *dto.Username
	}
	if dto.Email != nil {
		user.Email = *dto.Email
	}
	if dto.Name != nil {
		user.Name = *dto.Name
	}
	if dto.Surname != nil {
		user.Surname = *dto.Surname
	}
	if dto.Address != nil {
		user.Address = *dto.Address
	}

	err = s.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes the user's password.
func (s *service) ChangePassword(callerID uint, dto domain.PasswordUpdateDTO) (*domain.User, error) {
	log.Printf("User %d wants to change password of user %d", callerID, dto.Id)

	// User can only change his own password.

	if callerID != dto.Id {
		return nil, domain.ErrUnauthorized
	}

	// Search for the user.

	user, err := s.repo.FindById(dto.Id)
	if err != nil {
		return nil, err
	}

	// Check if confirm password is valid.

	if dto.NewPasswordConfirm != dto.NewPassword {
		return nil, domain.ErrPasswordsNotMatch
	}

	// Check if old password is correct.

	err = util.VerifyPassword(user.Password, dto.OldPassword)
	if err != nil {
		return nil, err
	}

	// Check if password is new.

	if dto.NewPassword == dto.OldPassword {
		return nil, domain.ErrPasswordNotChanged
	}

	// Hash new password.

	passwordHashed, err := util.HashPassword(dto.NewPassword)
	if err != nil {
		return nil, err
	}

	// Update.

	user.Password = passwordHashed

	err = s.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
