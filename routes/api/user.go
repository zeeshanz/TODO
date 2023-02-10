package routes

import (
	"github.com/gofiber/fiber/v2"
	apiControllers "github.com/zeeshanz/TODO/controllers/api"
)

func UserRoute(route fiber.Router) {
	route.Post("/signInUser", apiControllers.SignInUser)
	route.Post("/signUpUser", apiControllers.SignUpUser)
	route.Get("/signOutUser", apiControllers.SignOutUser)
}
