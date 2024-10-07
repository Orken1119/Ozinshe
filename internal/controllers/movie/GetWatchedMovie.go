package movie

import (
	"net/http"

	models "github.com/Orken1119/Ozinshe/internal/models/movie_models"
	"github.com/gin-gonic/gin"
)

// @Tags		movie
// @Accept		json
// @Produce	json
// @Security Bearer
// @Success	200		{object}	models.SuccessResponse
// @Failure	default	{object}	models.ErrorResponse
// @Router		/movie/watched-movies [get]
func (mr *MovieController) GetWatchedMovie(c *gin.Context) {
	UserID := c.GetUint("userID")

	profile, err := mr.MovieRepository.GetWatchedMoviesByUser(c, UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_GET_WATCHED_MOVIES",
					Message: "Can't get watched movie from db",
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
