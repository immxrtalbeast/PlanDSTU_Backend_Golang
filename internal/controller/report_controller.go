package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/immxrtalbeast/plandstu/internal/domain"
	"gorm.io/datatypes"
)

type ReportController struct {
	reportINT  domain.ReportInteractor
	roadmapINT domain.RoadmapInteractor
	userINT    domain.UserInteractor
}

func NewReportController(reportINT domain.ReportInteractor, roadmapINT domain.RoadmapInteractor, userINT domain.UserInteractor) *ReportController {
	return &ReportController{reportINT: reportINT, roadmapINT: roadmapINT, userINT: userINT}
}

func (c *ReportController) CreateReport(ctx *gin.Context) {
	disciplineIDStr := ctx.Query("discipline_id")
	disciplineName := ctx.Query("discipline_name")

	disciplineID, err := strconv.Atoi(disciplineIDStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing disciplineID", "detail": err.Error()})
		return
	}
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
	user, err := c.userINT.User(ctx, userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error getting user"})
		return
	}
	tests, err := c.roadmapINT.Report(ctx, userID, disciplineID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error getting tests"})
		return
	}
	reportData := gin.H{"report": tests}
	jsonBytes, err := json.Marshal(reportData)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling report data"})
		return
	}

	err = c.reportINT.CreateReport(ctx, disciplineID, datatypes.JSON(jsonBytes), userID, disciplineName, user.Group)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error creating report"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"report": tests})
}

func (c *ReportController) Report(ctx *gin.Context) {
	reportIDStr := ctx.Query("report_id")
	reportID, err := uuid.Parse(reportIDStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing reportID to uuid"})
		return
	}
	report, err := c.reportINT.Report(ctx, reportID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error getting report"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"report": report})

}
func (c *ReportController) ReportsByDisciplineID(ctx *gin.Context) {
	disciplineIDStr := ctx.Query("discipline_id")
	disciplineID, err := strconv.Atoi(disciplineIDStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing disciplineID", "detail": err.Error()})
		return
	}
	reports, err := c.reportINT.ReportsByDisciplineID(ctx, disciplineID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error getting report"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"reports": reports})

}
func (c *ReportController) ReportsDisciplines(ctx *gin.Context) {
	disciplines, err := c.reportINT.ReportDisciplines(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error getting disciplines"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"disciplines": disciplines})
}

func (c *ReportController) ReportsGroup(ctx *gin.Context) {
	disciplineName := ctx.Query("discipline_name")
	groups, err := c.reportINT.ReportGroups(ctx, disciplineName)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error getting disciplines"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"groups": groups})
}

func (c *ReportController) ReportsByGroupAndDiscipline(ctx *gin.Context) {
	disciplineName := ctx.Query("discipline_name")
	groupName := ctx.Query("group")
	reports, stats, err := c.reportINT.ReportsByGroupAndDiscipline(ctx, disciplineName, groupName)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error getting disciplines"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"reports": reports, "stats": stats})
}
