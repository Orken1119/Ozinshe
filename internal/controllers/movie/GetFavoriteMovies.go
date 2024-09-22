package movie

import (
	"net/http"

	models "github.com/Orken1119/Ozinshe/internal/models/movie_models"
	"github.com/gin-gonic/gin"
)

// @Tags		movie
// @Accept		json
// @Produce	json
// @Security	BearerAuth
// @Success	200		{object}	models.SuccessResponse
// @Failure	default	{object}	models.ErrorResponse
// @Router		/movie/favorite-movies [get]
func (mr *MovieController) GetFavoriteMovies(c *gin.Context) {
	UserID := c.GetUint("userID")

	profile, err := mr.MovieRepository.GetFavoriteMoviesByUser(c, UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_GET_FAVORITE_MOVIES",
					Message: "Can't get favorite movies from db",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Result: profile})
}
