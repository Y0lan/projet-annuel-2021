package controllers

import "github.com/gofiber/fiber/v2"

type ChallengesController interface {
	SendChallenge(ctx fiber.Ctx) error
}

func (c *authController) SendChallenge(ctx *fiber.Ctx) error {

}
