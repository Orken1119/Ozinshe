package movie

import (
	"net/http"
	"strconv"

	models "github.com/Orken1119/Ozinshe/internal/models/movie_models"
	"github.com/gin-gonic/gin"
)

// @Tags		movie
// @Accept		json
// @Produce	json
// @Param        id   path      int  true  "id"
// @Security	BearerAuth
// @Success	200		{object}	models.SuccessResponse
// @Failure	default	{object}	models.ErrorResponse
// @Router		/movie/get-episodes/{id} [get]
func (mr *MovieController) GetEpisodes(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_GET_MOVIE_PROFILE",
					Message: "Can't get episodes from db",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}
	episodes, err := mr.MovieRepository.GetEpisodes(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_GET_MOVIE_PROFILE",
					Message: "Can't get episodes from db",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Result: episodes})
}
