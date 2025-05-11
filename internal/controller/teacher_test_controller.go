package controller

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/immxrtalbeast/plandstu/internal/domain"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TeacherTestController struct {
	teacherTestINT domain.TeacherTestInteractor
}

func NewTeacherTestController(teacherTestINT domain.TeacherTestInteractor) *TeacherTestController {
	return &TeacherTestController{teacherTestINT: teacherTestINT}
}

func (c *TeacherTestController) TeacherTest(ctx *gin.Context) {
	disciplineIDStr := ctx.Query("discipline_id")

	disciplineID, err := strconv.Atoi(disciplineIDStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing disciplineID", "detail": err.Error()})
		return
	}
	test, err := c.teacherTestINT.TeacherTests(ctx, disciplineID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Теста не существует."})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error getting tests", "detail": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"test": test})
}
func (c *TeacherTestController) DeleteTeacherTest(ctx *gin.Context) {
	testIDStr := ctx.Query("test_id")

	testID, err := uuid.Parse(testIDStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing testID", "detail": err.Error()})
		return
	}
	if err := c.teacherTestINT.DeleteTeacherTest(ctx, testID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error deleting test", "detail": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
func (c *TeacherTestController) UpdateTeacherTest(ctx *gin.Context) {
	type UpdateTeacherTestRequest struct {
		TestID       uuid.UUID      `json:"ID"`
		Test         datatypes.JSON `json:"DetailsJSONB"`
		Answers      datatypes.JSON `json:"Answers"`
		UpdatedAt    time.Time      `json:"CreatedAt"`
		DisciplineID int            `json:"DisciplineID"`
	}
	var request UpdateTeacherTestRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.teacherTestINT.UpdateTeacherTest(ctx, request.TestID, request.Test, request.Answers); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error creating test", "detail": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
func (c *TeacherTestController) CreateTeacherTest(ctx *gin.Context) {
	type CreateTeacherTestRequest struct {
		Test    datatypes.JSON `json:"test"`
		Answers datatypes.JSON `json:"answers"`
	}
	var request CreateTeacherTestRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	disciplineIDStr := ctx.Query("discipline_id")

	disciplineID, err := strconv.Atoi(disciplineIDStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing disciplineID", "detail": err.Error()})
		return
	}

	if err = c.teacherTestINT.CreateTeacherTest(ctx, request.Test, request.Answers, disciplineID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error creating test", "detail": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *TeacherTestController) RandomTestTest(ctx *gin.Context) {
	disciplineIDStr := ctx.Query("discipline_id")

	disciplineID, err := strconv.Atoi(disciplineIDStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing disciplineID", "detail": err.Error()})
		return
	}
	test, err := c.teacherTestINT.TeacherTestForUser(ctx, disciplineID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error getting test", "detail": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"test": test})
}
