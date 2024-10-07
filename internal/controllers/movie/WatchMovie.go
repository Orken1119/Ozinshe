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
// @Security Bearer
// @Success	200		{object}	models.SuccessResponse
// @Failure	default	{object}	models.ErrorResponse
// @Router		/movie/watch-movie/{id} [post]
func (mr *MovieController) WatchMovie(c *gin.Context) {
	userID := c.GetUint("userID")

	episodeIDStr := c.Param("id")
	episodeID, err := strconv.Atoi(episodeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode"})
		return
	}

	err = mr.MovieRepository.AddWatchedMovie(c, int(userID), episodeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_WATCH_MOVIES",
					Message: "Can't watch movie",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Result: "Movie marked as watched"})

}
