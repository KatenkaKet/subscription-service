package main

import (
	"log"
	"net/http"
	"strconv"
	"subscription-service/internal/models"

	_ "subscription-service/internal/models"

	"github.com/gin-gonic/gin"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

//func (app *application) createRecord(c *gin.Context) {
//	var sub models.Subscription
//
//	if err := c.ShouldBindJSON(&sub); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	err := app.allModels.Subscriptions.Insert(&sub)
//
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record about subscription"})
//		return
//	}
//	c.JSON(http.StatusCreated, gin.H{"data": sub})
//}

//********************************************************************//
//  							 READ								  //
//********************************************************************//

func (app *application) getAllRecords(c *gin.Context) {
	rec, err := app.allModels.Subscriptions.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to return all records"})
		return
	}

	c.JSON(http.StatusOK, rec)
}

func (app *application) getRecordByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	sub, err := app.allModels.Subscriptions.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to return record"})
	}

	c.JSON(http.StatusOK, sub)
}

func (app *application) getRecordsByUserID(c *gin.Context) {
	idStr := c.Param("id")
	//log.Println(idStr)
	var id uuid.UUID
	if err := id.Scan(idStr); err != nil {
		log.Println(id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	sub, err := app.allModels.Subscriptions.GetByUserID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sub})
}

func (app *application) getRecordsByServiceName(c *gin.Context) {
	serviceName := c.Param("name")
	log.Println(serviceName)
	sub, err := app.allModels.Subscriptions.GetByUserSubscription(serviceName)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sub})
}

//********************************************************************//
//  							 CREATE								  //
//********************************************************************//

func (app *application) createRecord(c *gin.Context) {
	var sub models.Subscription

	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// ДОБАВИТЬ ПРОВЕРКУ, ЧТОБЫ НЕ БЫЛО ДУБЛИКАТОВ

	err := app.allModels.Subscriptions.Insert(&sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record about the subscription"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully created record about the subscription"})
}

//********************************************************************//
//  							 DELETE								  //
//********************************************************************//

func (app *application) deleteRecordByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err = app.allModels.Subscriptions.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record about the subscription"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted record about the subscription"})
}

func (app *application) deleteRecordByUserID(c *gin.Context) {
	idStr := c.Param("id")
	var id uuid.UUID
	if err := id.Scan(idStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err := app.allModels.Subscriptions.DeleteByUserID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success deleted record about the subscription by user_id"})
}

func (app *application) deleteRecordsByServiceName(c *gin.Context) {
	serviseName := c.Param("name")

	err := app.allModels.Subscriptions.DeleteByServiceName(serviseName)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success deleted record about the subscription by user_id"})
}

//********************************************************************//
//  							 UPDATE								  //
//********************************************************************//

func (app *application) updateRecordByID(c *gin.Context) {
	var sub models.Subscription

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	sub.ID = id

	serviceName := c.Param("service_name")
	sub.ServiceName = serviceName

	priceStr := c.Param("price")
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	sub.Price = price

	userIdStr := c.Param("user_id")
	var userId uuid.UUID
	if err := userId.Scan(userIdStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	sub.UserID = userId

	// ДОБАВИТЬ ДАТЫ!!!

	err = app.allModels.Subscriptions.Update(sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record about the subscription"})
	}
	c.JSON(http.StatusOK, sub)
}
