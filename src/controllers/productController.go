package controllers

import (
	"ecart/src/database"
	"ecart/src/middlewares"
	"ecart/src/models"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Product(c *fiber.Ctx) error {
	var products []models.Product
	database.DB.Debug().Find(&products)
	return c.JSON(products)
}

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return err
	}
	database.DB.Create(&product)
	return c.JSON(product)
}

func GetProduct(c *fiber.Ctx) error {
	productId := c.AllParams()["id"]
	var product models.Product
	database.DB.Debug().Where("id = ?", productId).Find(&product)
	if product.Id == 0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}
	return c.JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if id, err := strconv.Atoi(data["Id"]); err != nil {
		return err

	} else {
		const bitSize = 64
		if price, err := strconv.ParseFloat(data["Price"], bitSize); err != nil {
			return err
		} else {
			product := models.Product{
				Id:          uint(id),
				Name:        data["Title"],
				Description: data["Description"],

				Price: price,
			}
			database.DB.Model(&product).Where("id = ?", product.Id).Updates(&product)
			return c.JSON(fiber.Map{
				"message": "Updated Sucessfully",
			})
		}
	}
}

func GetAllProductsFromCart(c *fiber.Ctx) error {
	id, _ := middlewares.GetUserId(c)
	var user models.User
	var cart models.Cart
	database.DB.Debug().Where("id = ?", id).Preload("Cart").Find(&user)
	database.DB.Debug().Where("id = ?", user.Cart.Id).Preload("Products").Find(&cart)
	return c.JSON(cart)
}
func AddProductsToCart(c *fiber.Ctx) error {
	id, _ := middlewares.GetUserId(c)
	productId, _ := strconv.Atoi(c.Query("pid"))
	var user models.User
	var product models.Product
	var cart models.Cart
	database.DB.Debug().Where("id = ?", productId).Find(&product)
	database.DB.Debug().Where("id = ?", id).Preload("Cart").Find(&user)
	database.DB.Debug().Where("id = ?", user.Cart.Id).Preload("Products").Find(&cart)
	new_product_list := append(user.Cart.Products, product)
	user.Cart.Products = new_product_list
	database.DB.Save(&user)
	return c.JSON(user)
}
func CalculateTotal(c *fiber.Ctx) error {
	id, _ := middlewares.GetUserId(c)
	var user models.User
	var cart models.Cart
	database.DB.Debug().Where("id = ?", id).Preload("Cart").Find(&user)
	database.DB.Debug().Where("id = ?", user.Cart.Id).Preload("Products").Find(&cart)
	var price float64
	for _, val := range cart.Products {
		fmt.Println(val, price)
		price = price + val.Price
	}
	return c.JSON(price)
}
