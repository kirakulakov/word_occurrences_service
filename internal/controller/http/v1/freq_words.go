package v1

import (
	"fmt"
	"net/http"
	"npp_doslab/internal/entity"
	"npp_doslab/internal/usecase"
	"npp_doslab/pkg/logger"

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

type wordResponse struct {
	History []entity.Comm `json:"word"`
}

// @Summary     Get statistic
// @Description Get statistic of most frequently words.
// @ID          post_statistic
// @Tags  	    Post
// @Accept      json
// @Produce     json
// @Param       postId   path      int  true  "Id of specific post"
// @Success     200 {object} wordResponse
// @Failure     500 {object} response
// @Router      /post/{postId}/comments/statistics [get]
func (r *freqWordsRoutes) posts(c *gin.Context) {
	words, err := r.t.Scan(c.Request.Context())
	postId := c.Param("postId")
	fmt.Println(postId)

	if err != nil {
		r.l.Error(err, "http - v1 - posts")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	c.JSON(http.StatusOK, wordResponse{words})
}
