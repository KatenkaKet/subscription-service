package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Subscription struct {
	ServiceName string
	Price       int
	UserID      string
	StartDate   string
}

// func GetAllRecords(c *gin.Context);				// Вывод всех записей
// func GetRecordsByUserID(c *gin.Context);			// Вывод записей по айди
// func GetRecordsBySubscription(c *gin.Context);	// Вывод записей по названию подписки
// func CreateRecord(c *gin.Context);				// Создание новой записи
// func UpdateRecordByUserID(c *gin.Context);		// Обновление новой записи
// func DeleteRecordByUserID(c *gin.Context);		// Удаление записи

var subscriptions = []Subscription{
	{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		StartDate:   "07-2025",
	},
	{
		ServiceName: "Spotify Premium",
		Price:       299,
		UserID:      "a12f4d3b-8c77-4b2f-9c3e-123456789abc",
		StartDate:   "06-2025",
	},
	{
		ServiceName: "Netflix",
		Price:       500,
		UserID:      "b45e6f7a-1d23-4e6f-8c9d-987654321def",
		StartDate:   "05-2025",
	},
	{
		ServiceName: "Apple Music",
		Price:       199,
		UserID:      "c78g9h0j-3k45-6l7m-8n9o-234567890ghi",
		StartDate:   "08-2025",
	},
	{
		ServiceName: "Amazon Prime",
		Price:       350,
		UserID:      "d90i1j2k-4l56-7m8n-9o0p-345678901jkl",
		StartDate:   "09-2025",
	},
}

func main() {

	router := gin.Default()

	router.GET("/pong", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.GET("/subscriptions", GetAllRecords)
	router.GET("/subscriptions/user/:id", GetRecordsByUserID)
	router.GET("/subscriptions/service/:name", GetRecordsBySubscription)
	router.PUT("/subscriptions/user/:id", UpdateRecordByUserID)
	router.DELETE("/subscriptions/user/:id", DeleteRecordByUserID)
	router.POST("/subscriptions/newrecord", CreateRecord)

	router.Run(":8080")

}

/////////////////////////////////////////////////////////////////////////
////						 	basic CRUDL						     ////
/////////////////////////////////////////////////////////////////////////

// Вывод всех записей
func GetAllRecords(c *gin.Context) {
	c.JSON(http.StatusOK, subscriptions)
}

// Вывод записей по айди
func GetRecordsByUserID(c *gin.Context) {
	id := c.Param("id")
	var sub []Subscription
	for _, v := range subscriptions {
		if v.UserID == id {
			sub = append(sub, v)
		}
	}
	if len(sub) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
	} else {
		c.JSON(http.StatusOK, sub)
	}
}

// Вывод записей по названию подписки
func GetRecordsBySubscription(c *gin.Context) {
	name := c.Param("name")
	var sub []Subscription
	for _, v := range subscriptions {
		if v.ServiceName == name {
			sub = append(sub, v)
		}
	}
	if len(sub) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
	} else {
		c.JSON(http.StatusOK, sub)
	}
}

// Обновление новой записи
func UpdateRecordByUserID(c *gin.Context) {
	id := c.Param("id")
	var input Subscription

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	for i, v := range subscriptions {
		if v.UserID == id {
			subscriptions[i] = input
			c.JSON(http.StatusOK, gin.H{
				"message": "Record updated successfully",
				"record":  input,
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

// Удаление записи
func DeleteRecordByUserID(c *gin.Context) {
	id := c.Param("id")
	for i, v := range subscriptions {
		if v.UserID == id {
			subscriptions = append(subscriptions[:i], subscriptions[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": "Record deleted successfully",
				"user_id": id,
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

//Создание новой записи
func CreateRecord(c *gin.Context){
	var newSub Subscription
	if err:= c.ShouldBindJSON(&newSub); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
        return
	}

	for _, v := range subscriptions {
		if v.UserID == newSub.UserID && v.ServiceName == newSub.ServiceName {
			c.JSON(http.StatusConflict, gin.H{"error": "Record with this UserID already exists"})
			return
		}
	}

	subscriptions = append(subscriptions, newSub)
	c.JSON(http.StatusOK, gin.H{"message": "Record added successfuly"})
}
