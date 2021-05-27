package usecase

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"user/internal/app/models"
	review2 "user/internal/app/review"
	specialize2 "user/internal/app/specialize"
	"user/internal/app/user"
	"user/internal/app/user/tools"
	"user/internal/app/user/tools/passwordencrypt"
	customErrors "user/pkg/error"
	"user/pkg/types"
)

const (
	ctxParam types.CtxKey = 4
)

type UseCase struct {
	UserRepository       user.Repository
	SpecializeRepository specialize2.Repository
	ReviewsRepository    review2.Repository
	Encrypter            tools.PasswordEncrypter
}

func (useCase *UseCase) SetImg(ID uint64, img string, ctx context.Context) error {
	err := useCase.UserRepository.SetUserImg(ID, img, ctx)
	if err != nil {
		return err
	}
	return err
}

func (useCase *UseCase) Create(user models.NewUser, ctx context.Context) (*models.UserBasicInfo, error) {
	err := tools.ValidationCreateUser(&user)
	if err != nil {
		return nil, errors.Wrap(err, "Validation error")
	}
	if user, err = useCase.Encrypter.BeforeCreate(user); err != nil {
		return nil, errors.Wrap(err, "Encrypt password error")
	}
	ID, err := useCase.UserRepository.Create(user, ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error in user repository.")
	}

	for _, spec := range user.Specializes {
		specID, err := useCase.SpecializeRepository.FindByName(spec, ctx)
		if err != nil {
			if errors.Is(err, customErrors.ErrorNoRows) {
				specID, err = useCase.SpecializeRepository.Create(spec, ctx)
				if err != nil {
					return nil, errors.Wrap(err, "Error in data sourse")
				}
			}
		}
		err = useCase.SpecializeRepository.AssociateSpecializationWithUser(specID, ID, ctx)
		if err != nil {
			return nil, err
		}
	}

	return &models.UserBasicInfo{
		ID:       ID,
		Executor: user.Executor,
	}, nil
}

func (useCase *UseCase) Verification(email string,
	password string, ctx context.Context) (*models.UserBasicInfo, error) {
	user, err := useCase.UserRepository.FindUserByEmail(email, ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error in user repository.")
	}
	if !useCase.Encrypter.CompPass(user.Password, password) {
		return nil, customErrors.ErrorInvalidPassword
	}

	return &models.UserBasicInfo{
		ID:       user.ID,
		Executor: user.Executor,
	}, nil
}

func (useCase *UseCase) GetById(ID uint64, ctx context.Context) (*models.UserInfo, error) {
	userInfo, err := useCase.UserRepository.FindUserByID(ID, ctx)
	if err != nil {
		return nil, err
	}
	if userInfo.Executor {
		userInfo.Specializes, err = useCase.SpecializeRepository.FindByUserID(userInfo.ID, ctx)
		if err != nil {
			return nil, errors.Wrap(err, "Error in specialize repository.")
		}
	}
	userReviewInfo, err := useCase.ReviewsRepository.GetAvgScoreByUserId(ID, ctx)
	if err != nil {
		return nil, err
	}
	userInfo.Rating = userReviewInfo.Rating
	userInfo.ReviewCount = userReviewInfo.ReviewCount
	return userInfo, nil
}

func (useCase *UseCase) Change(user models.ChangeUser, ctx context.Context) (*models.UserBasicInfo, error) {
	fmt.Println(user)
	err := tools.ValidationChangeUser(user)
	if err != nil {
		return nil, err
	}
	oldUser, err := useCase.UserRepository.FindUserByID(user.ID, ctx)
	if err != nil {
		return nil, errors.Wrap(err, "User db error")
	}

	if !useCase.Encrypter.CompPass(oldUser.Password, user.Password) {
		return nil, customErrors.ErrorInvalidPassword
	}
	if user.NewPassword != "" {
		user.Password = user.NewPassword
	}
	fmt.Println(user)
	if user, err = useCase.Encrypter.BeforeChange(user); err != nil {
		return nil, err
	}
	if user.Email == "" {
		user.Email = oldUser.Email
	}

	if user.About == "" {
		user.About = oldUser.About
	}

	if user.Password == "" {
		user.EncryptPassword = oldUser.Password
	}

	if user.Login == "" {
		user.Login = oldUser.Login
	}

	if user.NameSurname == "" {
		user.NameSurname = oldUser.NameSurname
	}

	for i := range oldUser.Specializes {
		user.Specializes = append(user.Specializes, oldUser.Specializes[i])
	}
	err = useCase.UserRepository.Change(user, ctx)
	if err != nil {
		return nil, err
	}
	for i := range user.Specializes {
		specID, err := useCase.SpecializeRepository.FindByName(user.Specializes[i], ctx)
		if err != nil {
			if err == customErrors.ErrorNoRows {
				specID, err = useCase.SpecializeRepository.Create(user.Specializes[i], ctx)
				if err != nil {
					return nil, errors.Wrap(err, "Error in data sourse")
				}
			}
		}
		err = useCase.SpecializeRepository.AssociateSpecializationWithUser(specID, oldUser.ID, ctx)
		if err != nil {
			return nil, err
		}
	}

	return &models.UserBasicInfo{
		ID:          user.ID,
		Email:       user.Email,
		About:       user.About,
		Executor:    user.Executor,
		Login:       user.Login,
		NameSurname: user.NameSurname,
	}, nil
}

func (u *UseCase) SuggestUsersTitle(suggestWord string, ctx context.Context) ([]models.SuggestUsersTittle, error) {
	suggestTittles, err := u.UserRepository.SuggestUsersTitle(suggestWord, ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error in data sourse")
	}
	if suggestTittles == nil {
		return []models.SuggestUsersTittle{}, nil
	}
	return suggestTittles, nil
}

func (useCase *UseCase) GetUsers(ctx context.Context) ([]models.UserInfo, error) {
	uInf, err := useCase.UserRepository.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	param := ctx.Value(ctxParam).(map[string]interface{})
	category := param["category"].(string)

	var res []models.UserInfo

	for i := range uInf {
		uInf[i].Specializes, err = useCase.SpecializeRepository.FindByUserID(uInf[i].ID, ctx)
		if uInf[i].Specializes == nil {
			uInf[i].Specializes = []string{}
		}
		if err != nil {
			return nil, err
		}
		if category != "" {
			flag := false
			for _, spec := range uInf[i].Specializes {
				if spec == category {
					flag = true
				}
			}
			if flag {
				res = append(res, uInf[i])
				//uInf[i], uInf[len(uInf)-1] = uInf[len(uInf)-1], uInf[i]
				//uInf = uInf[:len(uInf)-1]
			}
		} else {
			res = append(res, uInf[i])
		}
	}
	if res == nil {
		return []models.UserInfo{}, nil
	}
	return res, nil
}

func New(userRep user.Repository, specRep specialize2.Repository, reviewsRepository review2.Repository) *UseCase {
	return &UseCase{
		SpecializeRepository: specRep,
		UserRepository:       userRep,
		ReviewsRepository:    reviewsRepository,
		Encrypter:            &passwordencrypt.PasswordEncrypter{},
	}
}
