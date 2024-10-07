package controller

import (
	_ "github.com/Orken1119/Ozinshe/docs"
	"github.com/Orken1119/Ozinshe/internal/controllers/auth_controller/middleware"
	"github.com/Orken1119/Ozinshe/internal/controllers/movie"
	pkg "github.com/Orken1119/Ozinshe/pkg"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"fmt"

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

	movieController := &movie.MovieController{
		MovieRepository: repository.NewMovieRepository(db),
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/forgot-password", loginController.ForgotPassword)
	router.POST("/change-forgotten-password", loginController.ChangeForgottenPassword)
	router.POST("/signup", loginController.Signup)
	router.POST("/signin", loginController.Signin)

	router.Use(middleware.JWTAuth(`access-secret-key`))

	fmt.Println("5154646")
	movieRouter := router.Group("/movie")
	{
		movieRouter.POST("/watch-movie/:id", movieController.WatchMovie)
		movieRouter.GET("/watched-movies", movieController.GetWatchedMovie)
		movieRouter.GET("/top-movies", movieController.GetTopMovies)
		movieRouter.GET("/movie-profile/:id", movieController.GetMovieProfile)
		movieRouter.GET("/category/:category", movieController.GetMovieByCategory)
		movieRouter.GET("/favorite-movies", movieController.GetFavoriteMovies)
		movieRouter.GET("/get-episodes/:id", movieController.GetEpisodes)
		movieRouter.POST("/search", movieController.FindMovie)
		movieRouter.POST("/add-in-favorits/:id", movieController.AddInFavorits)
	}

	userRouter := router.Group("/user")
	{
		userRouter.GET("/profile", userController.GetProfile)
		userRouter.POST("/edit-profile", userController.EditPersonalData)
		userRouter.PUT("/change-password", userController.ChangePassword)

	}

}
