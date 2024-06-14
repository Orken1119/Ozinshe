package auth

import (
	"fmt"
	"net/http"
	"unicode"

	"github.com/Orken1119/Ozinshe/internal/controller/tokenutil"
	"github.com/Orken1119/Ozinshe/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	UserRepository models.UserRepository
}

func (uc AuthController) Signup(c *gin.Context) {
	var request models.UserRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_BIND_JSON",
					Message: "Datas dont match with struct of signup",
				},
			},
		})
		return
	}
	user, _ := uc.UserRepository.GetUserByEmail(c, request.Email)
	if user.ID > 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "USER_EXISTS",
					Message: "User with this email already exists",
				},
			},
		})
		return
	}

	err := validatePassword(request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_PASSWORD_FORMAT",
					Message: err.Error(),
				},
			},
		})
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_ENCRYPTE_PASSWORD",
					Message: "Couldn't encrypte password",
				},
			},
		})
		return
	}
	request.Password = string(encryptedPassword)

	_, err = uc.UserRepository.CreateUser(c, request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_CREATE_USERS",
					Message: "Couldn't create user",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}

	user, err = uc.UserRepository.GetUserByEmail(c, request.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_GET_USER",
					Message: "User with this email doesn't found",
				},
			},
		})
		return
	}

	accessToken, err := tokenutil.CreateAccessToken(&user, `access-key`, 50)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "TOKEN_ERROR",
					Message: "Error to create access token",
				},
			},
		})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Result: accessToken})

}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	var (
		hasUpper, hasLower, hasDigit bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		return fmt.Errorf("password must contain at least one uppercase letter, one lowercase letter and one digit")
	}
	return nil
}
