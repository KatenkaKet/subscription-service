package main

import (
	"log"
	"net/http"
	"strconv"
	"subscription-service/internal/models"
	_ "subscription-service/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

//********************************************************************//
//  							 READ								  //
//********************************************************************//

// getAllRecords godoc
// @Summary Получить все записи подписок
// @Tags subscriptions, subscriptions-get
// @Produce json
// @Success 200 {array} models.Subscription
// @Failure 500 {object} map[string]string "Ошибка при получении данных"
// @Router /all [get]
func (app *application) getAllRecords(c *gin.Context) {
	rec, err := app.allModels.Subscriptions.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to return all records"})
		return
	}

	c.JSON(http.StatusOK, rec)
}

// getRecordByID godoc
// @Summary Получить запись подписки по ID
// @Tags subscriptions, subscriptions-get
// @Produce json
// @Param id path int true "ID записи"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string "Неверный формат ID"
// @Failure 500 {object} map[string]string "Ошибка при получении данных"
// @Router /{id} [get]
func (app *application) getRecordByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := app.allModels.Subscriptions.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to return record"})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// getRecordsByUserID godoc
// @Summary Получить подписки пользователя по ID
// @Tags subscriptions, subscriptions-get
// @Produce json
// @Param id path string true "UUID пользователя"
// @Success 200 {array} models.Subscription
// @Failure 400 {object} map[string]string "Неверный UUID"
// @Failure 404 {object} map[string]string "Подписки не найдены"
// @Router /user/{id} [get]
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

// getRecordsByServiceName godoc
// @Summary Получить подписки по имени сервиса
// @Tags subscriptions, subscriptions-get
// @Produce json
// @Param name path string true "Название сервиса"
// @Success 200 {array} models.Subscription
// @Failure 404 {object} map[string]string "Подписки не найдены"
// @Router /service/{name} [get]
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

// createRecord godoc
// @Summary Создать новую запись подписки
// @Tags subscriptions, subscriptions-post
// @Accept json
// @Produce json
// @Param subscription body models.MidwaySub true "Данные подписки"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Ошибка при сохранении"
// @Router /newrecord [post]
func (app *application) createRecord(c *gin.Context) {
	var mid models.MidwaySub

	if err := c.ShouldBindJSON(&mid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := mid.FromMidwaySub()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = app.allModels.Subscriptions.Insert(&sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully created record about the subscription"})
}

//********************************************************************//
//  							 DELETE								  //
//********************************************************************//

// deleteRecordByID godoc
// @Summary Удалить запись подписки по ID
// @Tags subscriptions, subscriptions-delete
// @Produce json
// @Param id path int true "ID записи"
// @Success 200 {object} map[string]string "Сообщение об успешном удалении"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 500 {object} map[string]string "Ошибка при удалении"
// @Router /{id} [delete]
func (app *application) deleteRecordByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.allModels.Subscriptions.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record about the subscription"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted record about the subscription"})
}

// deleteRecordByUserID godoc
// @Summary Удалить все записи подписок пользователя по его UUID
// @Tags subscriptions, subscriptions-delete
// @Produce json
// @Param id path string true "UUID пользователя"
// @Success 200 {object} map[string]string "Сообщение об успешном удалении"
// @Failure 400 {object} map[string]string "Неверный UUID"
// @Failure 404 {object} map[string]string "Подписки не найдены"
// @Router /user/{id} [delete]
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

// deleteRecordsByServiceName godoc
// @Summary Удалить все записи подписок по имени сервиса
// @Tags subscriptions, subscriptions-delete
// @Produce json
// @Param name path string true "Название сервиса"
// @Success 200 {object} map[string]string "Сообщение об успешном удалении"
// @Failure 404 {object} map[string]string "Подписки не найдены"
// @Router /service/{name} [delete]
func (app *application) deleteRecordsByServiceName(c *gin.Context) {
	serviseName := c.Param("name")

	err := app.allModels.Subscriptions.DeleteByServiceName(serviseName)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success deleted record about the subscription by service_name"})
}

//********************************************************************//
//  							 UPDATE								  //
//********************************************************************//

// updateRecordByID godoc
// @Summary Обновить запись подписки по ID
// @Tags subscriptions, subscriptions-put
// @Accept json
// @Produce json
// @Param id path int true "ID записи"
// @Param subscription body models.MidwaySub true "Данные для обновления"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Ошибка при обновлении"
// @Router /{id} [put]
func (app *application) updateRecordByID(c *gin.Context) {
	var sub models.Subscription

	var mid models.MidwaySub
	if err := c.ShouldBindJSON(&mid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := mid.FromMidwaySub()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sub.ID = id

	err = app.allModels.Subscriptions.Update(sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record about the subscription"})
		return
	}
	c.JSON(http.StatusOK, sub)
}

//********************************************************************//
//  							 Filter								  //
//********************************************************************//

// getRecordsByFilter godoc
// @Summary Получить сводку подписок за период
// @Tags subscriptions, subscriptions-filter
// @Produce json
// @Param from query string true "Начало периода (формат: MM-YYYY)" example:"01-2023"
// @Param to query string true "Конец периода (формат: MM-YYYY)" example:"12-2023"
// @Param user_id query string false "UUID пользователя" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param service_name query string false "Название сервиса" example:"Netflix"
// @Success 200 {object} map[string]interface{} "total_cost и список подписок"
// @Failure 400 {object} map[string]string "Неверные параметры"
// @Failure 500 {object} map[string]string "Ошибка при получении данных"
// @Router /summary [get]
func (app *application) getRecordsByFilter(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")

	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to parameters are required"})
		return
	}

	from, err := time.Parse("01-2006", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from date format. Use MM-YYYY"})
		return
	}

	to, err := time.Parse("01-2006", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to date format. Use MM-YYYY"})
		return
	}

	subscriptions, totalCost, err := app.allModels.Subscriptions.GetSummary(from, to, userID, serviceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate summary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_cost":    totalCost,
		"subscriptions": subscriptions,
	})
}
