package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/immxrtalbeast/plandstu/internal/domain"
	"gorm.io/gorm"
)

type RoadmapController struct {
	interactor domain.RoadmapInteractor
}

func NewRoadmapController(RoadmapINT domain.RoadmapInteractor) *RoadmapController {
	return &RoadmapController{interactor: RoadmapINT}
}

func (c *RoadmapController) History(ctx *gin.Context) {
	userIDStr, ok := ctx.Keys["userID"].(string)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing userID"})
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing userID to uuid"})
		return
	}
	disciplineIDStr := ctx.Query("link")
	disciplineID, err := strconv.Atoi(disciplineIDStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing disciplineID", "detail": err.Error()})
		return
	}
	history, err := c.interactor.History(ctx, userID, disciplineID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "У пользователя нет истории"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error while getting roadmap history", "detail": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"roadmap_history": history})

}
func (c *RoadmapController) Report(ctx *gin.Context) {
	disciplineIDStr := ctx.Query("discipline_id")
	disciplineID, err := strconv.Atoi(disciplineIDStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing disciplineID", "detail": err.Error()})
		return
	}
	// themes := ctx.Query("discipline")
	userIDStr, ok := ctx.Keys["userID"].(string)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing userID"})
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing userID to uuid"})
		return
	}

	tests, err := c.interactor.Report(ctx, userID, disciplineID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error creating report"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"report": tests})
}
