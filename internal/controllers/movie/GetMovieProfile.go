package movie

import (
	"net/http"
	"strconv"

	models "github.com/Orken1119/Ozinshe/internal/models/movie_models"
	"github.com/gin-gonic/gin"
)

type MovieController struct {
	MovieRepository models.MovieRepository
}

// @Tags		movie
// @Accept		json
// @Produce	json
// @Param        id   path      int  true  "id"
// @Security	BearerAuth
// @Success	200		{object}	models.SuccessResponse
// @Failure	default	{object}	models.ErrorResponse
// @Router		/movie/movie-profile/{id} [get]
func (mr *MovieController) GetMovieProfile(c *gin.Context) {
	movieIDStr := c.Param("id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_GET_MOVIE_PROFILE",
					Message: "Can't get profile from db",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}

	profile, err := mr.MovieRepository.GetMovieProfile(c, movieID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_GET_MOVIE_PROFILE",
					Message: "Can't get profile from db",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}

	response := map[string]interface{}{
		"name":          profile.Name,
		"year_of_issue": profile.YearOfIssue,
		"director":      profile.Director,
		"producer":      profile.Producer,
		"hours":         profile.Hours,
		"minutes":       profile.Minutes,
		"description":   profile.Description,
		"image":         profile.Image,
		"similarMovies": profile.SimilarMovie,
		"seasonCount":   profile.Seasons,
		"episodeCount":  profile.Episodes,
		"categories":    profile.Categories,
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Result: response})
}
