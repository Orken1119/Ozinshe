package movie

import (
	"net/http"

	models "github.com/Orken1119/Ozinshe/internal/models/movie_models"
	"github.com/gin-gonic/gin"
)

// @Tags		movie
// @Accept		json
// @Produce	json
// @Param        category   path      string  true "category"
// @Security	BearerAuth
// @Success	200		{object}	models.SuccessResponse
// @Failure	default	{object}	models.ErrorResponse
// @Router		/movie/category/{category} [get]
func (mr *MovieController) GetMovieByCategory(c *gin.Context) {
	category := c.Param("category")
	movies, err := mr.MovieRepository.GetByCategory(c, category)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_GET_MOVIE_BY_CATEGORY",
					Message: "Can't get movie with this category",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Result: movies})
}
