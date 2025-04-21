package controller

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ParserController struct {
	parserURL string
}

func NewParserController(parserURL string) *ParserController {
	return &ParserController{parserURL: parserURL}
}

func (c *ParserController) Faculties(ctx *gin.Context) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", c.parserURL+"api/faculties/", nil)
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unknown status code", "details": err.Error()})
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
func (c *ParserController) FacultyByID(ctx *gin.Context) {
	id := ctx.Param("id")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", c.parserURL+"api/faculties/"+id, nil)
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unknown status code", "details": err.Error()})
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

func (c *ParserController) Disciplines(ctx *gin.Context) {
	direction := ctx.Param("direction")
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	req, err := http.NewRequest("GET", c.parserURL+"api/get_disciplines/"+direction, nil)
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unknown status code", "details": err.Error()})
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

func (c *ParserController) Roadmap(ctx *gin.Context) {
	discipline := ctx.Param("discipline")
	link := ctx.Param("link")
	client := &http.Client{
		Timeout: 40 * time.Second,
	}

	req, err := http.NewRequest("GET", c.parserURL+"api/roadmaps/"+discipline+"/"+link, nil)
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unknown status code", "details": err.Error()})
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
