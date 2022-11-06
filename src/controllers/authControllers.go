package controllers

import (
	"ecart/src/database"
	"ecart/src/middlewares"
	"ecart/src/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	var existing_user models.User = models.User{}

	database.DB.Where("email = ?", data["Email"]).Find(&existing_user)

	if existing_user.Id != 0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "users already exists",
		})
	}
	if !models.MatchPassword((data["Password"]), data["Password_Confirm"]) {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	user := models.User{
		FirstName: fmt.Sprint(data["FirstName"]),
		LastName:  fmt.Sprint(data["LastName"]),
		Email:     fmt.Sprint(data["Email"]),
	}
	value, _ := data["IsAdmin"].(bool)
	//you can use variable `value`
	user.IsAdmin = value

	if err := user.SetPassword(fmt.Sprint(data["Password"])); err != nil {
		return err
	}
	cart := models.Cart{
		Products: []models.Product{},
	}
	user.Cart = cart
	database.DB.Create(&user)
	return c.JSON(
		user,
	)
}
func Login(c *fiber.Ctx) error {
	var data map[string]string
	c.BodyParser(&data)
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	fName := data["FirstName"]
	password := data["Password"]
	user := models.User{
		FirstName: fName,
	}

	database.DB.Where("first_name = ?", fName).Find(&user)

	if user.Id != 0 {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			c.Status(400)
			return c.JSON(fiber.Map{
				"message": "Invalid users",
			})
		}

		payload := jwt.StandardClaims{
			Subject:   strconv.Itoa(int(user.Id)),
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		}

		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte("secret"))

		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "Invalid Creadentials",
			})
		} else {
			cookie := fiber.Cookie{
				Name:     "jwt",
				Value:    token,
				Expires:  time.Now().Add(time.Hour * 24),
				HTTPOnly: true,
			}
			c.Cookie(&cookie)
			return c.JSON(fiber.Map{
				"message": "success",
			})
		}
	}
	return c.JSON(fiber.Map{
		"message": "User Not found",
	})
}
func User(c *fiber.Ctx) error {
	id, _ := middlewares.GetUserId(c)
	var user models.User
	database.DB.Where("id = ?", id).Preload("Cart").Find(&user)

	return c.JSON(user)
}
func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "You have logged out successfully",
	})
}

func UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid users",
		})
	}
	id, _ := middlewares.GetUserId(c)

	user := models.User{
		Id:        id,
		FirstName: data["FirstName"],
		LastName:  data["LastName"],
		Email:     data["Email"],
	}
	database.DB.Model(&user).Updates(&user)
	return c.JSON(fiber.Map{
		"message": "Updated Sucessfully",
	})
}

func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid users",
		})
	}
	id, _ := middlewares.GetUserId(c)
	if !models.MatchPassword(data["Password"], data["Password_Confirm"]) {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	user := models.User{
		Id: id,
	}
	user.SetPassword(data["Password"])
	database.DB.Model(&user).Updates(&user)
	return c.JSON(fiber.Map{
		"message": "Updated Sucessfully",
	})
}
