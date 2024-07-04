package auth

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/smtp"

	models "github.com/Orken1119/Ozinshe/internal/models/auth_models"
	"github.com/gin-gonic/gin"
)

func (cp *AuthController) ForgotPassword(c *gin.Context) {
	var request models.ForgotPasswordRequest

	code := generateCode()

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_BIND_JSON",
					Message: "Datas dont match with struct of forgotPassword",
				},
			},
		})
		return
	}

	if request.Email == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "EMPTY_VALUES",
					Message: "Empty values are written in the form",
				},
			},
		})
		return
	}

	// Sender data.
	from := "ork.en389@gmail.com"
	password := "hvcl fdqt krgt hlji"

	// Receiver email address.
	to := []string{
		request.Email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	message := []byte("Subject: Your Authentication Code\n\nYour code is: " + code)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	sendingErr := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if sendingErr != nil {
		fmt.Println(sendingErr)
		return
	}
	email := request.Email
	err = cp.UserRepository.CreatePasswordResetCode(c, email, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_WITH_DATA",
					Message: "Invalid code was sent",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}
	fmt.Println("Email Sent Successfully!")
}

func generateCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
