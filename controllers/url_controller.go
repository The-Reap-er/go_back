package controllers

import (
	"context"
	"encoding/json"
	"go_back/config"
	"go_back/database"
	"go_back/models"
	"io"
	"net/http"
	"strings"
	"time"

	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLInput struct {
	Service string `json:"service" binding:"required"`
	URLList string `json:"url_list" binding:"required"`
}

func CreateURL(c *gin.Context) {
	var input URLInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Ensure userID is of type primitive.ObjectID
	userIDObj, ok := userID.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID has invalid type"})
		return
	}

	url := models.URL{
		UserID:  userIDObj,
		Service: input.Service,
		URLList: input.URLList,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := database.URLCollection.InsertOne(ctx, url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create URL list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": result.InsertedID})
}

func GetURLs(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDStr.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID has invalid type"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.URLCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URLs"})
		return
	}
	defer cursor.Close(ctx)

	var urls []models.URL
	if err = cursor.All(ctx, &urls); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse URLs"})
		return
	}

	c.JSON(http.StatusOK, urls)
}

func DeleteURL(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL ID"})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, ok := userIDStr.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID has invalid type"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := database.URLCollection.DeleteOne(ctx, bson.M{"_id": objID, "user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete URL"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL deleted"})
}

func UpdateURL(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL ID"})
		return
	}

	var input URLInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, ok := userIDStr.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID has invalid type"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": objID, "user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"service":  input.Service,
			"url_list": input.URLList,
			// Add other fields if necessary
		},
	}

	result, err := database.URLCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update URL"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL updated"})
}

func GetURLsByService(c *gin.Context) {
	serviceName := c.Param("service")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service name is required"})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDStr.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID has invalid type"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.URLCollection.Find(ctx, bson.M{"user_id": userID, "service": serviceName})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URLs"})
		return
	}
	defer cursor.Close(ctx)

	var urls []models.URL
	if err = cursor.All(ctx, &urls); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse URLs"})
		return
	}

	c.JSON(http.StatusOK, urls)
}

func logScanDetails(userID primitive.ObjectID, url string, scanType string, status string, message string) {
	logEntry := models.ScanLog{
		UserID:    userID,
		URL:       url,
		ScanType:  scanType,
		Status:    status,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}

	_, err := database.ScanLogCollection.InsertOne(context.Background(), logEntry)
	if err != nil {
		// Handle logging error (optional)
	}
}

func StartSpiderScan(c *gin.Context) {
	serviceName := c.Param("service")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service name is required"})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDStr.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID has invalid type"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.URLCollection.Find(ctx, bson.M{"user_id": userID, "service": serviceName})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URLs"})
		return
	}
	defer cursor.Close(ctx)

	var urls []models.URL
	if err = cursor.All(ctx, &urls); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse URLs"})
		return
	}

	// Load ZAP configuration
	cfg := config.LoadConfig()

	// Create a slice to hold the results of the scans
	var scanResults []string

	// Iterate over each URL and send a separate request to ZAP
	for _, url := range urls {
		// Split the URLList by comma to get individual URLs
		individualURLs := strings.Split(url.URLList, ",")
		for _, targetURL := range individualURLs {
			targetURL = strings.TrimSpace(targetURL) // Trim any whitespace
			zapAPIURL := fmt.Sprintf("%s/JSON/spider/action/scan?apikey=%s&url=%s", cfg.ZAPAPIURL, cfg.ZAPAPIKey, targetURL)

			// Make the request to ZAP without a body
			resp, err := http.Get(zapAPIURL) // Use GET instead of POST
			if err != nil {
				logScanDetails(userID, targetURL, "spider", "failure", err.Error())
				continue // Continue to the next URL even if one fails
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				logScanDetails(userID, targetURL, "spider", "failure", fmt.Sprintf("ZAP returned an error: %s", resp.Status))
				continue // Continue to the next URL even if one fails
			}

			logScanDetails(userID, targetURL, "spider", "success", "Spider scan started successfully")
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Spider scan initiated", "results": scanResults})
}

func StartActiveScan(c *gin.Context) {
	serviceName := c.Param("service")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service name is required"})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDStr.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID has invalid type"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.URLCollection.Find(ctx, bson.M{"user_id": userID, "service": serviceName})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URLs"})
		return
	}
	defer cursor.Close(ctx)

	var urls []models.URL
	if err = cursor.All(ctx, &urls); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse URLs"})
		return
	}

	// Load ZAP configuration
	cfg := config.LoadConfig()

	// Create a slice to hold the results of the scans
	var scanResults []string

	// Iterate over each URL and send a separate request to ZAP for active scan
	for _, url := range urls {
		// Split the URLList by comma to get individual URLs
		individualURLs := strings.Split(url.URLList, ",")
		for _, targetURL := range individualURLs {
			targetURL = strings.TrimSpace(targetURL) // Trim any whitespace
			zapAPIURL := fmt.Sprintf("%s/JSON/ascan/action/scan?apikey=%s&url=%s", cfg.ZAPAPIURL, cfg.ZAPAPIKey, targetURL)

			// Make the request to ZAP without a body
			resp, err := http.Get(zapAPIURL) // Use GET instead of POST
			if err != nil {
				logScanDetails(userID, targetURL, "active", "failure", err.Error())
				continue // Continue to the next URL even if one fails
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				logScanDetails(userID, targetURL, "active", "failure", fmt.Sprintf("ZAP returned an error: %s", resp.Status))
				continue // Continue to the next URL even if one fails
			}

			logScanDetails(userID, targetURL, "active", "success", "Active scan started successfully")
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Active scan initiated", "results": scanResults})
}

func CheckUrl(c *gin.Context) {
	serviceName := c.Param("service")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service name is required"})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDStr.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID has invalid type"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.URLCollection.Find(ctx, bson.M{"user_id": userID, "service": serviceName})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URLs"})
		return
	}
	defer cursor.Close(ctx)

	var urls []models.URL
	if err = cursor.All(ctx, &urls); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse URLs"})
		return
	}

	// Load ZAP configuration
	cfg := config.LoadConfig()

	// Initialize counts for each risk level
	alertCounts := map[string]int{
		"Medium":   0,
		"High":     0,
		"Critical": 0,
	}

	// Iterate over each URL and send a request to ZAP for alerts
	for _, url := range urls {
		individualURLs := strings.Split(url.URLList, ",")
		for _, targetURL := range individualURLs {
			targetURL = strings.TrimSpace(targetURL)
			zapAPIURL := fmt.Sprintf("%s/JSON/alert/view/alerts/?baseurl=%s&apikey=%s", cfg.ZAPAPIURL, targetURL, cfg.ZAPAPIKey)

			resp, err := http.Get(zapAPIURL)
			if err != nil {
				logScanDetails(userID, targetURL, "alert", "failure", err.Error())
				continue
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logScanDetails(userID, targetURL, "alert", "failure", err.Error())
				continue
			}

			if resp.StatusCode != http.StatusOK {
				logScanDetails(userID, targetURL, "alert", "failure", fmt.Sprintf("ZAP returned an error: %s", string(body)))
				continue
			}

			// Parse the JSON response
			var zapResponse struct {
				Alerts []struct {
					Risk string `json:"risk"`
				} `json:"alerts"`
			}
			if err := json.Unmarshal(body, &zapResponse); err != nil {
				logScanDetails(userID, targetURL, "alert", "failure", "Failed to parse JSON response")
				continue
			}

			// Count the alerts by risk level
			for _, alert := range zapResponse.Alerts {
				alertCounts[alert.Risk]++
			}

			logScanDetails(userID, targetURL, "alert", "success", string(body))
		}
	}

	// Check if the project passes or fails security
	if alertCounts["Medium"] > 0 || alertCounts["High"] > 0 || alertCounts["Critical"] > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Your project failed the security check.",
			"counts":  alertCounts,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Your project is safe. Proceed.",
			"counts":  alertCounts,
		})
	}
}

func GetAlerts(c *gin.Context) {
	serviceName := c.Param("service")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service name is required"})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDStr.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID has invalid type"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.URLCollection.Find(ctx, bson.M{"user_id": userID, "service": serviceName})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URLs"})
		return
	}
	defer cursor.Close(ctx)

	var urls []models.URL
	if err = cursor.All(ctx, &urls); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse URLs"})
		return
	}

	// Load ZAP configuration
	cfg := config.LoadConfig()

	// Initialize counts for each risk level
	alertCounts := map[string]int{
		"Informational": 0,
		"Low":           0,
		"Medium":        0,
		"High":          0,
		"Critical":      0,
	}

	// Iterate over each URL and send a separate request to ZAP for alerts
	for _, url := range urls {
		// Split the URLList by comma to get individual URLs
		individualURLs := strings.Split(url.URLList, ",")
		for _, targetURL := range individualURLs {
			targetURL = strings.TrimSpace(targetURL) // Trim any whitespace
			zapAPIURL := fmt.Sprintf("%s/JSON/alert/view/alerts/?baseurl=%s&apikey=%s", cfg.ZAPAPIURL, targetURL, cfg.ZAPAPIKey)

			// Make the request to ZAP without a body
			resp, err := http.Get(zapAPIURL) // Use GET instead of POST
			if err != nil {
				logScanDetails(userID, targetURL, "alert", "failure", err.Error())
				continue // Continue to the next URL even if one fails
			}
			defer resp.Body.Close()

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logScanDetails(userID, targetURL, "alert", "failure", err.Error())
				continue
			}

			if resp.StatusCode != http.StatusOK {
				logScanDetails(userID, targetURL, "alert", "failure", fmt.Sprintf("ZAP returned an error: %s", string(body)))
				continue // Continue to the next URL even if one fails
			}

			// Parse the JSON response
			var zapResponse struct {
				Alerts []struct {
					Risk string `json:"risk"`
				} `json:"alerts"`
			}

			if err := json.Unmarshal(body, &zapResponse); err != nil {
				logScanDetails(userID, targetURL, "alert", "failure", "Failed to parse JSON response")
				continue
			}

			// Count the alerts by risk level
			for _, alert := range zapResponse.Alerts {
				alertCounts[alert.Risk]++
			}

			logScanDetails(userID, targetURL, "alert", "success", string(body))
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alerts retrieved", "counts": alertCounts})
}

// Import API Target to ZAP
func ImportApiTarget(c *gin.Context) {
	serviceName := c.Param("service")

	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service name is required"})
		return
	}

	// Retrieve API target from MongoDB
	var apiTarget models.APITarget
	err := database.APITargetCollection.FindOne(context.Background(), bson.M{"service": serviceName}).Decode(&apiTarget)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API target not found"})
		return
	}

	// Load ZAP configuration
	cfg := config.LoadConfig()

	// Construct ZAP import URL request
	zapAPIURL := fmt.Sprintf("%s/JSON/openapi/action/importUrl/", cfg.ZAPAPIURL)
	data := fmt.Sprintf("url=%s&contextName=%s&apikey=%s", apiTarget.APISchema, serviceName, cfg.ZAPAPIKey)

	// Make the request
	resp, err := http.Post(zapAPIURL, "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import API target"})
		return
	}
	defer resp.Body.Close()

	// Read response
	body, _ := io.ReadAll(resp.Body)

	c.JSON(http.StatusOK, gin.H{"message": "API target imported", "response": string(body)})
}

// Start Active Scan on API Target
func StartApiScan(c *gin.Context) {
	serviceName := c.Param("service")

	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service name is required"})
		return
	}

	// Retrieve API target from MongoDB
	var apiTarget models.APITarget
	err := database.APITargetCollection.FindOne(context.Background(), bson.M{"service": serviceName}).Decode(&apiTarget)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API target not found"})
		return
	}

	// Load ZAP configuration
	cfg := config.LoadConfig()

	// Construct ZAP active scan request
	zapAPIURL := fmt.Sprintf("%s/JSON/ascan/action/scan/", cfg.ZAPAPIURL)
	data := fmt.Sprintf("url=%s&recurse=true&inScopeOnly=false&apikey=%s", apiTarget.APISchema, cfg.ZAPAPIKey)

	// Make the request
	resp, err := http.Post(zapAPIURL, "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start API scan"})
		return
	}
	defer resp.Body.Close()

	// Read response
	body, _ := io.ReadAll(resp.Body)

	c.JSON(http.StatusOK, gin.H{"message": "API scan started", "response": string(body)})
}

// Fetch API Alerts from ZAP
func GetApiAlerts(c *gin.Context) {
	serviceName := c.Param("service")

	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service name is required"})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDStr.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID has invalid type"})
		return
	}

	// Retrieve API target from MongoDB
	var apiTarget models.APITarget
	err := database.APITargetCollection.FindOne(context.Background(), bson.M{"user_id": userID, "service": serviceName}).Decode(&apiTarget)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API target not found"})
		return
	}

	// Load ZAP configuration
	cfg := config.LoadConfig()

	// Initialize counts for each risk level
	alertCounts := map[string]int{
		"Medium":   0,
		"High":     0,
		"Critical": 0,
	}

	// Construct ZAP API URL to fetch alerts
	zapAPIURL := fmt.Sprintf("%s/JSON/alert/view/alerts/?baseurl=%s&apikey=%s", cfg.ZAPAPIURL, apiTarget.APISchema, cfg.ZAPAPIKey)

	// Make the request to ZAP
	resp, err := http.Get(zapAPIURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch API alerts"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read API alert response"})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("ZAP returned an error: %s", string(body))})
		return
	}

	// Parse the JSON response
	var zapResponse struct {
		Alerts []struct {
			Risk string `json:"risk"`
		} `json:"alerts"`
	}
	if err := json.Unmarshal(body, &zapResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON response"})
		return
	}

	// Count the alerts by risk level
	for _, alert := range zapResponse.Alerts {
		alertCounts[alert.Risk]++
	}

	// Save alerts to MongoDB
	apiAlert := models.APIAlert{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Service:   serviceName,
		APISchema: apiTarget.APISchema,
		Alerts:    alertCounts,
		Timestamp: time.Now(),
	}

	_, err = database.APIAlertCollection.InsertOne(context.Background(), apiAlert)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save API alerts"})
		return
	}

	// Determine security check result
	securityMessage := "Your API is safe. Proceed."
	if alertCounts["Medium"] > 0 || alertCounts["High"] > 0 || alertCounts["Critical"] > 0 {
		securityMessage = "Your API failed the security check."
	}

	c.JSON(http.StatusOK, gin.H{
		"message": securityMessage,
		"counts":  alertCounts,
	})
}
