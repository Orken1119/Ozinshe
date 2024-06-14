package controller

import (
	"github.com/Orken1119/Ozinshe/internal/controller/middleware"
	pkg "github.com/Orken1119/Ozinshe/pkq"
	"github.com/gin-gonic/gin"

	"github.com/Orken1119/Ozinshe/internal/controller/auth"
	"github.com/Orken1119/Ozinshe/internal/controller/user"
	"github.com/Orken1119/Ozinshe/internal/repository"
)

func Setup(app pkg.Application, router *gin.Engine) {
	db := app.DB

	loginController := &auth.AuthController{
		UserRepository: repository.NewUserRepository(db),
	}

	userController := &user.UserController{
		UserRepository: repository.NewUserRepository(db),
	}

	router.POST("/signup", loginController.Signup)
	router.POST("/signin", loginController.Signin)

	router.Use(middleware.JWTAuth(`access-secret-key`))

	userRouter := router.Group("/user")
	{
		userRouter.GET("/profile", userController.GetProfile)
	}

}
