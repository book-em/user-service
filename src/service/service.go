package api

import (
	"bookem-user-service/client/roomclient"
	"bookem-user-service/domain"
	repo "bookem-user-service/repo"
	util "bookem-user-service/util"
	"context"
	"fmt"
	"strings"
)

type Service interface {
	Register(ctx context.Context, input *domain.UserCreateDTO) (*domain.User, error)
	Login(ctx context.Context, dto domain.LoginDTO) (string, error)
	Update(ctx context.Context, callerID uint, dto domain.UserUpdateDTO) (*domain.User, error)
	ChangePassword(ctx context.Context, callerID uint, dto domain.PasswordUpdateDTO) (*domain.User, error)
	FindById(ctx context.Context, id uint) (*domain.User, error)
	Delete(ctx context.Context, callerID uint, id uint) error

	/// canDeleteUser returns an error if the user cannot be deleted right now.
	/// The error specifies the reason why the operation cannot be done.
	canDeleteUser(ctx context.Context, user *domain.User) error
}

type service struct {
	repo       repo.Repository
	roomClient roomclient.RoomClient
}

func NewService(r repo.Repository, roomClient roomclient.RoomClient) Service {
	return &service{r, roomClient}
}

func (s *service) Register(ctx context.Context, dto *domain.UserCreateDTO) (*domain.User, error) {
	util.TEL.Push(ctx, "hash-password")
	defer util.TEL.Pop()

	hashed, err := util.HashPassword(dto.Password)
	if err != nil {
		util.TEL.Event("failed hashing password", err)
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

	util.TEL.Push(ctx, "db-query-user")
	defer util.TEL.Pop()
	existing, _ := s.repo.FindByUsernameOrEmail(dto.Username, dto.Email)
	if existing != nil {
		if existing.Username == dto.Username {
			util.TEL.Event("username exists", nil)
			return nil, domain.ErrUsernameExists
		}
		if existing.Email == dto.Email {
			util.TEL.Event("email exists", nil)
			return nil, domain.ErrEmailExists
		}
	}

	util.TEL.Push(ctx, "db-insert-user")
	defer util.TEL.Pop()
	err = s.repo.Create(user)
	if err != nil {
		util.TEL.Event("failed inserting user", err)
		return nil, fmt.Errorf("%w: %v", domain.ErrDBInternal, err)
	}

	return user, nil
}

// Login signs the user in.
// It can accept both an email or a username.
// On success, it returns a JWT string.
// On error, it returns an empty string.
func (s *service) Login(ctx context.Context, dto domain.LoginDTO) (string, error) {
	util.TEL.Push(ctx, "find-user")
	defer util.TEL.Pop()

	user, _ := s.repo.FindByUsernameOrEmail(dto.UsernameOrEmail, dto.UsernameOrEmail)

	if user == nil {
		util.TEL.Event(fmt.Sprintf("User %s not found", dto.UsernameOrEmail), nil)
		return "", domain.ErrLoginFailed
	}

	util.TEL.Push(ctx, "verify-password")
	defer util.TEL.Pop()

	err := util.VerifyPassword(user.Password, dto.Password)
	if err != nil {
		util.TEL.Event("Password verification failed", err)
		return "", domain.ErrLoginFailed
	}

	util.TEL.Push(ctx, "create-jwt")
	defer util.TEL.Pop()

	jwt, err := util.CreateJWT(int(user.ID), user.Username, user.Role)
	if err != nil {
		util.TEL.Event("JWT Creation failed", err)
		return "", domain.ErrLoginFailed
	}

	return jwt, nil
}

// Update updates the user (specified by his ID in the dto) with the new values
// in the DTO. Fields with null values are skipped.
func (s *service) Update(ctx context.Context, callerID uint, dto domain.UserUpdateDTO) (*domain.User, error) {
	util.TEL.Eventf("User %d wants to update user %d", nil, callerID, dto.Id)

	// Users can only update themselves.

	if callerID != dto.Id {
		util.TEL.Eventf("User %d trying to update someone else", nil, callerID)
		return nil, domain.ErrUnauthorized
	}

	// Search for the user.

	util.TEL.Push(ctx, "find-user")
	defer util.TEL.Pop()

	user, err := s.FindById(util.TEL.Ctx(), dto.Id)
	if err != nil {
		util.TEL.Eventf("User %d not found", err, dto.Id)
		return nil, domain.ErrNotFound
	}

	// Check if the username or email is already taken by someone else.

	util.TEL.Push(ctx, "assert-unique-credentials")
	defer util.TEL.Pop()

	if dto.Username != nil || dto.Email != nil {
		usernameSafe := ""
		if dto.Username != nil {
			usernameSafe = *dto.Username
		}
		emailSafe := ""
		if dto.Email != nil {
			emailSafe = *dto.Email
		}

		otherUserWithUsernameOrEmail, _ := s.repo.FindByUsernameOrEmailNotId(usernameSafe, emailSafe, dto.Id)

		if otherUserWithUsernameOrEmail != nil {
			if usernameSafe == otherUserWithUsernameOrEmail.Username {
				return nil, domain.ErrUsernameExists
			} else if emailSafe == otherUserWithUsernameOrEmail.Email {
				return nil, domain.ErrEmailExists
			} else {
				util.TEL.Eventf("DB found user matching [%s] or [%s] but the in-memory comparison failed.\nFound user: %+v", nil, usernameSafe, emailSafe, otherUserWithUsernameOrEmail)

				return nil, domain.ErrDBInternal
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

	util.TEL.Push(ctx, "update-user")
	defer util.TEL.Pop()

	err = s.repo.Update(user)
	if err != nil {
		util.TEL.Event("Could not update user in DB", err)
		return nil, err
	}

	return user, nil
}

// ChangePassword changes the user's password.
func (s *service) ChangePassword(ctx context.Context, callerID uint, dto domain.PasswordUpdateDTO) (*domain.User, error) {
	util.TEL.Eventf("User %d wants to change password of user %d", nil, callerID, dto.Id)

	// User can only change his own password.

	if callerID != dto.Id {
		util.TEL.Eventf("User %d trying to update someone else", nil, callerID)
		return nil, domain.ErrUnauthorized
	}

	// Search for the user.

	util.TEL.Push(ctx, "find-user")
	defer util.TEL.Pop()

	user, err := s.FindById(util.TEL.Ctx(), dto.Id)
	if err != nil {
		util.TEL.Eventf("User %d not found", err, dto.Id)
		return nil, err
	}

	// Check if confirm password is valid.

	util.TEL.Push(ctx, "password-validation")
	defer util.TEL.Pop()

	if dto.NewPasswordConfirm != dto.NewPassword {
		util.TEL.Eventf("Passwords do not match", nil)
		return nil, domain.ErrPasswordsNotMatch
	}

	// Check if old password is correct.

	err = util.VerifyPassword(user.Password, dto.OldPassword)
	if err != nil {
		util.TEL.Eventf("Old password is incorrect", err)
		return nil, domain.ErrWrongPassword
	}

	// Check if password is new.

	if dto.NewPassword == dto.OldPassword {
		util.TEL.Eventf("New password hasn't changed", nil)
		return nil, domain.ErrPasswordNotChanged
	}

	// Hash new password.

	passwordHashed, err := util.HashPassword(dto.NewPassword)
	if err != nil {
		util.TEL.Eventf("Password hashing failed", err)
		return nil, err
	}

	// Update.

	util.TEL.Push(ctx, "update-user")
	defer util.TEL.Pop()

	user.Password = passwordHashed

	err = s.repo.Update(user)
	if err != nil {
		util.TEL.Event("Could not update user in DB", err)
		return nil, err
	}

	return user, nil
}

func (s *service) FindById(ctx context.Context, id uint) (*domain.User, error) {
	util.TEL.Eventf("Find user by ID %d", nil, id)

	util.TEL.Push(ctx, "find-user-in-db")
	defer util.TEL.Pop()

	user, err := s.repo.FindById(id)
	if err != nil {
		util.TEL.Eventf("User %d not found", err, id)
		return nil, domain.ErrNotFound
	}
	return user, nil
}

func (s *service) Delete(ctx context.Context, callerID uint, id uint) error {
	util.TEL.Eventf("User %d wants to delete user %d", nil, callerID, id)

	// User can only delete himself.

	if id != callerID {
		util.TEL.Eventf("User %d trying to delete someone else", nil, callerID)
		return domain.ErrUnauthorized
	}

	// Search for the user.

	util.TEL.Push(ctx, "find-user")
	defer util.TEL.Pop()

	user, err := s.FindById(util.TEL.Ctx(), id)
	if err != nil {
		util.TEL.Eventf("User %d not found", err, id)
		return err
	}

	// Check if user can be deleted.

	util.TEL.Push(ctx, "delete-safety-check")
	defer util.TEL.Pop()

	err = s.canDeleteUser(util.TEL.Ctx(), user)
	if err != nil {
		util.TEL.Eventf("Cannot delete user %d", err, id)
		return err
	}

	// Delete user

	util.TEL.Push(ctx, "delete-user-in-db")
	defer util.TEL.Pop()

	s.repo.Delete(user.ID)
	util.TEL.Eventf("User %d deleted", nil, id)

	return nil
}

func (s *service) canDeleteUser(ctx context.Context, user *domain.User) error {
	util.TEL.Eventf("Check if user %d can be deleted", nil, user.ID)

	switch user.Role {
	case domain.Guest:
		util.TEL.Eventf("User is guest - must not have any reservations", nil)
		reservations, err := s.roomClient.GetPendingGuestReservations(ctx, user)
		if err != nil {
			util.TEL.Eventf("Could not check", err)
			return err
		}
		if len(reservations) > 0 {
			util.TEL.Eventf("Guest has reservations, cannot delete user", nil)
			return domain.ErrGuestHasReservations
		}
		return nil
	case domain.Host:
		util.TEL.Eventf("User is host - rooms must not have any reservations", nil)
		reservations, err := s.roomClient.GetActiveHostReservations(ctx, user)
		if err != nil {
			return err
		}
		if len(reservations) > 0 {
			util.TEL.Eventf("Host's rooms have reservations, cannot delete user", nil)
			return domain.ErrHostHasReservations
		}
		return nil
	default:
		util.TEL.Eventf("Users with role %d cannot be deleted", nil, user.Role)
		return domain.ErrCannotDeleteAdmin
	}
}
