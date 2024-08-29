package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Docket struct {
	OrderNo       string  `json:"OrderNo"`
	Customer      string  `json:"Customer"`
	PickUpPoint   string  `json:"PickUpPoint"`
	DeliveryPoint string  `json:"DeliveryPoint"`
	Quantity      int     `json:"Quantity"`
	Volume        float64 `json:"Volume"`
	Status        string  `json:"Status"`
	TruckNo       string  `json:"TruckNo"`
	LogsheetNo    string  `json:"LogsheetNo"`
}

var dockets = []Docket{}

type Logsheet struct {
	LogsheetNo string   `json:"LogsheetNo"`
	Dockets    []string `json:"Dockets"`
	TruckNo    string   `json:"TruckNo"`
}

var logsheets = []Logsheet{}

func createDocket(c *gin.Context) {
	var newDocket Docket
	if err := c.ShouldBindJSON(&newDocket); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	newDocket.OrderNo = fmt.Sprintf("TDN%04d", len(dockets)+1)
	newDocket.Status = "Created"
	dockets = append(dockets, newDocket)
	c.JSON(201, newDocket)
}

func getDocket(c *gin.Context) {
	orderNo := c.Param("orderNo")
	for _, docket := range dockets {
		if docket.OrderNo == orderNo {
			c.JSON(200, docket)
			return
		}
	}
	c.JSON(404, gin.H{"error": "Docket not found"})
}

func listDockets(c *gin.Context) {
	c.JSON(200, dockets)
}

func formatTruckNo(truckNo string) string {
	if len(truckNo) > 3 {
		return truckNo[:3] + " " + truckNo[3:]
	}
	return truckNo
}

func createLogsheet(c *gin.Context) {
	var newLogsheet Logsheet
	var updatedDockets []Docket

	if err := c.ShouldBindJSON(&newLogsheet); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	newLogsheet.LogsheetNo = fmt.Sprintf("DT%04d", len(logsheets)+1)
	formattedTruckNo := formatTruckNo(newLogsheet.TruckNo)

	for i, docket := range dockets {
		for _, docketNo := range newLogsheet.Dockets {
			if docket.OrderNo == docketNo {
				dockets[i].TruckNo = formattedTruckNo
				dockets[i].LogsheetNo = newLogsheet.LogsheetNo
				updatedDockets = append(updatedDockets, dockets[i])
			}
		}
	}

	logsheets = append(logsheets, newLogsheet)
	c.JSON(201, updatedDockets)
}

func getLogsheet(c *gin.Context) {
	logsheetNo := c.Param("logsheetNo")
	var responseDockets []Docket

	for _, logsheet := range logsheets {
		if logsheet.LogsheetNo == logsheetNo {
			for _, docketNo := range logsheet.Dockets {
				for _, docket := range dockets {
					if docket.OrderNo == docketNo {
						docket.TruckNo = formatTruckNo(docket.TruckNo)
						responseDockets = append(responseDockets, docket)
					}
				}
			}
			c.JSON(200, responseDockets)
			return
		}
	}
	c.JSON(404, gin.H{"error": "Logsheet not found"})
}

func main() {
	router := gin.Default()

	router.POST("/docket", createDocket)
	router.GET("/docket/:orderNo", getDocket)
	router.GET("/docket", listDockets)
	router.POST("/logsheet", createLogsheet)
	router.GET("/logsheet/:logsheetNo", getLogsheet)

	router.Run(":8080")
}
