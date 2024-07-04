package controller

import (
	"github.com/Orken1119/Ozinshe/internal/controllers/auth_controller/middleware"
	pkg "github.com/Orken1119/Ozinshe/pkg"
	"github.com/gin-gonic/gin"

	"github.com/Orken1119/Ozinshe/internal/controllers/auth_controller/auth"
	"github.com/Orken1119/Ozinshe/internal/controllers/auth_controller/user"
	repository "github.com/Orken1119/Ozinshe/internal/repositories"
)

func Setup(app pkg.Application, router *gin.Engine) {
	db := app.DB

	loginController := &auth.AuthController{
		UserRepository: repository.NewUserRepository(db),
	}

	userController := &user.UserController{
		UserRepository: repository.NewUserRepository(db),
	}

	router.POST("/forgot-password", loginController.ForgotPassword)
	router.POST("/change-password", loginController.ChangePassword)
	router.POST("/signup", loginController.Signup)
	router.POST("/signin", loginController.Signin)

	router.Use(middleware.JWTAuth(`access-secret-key`))

	userRouter := router.Group("/user")
	{
		userRouter.GET("/profile", userController.GetProfile)
	}

}
