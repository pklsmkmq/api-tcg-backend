package controllers

import (
	"encoding/json"

	"go-supabase-api/config"
	"go-supabase-api/models"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	resp, _ := config.SupabaseRequest("GET", "products?select=*", nil)
	var products []models.Product
	json.NewDecoder(resp.Body).Decode(&products)

	c.JSON(200, products)
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")
	resp, _ := config.SupabaseRequest("GET", "products?id=eq."+id+"&select=*", nil)
	var products []models.Product
	json.NewDecoder(resp.Body).Decode(&products)

	if len(products) == 0 {
		c.JSON(404, gin.H{"message": "Product tidak ditemukan"})
		return
	}

	c.JSON(200, products[0])
}

func CreateProduct(c *gin.Context) {
	var p models.Product
	c.BindJSON(&p)

	resp, _ := config.SupabaseRequest("POST", "products", p)
	var result any
	json.NewDecoder(resp.Body).Decode(&result)

	c.JSON(201, result)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var p models.Product
	c.BindJSON(&p)

	resp, _ := config.SupabaseRequest("PATCH", "products?id=eq."+id, p)
	var result any
	json.NewDecoder(resp.Body).Decode(&result)

	c.JSON(200, result)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	config.SupabaseRequest("DELETE", "products?id=eq."+id, nil)
	c.JSON(200, gin.H{"message": "Product deleted"})
}
