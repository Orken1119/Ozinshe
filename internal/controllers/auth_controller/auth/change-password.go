package auth

import (
	"net/http"

	models "github.com/Orken1119/Ozinshe/internal/models/auth_models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (fp *AuthController) ChangePassword(c *gin.Context) {
	var request models.ChangePasswordRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_BIND_JSON",
					Message: "Datas dont match with struct of changePassword",
				},
			},
		})
		return
	}

	user, _ := fp.UserRepository.GetUserByEmail(c, request.Email)
	if user.ID > 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "USER_EXISTS",
					Message: "User with this email doesn't exists",
				},
			},
		})
		return
	}

	err = ValidatePassword(request.Password.Password)
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
	// Подтверждение пароля
	err = ConfirmPassword(request.Password.Password, request.Password.ConfirmPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_PASSWORD_MISMATCH",
					Message: err.Error(),
				},
			},
		})
		return
	}
	//

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password.Password), bcrypt.DefaultCost)
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
	request.Password.Password = string(encryptedPassword)

	code, err := fp.UserRepository.GetCodeByEmail(c, request.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_RECEIVE_CODE",
					Message: "Couldn't receive code",
				},
			},
		})
		return
	}

	err = fp.UserRepository.ChangePassword(c, code, request.Email, request.Password.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_CHANGE_RASSWORD",
					Message: "Couldn't change password",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}
}
