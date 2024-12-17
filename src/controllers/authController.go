package controllers

import (
	"be-kreditkita/src/config"
	"be-kreditkita/src/middlewares"
	"be-kreditkita/src/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var consumer models.Consumer

	if err := config.DB.Where("nik = ?", data["NIK"]).First(&consumer).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "NIK not found"})
	}

	if consumer.FullName != data["FullName"] {
		return c.Status(400).JSON(fiber.Map{"message": "Full Name does not match"})
	}
	if consumer.LegalName != data["LegalName"] {
		return c.Status(400).JSON(fiber.Map{"message": "Legal Name does not match"})
	}

	token, err := middlewares.GenerateJwt(strconv.Itoa(int(consumer.Id)))

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	item := map[string]string{
		"NIK":          data["NIK"],
		"FullName":     consumer.FullName,
		"LegalName":    consumer.LegalName,
		"TempatLahir":  consumer.PlaceOfBirth,
		"TanggalLahir": consumer.DateOfBirth,
		"Id":           strconv.Itoa(int(consumer.Id)),
		"Token":        token,
	}

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(item)
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	id, _ := middlewares.ParseJwt(cookie)

	var consumer models.Consumer

	config.DB.Where("id = ?", id).First(&consumer)

	return c.JSON(consumer)
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
		"Message": "Logout Success",
	})
}
