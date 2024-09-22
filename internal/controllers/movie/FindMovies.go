package movie

import (
	"net/http"

	models "github.com/Orken1119/Ozinshe/internal/models/movie_models"
	"github.com/gin-gonic/gin"
)

// @Tags		movie
// @Accept		json
// @Produce	json
// @Param request body models.SearchRequest true "query params"
// @Security	BearerAuth
// @Success	200		{object}	models.SuccessResponse
// @Failure	default	{object}	models.ErrorResponse
// @Router		/movie/search [get]
func (mr *MovieController) FindMovie(c *gin.Context) {
	var searchRequest models.SearchRequest

	err := c.ShouldBindJSON(&searchRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_BIND_JSON",
					Message: "Datas dont match with struct of search",
				},
			},
		})
		return
	}

	movies, err := mr.MovieRepository.FindMovie(c, searchRequest.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_GET_MOVIES_WITH_THIS_NAME",
					Message: "Can't get movie with this name from db",
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
