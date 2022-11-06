package routes

import (
	"ecart/src/controllers"
	"ecart/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/admin/register", controllers.Register)
	api.Post("/admin/login", controllers.Login)
	Authenticated := api.Use(middlewares.IsAuthenticated)
	Authenticated.Get("/admin/user", controllers.User)
	Authenticated.Post("/admin/logout", controllers.Logout)
	Authenticated.Patch("admin/user/update", controllers.UpdateInfo)
	Authenticated.Patch("admin/user/password", controllers.UpdatePassword)
	Authenticated.Get("products", controllers.Product)
	Authenticated.Get("product/:id", controllers.GetProduct)
	Authenticated.Get("getAllproductsFromCart", controllers.GetAllProductsFromCart)
	Authenticated.Get("addProductsToCart", controllers.AddProductsToCart)
	Authenticated.Get("calculateTotalOfCart", controllers.CalculateTotal)
	adminAuthenticated := Authenticated.Use(middlewares.IsAdmin)
	adminAuthenticated.Post("createproducts", controllers.CreateProduct)
	adminAuthenticated.Patch("updateproduct", controllers.UpdateProduct)

}
