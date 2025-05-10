package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/immxrtalbeast/plandstu/internal/domain"
)

type LLMController struct {
	llmURL     string
	interactor domain.LLMInteractor
}

func NewLLMController(llmURL string, LLMINT domain.LLMInteractor) *LLMController {
	return &LLMController{llmURL: llmURL, interactor: LLMINT}
}

func (c *LLMController) Chat(ctx *gin.Context) {
	type ChatRequest struct {
		UserID string `json:"user_id"`
		Prompt string `json:"prompt"`
	}
	message := ctx.Query("message")
	userIDStr, ok := ctx.Keys["userID"].(string)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing userID"})
		return
	}

	fastAPIReq := ChatRequest{
		UserID: userIDStr,
		Prompt: message,
	}

	requestBody, err := json.Marshal(fastAPIReq)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
		return
	}
	// Создаем запрос к FastAPI
	req, err := http.NewRequest("GET", c.llmURL+"chat-stream/", bytes.NewBuffer(requestBody))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Настраиваем заголовки SSE
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	// Потоково копируем данные
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			ctx.Writer.Write(buf[:n])
			ctx.Writer.Flush()
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

func (c *LLMController) SaveHistory(ctx *gin.Context) {
	var req domain.SaveHistoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Error while parsing userID", "details": err.Error()})
		return
	}

	id, err := c.interactor.SaveHistory(ctx, req, userID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Error while saving history", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"HistoryID": id,
	})

}

func (c *LLMController) History(ctx *gin.Context) {
	type HistoryRequest struct {
		UserID string `json:"user_id"`
	}
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	userIDStr, ok := ctx.Keys["userID"].(string)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing userID"})
		return
	}
	historyReq := HistoryRequest{
		UserID: userIDStr,
	}
	requestBody, err := json.Marshal(historyReq)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
		return
	}
	req, err := http.NewRequest("GET", c.llmURL+"get_history/", bytes.NewBuffer(requestBody))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create request", "details": err.Error()})
		return

	}
	req.Header.Set("User-Agent", "Parser/1.0")
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request failed", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorBody, _ := io.ReadAll(resp.Body)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unknown status code", "details": string(errorBody)})
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
	ctx.Data(resp.StatusCode, "application/json", data)
}

func (c *LLMController) ClearHistory(ctx *gin.Context) {
	type HistoryRequest struct {
		UserID string `json:"user_id"`
	}
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	userIDStr, ok := ctx.Keys["userID"].(string)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing userID"})
		return
	}
	historyReq := HistoryRequest{
		UserID: userIDStr,
	}
	requestBody, err := json.Marshal(historyReq)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
		return
	}
	req, err := http.NewRequest("GET", c.llmURL+"clear_history/", bytes.NewBuffer(requestBody))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create request", "details": err.Error()})
		return

	}
	req.Header.Set("User-Agent", "Parser/1.0")
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request failed", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorBody, _ := io.ReadAll(resp.Body)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unknown status code", "details": string(errorBody)})
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
	ctx.Data(resp.StatusCode, "application/json", data)
}
