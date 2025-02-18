package controllers

import (
	"context"
	"go_back/database"
	"go_back/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type APITargetInput struct {
	Service   string `json:"service" binding:"required"`
	APISchema string `json:"api_schema" binding:"required"`
}

func CreateAPITarget(c *gin.Context) {
	var input APITargetInput
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

	apiTarget := models.APITarget{
		ID:        primitive.NewObjectID(),
		UserID:    userIDObj, // Associate the API target with the user
		Service:   input.Service,
		APISchema: input.APISchema,
	}

	_, err := database.APITargetCollection.InsertOne(context.Background(), apiTarget)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API target"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API target created", "id": apiTarget.ID})
}

func GetAPITargets(c *gin.Context) {
	var apiTargets []models.APITarget
	cursor, err := database.APITargetCollection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve API targets"})
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &apiTargets); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API targets"})
		return
	}

	c.JSON(http.StatusOK, apiTargets)
}

func UpdateAPITarget(c *gin.Context) {
	id := c.Param("id")
	var input APITargetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	_, err = database.APITargetCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": input})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update API target"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API target updated"})
}

func DeleteAPITarget(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	_, err = database.APITargetCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete API target"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API target deleted"})
}

func GetAPITargetsByService(c *gin.Context) {
	serviceName := c.Param("service")
	var apiTargets []models.APITarget
	cursor, err := database.APITargetCollection.Find(context.Background(), bson.M{"service": serviceName})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve API targets"})
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &apiTargets); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API targets"})
		return
	}

	c.JSON(http.StatusOK, apiTargets)
}
