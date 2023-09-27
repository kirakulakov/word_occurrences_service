package v1

import (
	"fmt"
	"net/http"
	"npp_doslab/internal/entity"
	"npp_doslab/internal/usecase"
	"npp_doslab/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
)

type freqWordsRoutes struct {
	t usecase.FrequentlyWordsUseCase
	l logger.Interface
}

func NewWPostRouter(handler *gin.RouterGroup, t usecase.FrequentlyWordsUseCase, l logger.Interface) {
	r := &freqWordsRoutes{t, l}

	h := handler.Group("/post")
	{
		h.GET("/:postId/comments/statistics", r.posts)
	}
}

type commentResponse struct {
	Words []entity.Comm `json:"words"`
}

// @Summary     Get statistic
// @Description Get statistic of most frequently words.
// @ID          post_statistic
// @Tags  	    Post
// @Accept      json
// @Produce     json
// @Param       postId   path      int  true  "Id of specific post"
// @Success     200 {object} commentResponse
// @Failure     500 {object} response
// @Router      /post/{postId}/comments/statistics [get]
func (r *freqWordsRoutes) posts(c *gin.Context) {
	postId := c.Param("postId")
	postIdInt, err := strconv.Atoi(postId)

	words, err := r.t.GetByPostId(postIdInt, c.Request.Context())
	fmt.Println(postId)

	if err != nil {
		r.l.Error(err, "http - v1 - posts")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	c.JSON(http.StatusOK, commentResponse{words})
}
