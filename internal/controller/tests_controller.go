package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/immxrtalbeast/plandstu/internal/domain"
	"github.com/immxrtalbeast/plandstu/internal/task"
	"gorm.io/datatypes"
)

type TestsController struct {
	llmURL         string
	roadmapINT     domain.RoadmapInteractor
	testINT        domain.TestInteractor
	redisURL       string
	teacherTestINT domain.TeacherTestInteractor
}

func NewTestsController(llmURL string, roadmapINT domain.RoadmapInteractor, testINT domain.TestInteractor, redisURL string, teacherTestINT domain.TeacherTestInteractor) *TestsController {
	return &TestsController{llmURL: llmURL, roadmapINT: roadmapINT, testINT: testINT, redisURL: redisURL, teacherTestINT: teacherTestINT}
}

// TODO: Можно объеденить FirtsTest и просто Test
func (c *TestsController) FirstTest(ctx *gin.Context) {
	type GenerateTestRequest struct {
		TestID uuid.UUID `json:"test_id"`
		Themes []string  `json:"themes"`
	}

	type CreateTestRequets struct {
		Themes []string `json:"themes"`
	}
	var request CreateTestRequets
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}
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
	var history *domain.RoadmapHistory
	existingHistory, err := c.roadmapINT.History(ctx, userID, disciplineID)
	if err != nil {
		newHistory, err := c.roadmapINT.CreateHistory(ctx, userID, disciplineID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error creating history", "detail": err.Error()})
			return
		}
		history = newHistory
	} else {
		history = existingHistory
	}
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	teacherTest, _ := c.teacherTestINT.TeacherTestForUser(ctx, disciplineID)
	if teacherTest != nil {
		generatedTestID := uuid.New()
		result, err := c.SendAnswers(ctx, client, teacherTest, history.ID, generatedTestID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error with preload test", "detail": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, result)
		return
	}

	generatedTestID := uuid.New()
	genReq := GenerateTestRequest{
		TestID: generatedTestID,
		Themes: request.Themes,
	}
	requestBody, err := json.Marshal(genReq)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
		return
	}
	req, err := http.NewRequest("GET", c.llmURL+"test-exmpl/", bytes.NewBuffer(requestBody))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create request", "details": err.Error()})
		return

	}
	req.Header.Set("User-Agent", "LLM/1.0")
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request failed", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unknown status code", "details": resp.Status})
		return
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to read response body",
			"details": err.Error(),
		})
		return
	}
	if !json.Valid(data) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON format in response",
		})
		return
	}
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to parse response JSON",
			"details": err.Error(),
		})
		return
	}
	jsonData["id"] = generatedTestID.String() // Добавление ID
	modifiedData, err := json.Marshal(jsonData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to marshal modified JSON",
			"details": err.Error(),
		})
		return
	}
	_, err = c.testINT.CreateTest(ctx, generatedTestID, datatypes.JSON(data), history.ID, true)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error creating test", "detail": err.Error()})
		return
	}
	ctx.Data(resp.StatusCode, "application/json", modifiedData)
}

func (c *TestsController) CreateTest(ctx *gin.Context) {
	type GenerateTestRequest struct {
		TestID uuid.UUID `json:"test_id"`
		Themes []string  `json:"themes"`
	}
	type CreateTestRequets struct {
		Themes []string `json:"themes"`
	}
	var request CreateTestRequets
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}
	disciplineIDStr := ctx.Query("discipline_id")
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
	history, err := c.roadmapINT.History(ctx, userID, disciplineID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error creating history", "detail": err.Error()})
		return
	}

	generatedTestID := uuid.New()
	payload := task.GenerateTestPayload{
		TestID:        generatedTestID.String(),
		Themes:        request.Themes,
		UserID:        userID.String(),
		DisciplineID:  disciplineID,
		HistoryID:     history.ID,
		LLMServiceURL: c.llmURL,
	}
	// [4] Ставим задачу в очередь
	t, err := task.NewGenerateTestTask(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}

	info, err := task.RedisClient.Enqueue(t)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue task"})
		return
	}

	// [6] Возвращаем ответ сразу
	ctx.JSON(http.StatusAccepted, gin.H{
		"task_id": info.ID,
	})
	// genReq := GenerateTestRequest{
	// 	TestID: generatedTestID,
	// 	Themes: request.Themes,
	// }
	// client := &http.Client{
	// 	Timeout: 60 * time.Minute,
	// }
	// requestBody, err := json.Marshal(genReq)
	// if err != nil {
	// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
	// 	return
	// }
	// req, err := http.NewRequest("POST", c.llmURL+"api/test-workflow", bytes.NewBuffer(requestBody))
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create request", "details": err.Error()})
	// 	return

	// }
	// req.Header.Set("User-Agent", "LLM/1.0")
	// resp, err := client.Do(req)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request failed", "details": err.Error()})
	// 	return
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode < 200 || resp.StatusCode >= 300 {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unknown status code", "details": resp.Status})
	// 	return
	// }
	// data, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"error":   "Failed to read response body",
	// 		"details": err.Error(),
	// 	})
	// 	return
	// }
	// if !json.Valid(data) {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Invalid JSON format in response",
	// 	})
	// 	return
	// }
	// var jsonData map[string]interface{}
	// if err := json.Unmarshal(data, &jsonData); err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"error":   "Failed to parse response JSON",
	// 		"details": err.Error(),
	// 	})
	// 	return
	// }
	// jsonData["id"] = generatedTestID.String() // Добавление ID
	// modifiedData, err := json.Marshal(jsonData)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"error":   "Failed to marshal modified JSON",
	// 		"details": err.Error(),
	// 	})
	// 	return
	// }
	// _, err = c.testINT.CreateTest(ctx, generatedTestID, datatypes.JSON(data), history.ID, false)
	// if err != nil {
	// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error creating test", "detail": err.Error()})
	// 	return
	// }
	// ctx.Data(resp.StatusCode, "application/json", modifiedData)
}

func (c *TestsController) Answers(ctx *gin.Context) {
	type AnswersRequest struct {
		TestID  uuid.UUID `json:"test_id"`
		Answers []string  `json:"answers"`
	}
	var req AnswersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}
	test_result, err := c.testINT.Answers(ctx, req.TestID, req.Answers)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to post answers",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, datatypes.JSON(test_result))
}
func (c *TestsController) GetTaskStatus(ctx *gin.Context) {
	taskID := ctx.Query("task_id")

	inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: c.redisURL})
	taskInfo, err := inspector.GetTaskInfo("default", taskID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": taskInfo.State.String(),
	})
}

func (c *TestsController) MyHistory(ctx *gin.Context) {
	disciplineIDStr := ctx.Query("discipline_id")
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

	tests, err := c.roadmapINT.Report(ctx, userID, disciplineID)
	if err != nil {
		var tests []string
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"history": tests, "error": "Error getting tests"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"history": tests})
}

func (c *TestsController) SendAnswers(ctx *gin.Context, client *http.Client, teacherTest *domain.TestResponse, historyID uuid.UUID, generatedTestID uuid.UUID) (*domain.TestResponse, error) {
	type SetAnswersRequest struct {
		TestID  uuid.UUID `json:"test_id"`
		Answers []string  `json:"answers"`
	}
	teacherTestFull, err := c.teacherTestINT.TeacherTestByID(ctx, teacherTest.ID)
	if err != nil {
		return nil, err
	}
	var answers []string

	if err := json.Unmarshal(teacherTestFull.Answers, &answers); err != nil {
		return nil, err
	}
	setAnswersReq := SetAnswersRequest{
		TestID:  generatedTestID,
		Answers: answers,
	}
	requestBody, err := json.Marshal(setAnswersReq)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.llmURL+"test/set-answers", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "LLM/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, err
	}
	wrappedJSON := []byte(`{"test":` + string(teacherTestFull.DetailsJSONB) + `}`)

	// Проверяем валидность JSON
	if !json.Valid(wrappedJSON) {
		return nil, fmt.Errorf("failed to wrapp json")
	}
	wrappedDetails := datatypes.JSON(wrappedJSON)
	_, err = c.testINT.CreateTest(ctx, generatedTestID, wrappedDetails, historyID, true)
	if err != nil {
		return nil, err
	}
	teacherTest.ID = generatedTestID
	return teacherTest, err
}
