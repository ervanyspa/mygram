package service

import (
	"context"
	"errors"
	"fmt"
	"mygram/internal/model"
	"mygram/internal/repository"
	"mygram/pkg/helper"
	"time"
)

type UserService interface {	
	SignIn(ctx context.Context, userSignIn model.UserSignIn) (model.User, error)
	GetUsersById(ctx context.Context, id uint64) (model.User, error)
	EditUser(ctx context.Context, user model.User) (model.UserResponse, error)
	DeleteUsersById(ctx context.Context, id uint64) (model.User, error)

	SignUp(ctx context.Context, userSignUp model.UserSignUp) (model.User, error)	

	GenerateUserAccessToken(ctx context.Context, user model.User) (token string, err error)
}

type userServiceImpl struct{
	repo repository.UserQuery
}

func NewUserService(repo repository.UserQuery) UserService{
	return &userServiceImpl{repo: repo}

}

func (u *userServiceImpl) SignUp(ctx context.Context, userSignUp model.UserSignUp) (model.User, error) {
	user := model.User{
		Username: userSignUp.Username,
		Email:    userSignUp.Email,
	}

	dob, err := time.Parse("2006-01-02", userSignUp.DoB)
	if err != nil {
		return model.User{}, err
	}
	user.DoB = dob

	pass, err := helper.GenerateHash(userSignUp.Password)
	if err != nil {
		return model.User{}, err
	}

	user.Password = pass

	res, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return model.User{}, err
	}
	return res, err
}

func (u *userServiceImpl) SignIn(ctx context.Context, userSignIn model.UserSignIn) (model.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, userSignIn.Email)
	if err != nil {
		return model.User{}, err
	}
	if user.ID == 0 {
		return model.User{}, errors.New("email not found")
	}

	isValidLogin := helper.CheckPasswordHash(userSignIn.Password, user.Password)
	if !isValidLogin {
		return model.User{}, errors.New("password incorrect")
	}

	return user, nil
}

func (u *userServiceImpl) GetUsersById(ctx context.Context, id uint64) (model.User, error) {
	user, err := u.repo.GetUsersByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	return user, err
}

func (u *userServiceImpl) EditUser(ctx context.Context, user model.User) (model.UserResponse, error) {
	cekEmail, err := u.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return model.UserResponse{}, err
	}
	
	if cekEmail.ID != 0 && user.ID != cekEmail.ID {
		return model.UserResponse{}, errors.New("email already in use")
	}
	
	err = u.repo.EditUser(ctx, &user)
	if err != nil {
		return model.UserResponse{}, err
	}

	getDob, err := u.repo.GetUsersByID(ctx, user.ID)
	if err != nil {
		return model.UserResponse{}, err
	}

	user.UpdatedAt = time.Now()
	
	DobTime := getDob.DoB

	DobString := DobTime.Format("2006-01-02")

	UserResponse := model.UserResponse{}
	UserResponse.ID = user.ID
	UserResponse.Email = user.Email
	UserResponse.Username = user.Username
	UserResponse.DoB = DobString
	UserResponse.UpdatedAt = user.UpdatedAt
	return UserResponse, nil
}

func (u *userServiceImpl) DeleteUsersById(ctx context.Context, id uint64) (model.User, error) {
	user, err := u.repo.GetUsersByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	if user.ID == 0 {
		return model.User{}, nil
	}

	err = u.repo.DeleteUsersByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return user, err
}


func (u *userServiceImpl) GenerateUserAccessToken(ctx context.Context, user model.User) (token string, err error) {
	now := time.Now()

	claim := model.StandardClaim{
		Jti: fmt.Sprintf("%v", time.Now().UnixNano()),
		Iss: "go-middleware",
		Aud: "golang-006",
		Sub: "access-token",
		Exp: uint64(now.Add(time.Hour).Unix()),
		Iat: uint64(now.Unix()),
		Nbf: uint64(now.Unix()),
	}

	userClaim := model.AccessClaim{
		StandardClaim: claim,
		UserID:        user.ID,
		Username:      user.Username,
		Dob:           user.DoB,
	}

	token, err = helper.GenerateToken(userClaim)
	return
}

