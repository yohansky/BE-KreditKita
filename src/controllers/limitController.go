package controllers

import (
	"be-kreditkita/src/config"
	"be-kreditkita/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AllLimit(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))

	limits := []models.Limit{}
	result := config.DB.Find(&limits)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch limits",
		})
	}

	return c.JSON(models.Paginate(config.DB, &models.Limit{}, page))
}

func GetLimit(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var limit models.Limit

	result := config.DB.Preload("Consumer").First(&limit, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Limit not found",
		})
	}

	return c.JSON(limit)
}

func GetLimitByConsumerId(c *fiber.Ctx) error {
	id := c.Params("id")

	var limit models.Limit
	if err := config.DB.Where("consumer_id = ?", id).Preload("Consumer").First(&limit).Error; err != nil {
		return c.JSON(fiber.Map{"Error": "User not found"})
	}

	return c.JSON(limit)
}

func UpdateLimit(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var limit models.Limit

	limit.Id = uint(id)

	tenor1, _ := strconv.ParseFloat(c.FormValue("tenor_1"), 64)
	tenor2, _ := strconv.ParseFloat(c.FormValue("tenor_2"), 64)
	tenor3, _ := strconv.ParseFloat(c.FormValue("tenor_3"), 64)
	tenor4, _ := strconv.ParseFloat(c.FormValue("tenor_4"), 64)
	Rtenor1, _ := strconv.ParseFloat(c.FormValue("remaining_tenor_1"), 64)
	Rtenor2, _ := strconv.ParseFloat(c.FormValue("remaining_tenor_2"), 64)
	Rtenor3, _ := strconv.ParseFloat(c.FormValue("remaining_tenor_3"), 64)
	Rtenor4, _ := strconv.ParseFloat(c.FormValue("remaining_tenor_4"), 64)

	limit.Tenor1 = tenor1
	limit.Tenor2 = tenor2
	limit.Tenor3 = tenor3
	limit.Tenor4 = tenor4
	limit.RemainingTenor1 = Rtenor1
	limit.RemainingTenor2 = Rtenor2
	limit.RemainingTenor3 = Rtenor3
	limit.RemainingTenor4 = Rtenor4

	consumerIdStr := c.FormValue("consumer_id")
	consumerId, err := strconv.ParseUint(consumerIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid consumer ID",
		})
	}
	limit.ConsumerId = uint(consumerId)

	result := config.DB.Model(&limit).Updates(limit)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update limit",
		})
	}

	return c.JSON(limit)
}

type LimitUpdateRequest struct {
	Tenor   int     `json:"tenor"`
	Cicilan float64 `json:"cicilan"`
}

func UpdateLimitConsumer(c *fiber.Ctx) error {
	type PayInput struct {
		Tenor  int     `json:"tenor"`
		Amount float64 `json:"amount"`
	}

	id := c.Params("id")
	var input PayInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(502).JSON(fiber.Map{
			"message": "Invalid input data",
			"error":   err.Error(),
		})
	}

	var limit models.Limit
	if err := config.DB.Where("consumer_id = ?", id).Preload("Consumer").First(&limit).Error; err != nil {
		return c.JSON(fiber.Map{"Error": "User not found"})
	}

	var remaining float64
	var tenor float64
	var newRemaining float64

	switch input.Tenor {
	case 1:
		tenor = limit.Tenor1
		remaining = limit.RemainingTenor1
	case 2:
		tenor = limit.Tenor2
		remaining = limit.RemainingTenor2
	case 3:
		tenor = limit.Tenor3
		remaining = limit.RemainingTenor3
	case 6:
		tenor = limit.Tenor4
		remaining = limit.RemainingTenor4
	default:
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid tenor value",
		})
	}

	newRemaining = remaining + input.Amount

	switch input.Tenor {
	case 1:
		limit.RemainingTenor1 = newRemaining
	case 2:
		limit.RemainingTenor2 = newRemaining
	case 3:
		limit.RemainingTenor3 = newRemaining
	case 6:
		limit.RemainingTenor4 = newRemaining
	}

	if err := config.DB.Save(&limit).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error updating limit",
			"error":   err.Error(),
		})
	}

	if newRemaining >= tenor {
		return c.JSON(fiber.Map{
			"message": "Payment successfully processed",
			"updated_limit": fiber.Map{
				"tenor":              input.Tenor,
				"original_tenor":     tenor,
				"amount_paid":        input.Amount,
				"sisa bayar cicilan": "Lunas",
			},
		})
	}

	return c.JSON(fiber.Map{
		"message": "Payment successfully processed",
		"updated_limit": fiber.Map{
			"tenor":              input.Tenor,
			"remaining_limit":    newRemaining - tenor,
			"original_tenor":     tenor,
			"amount_paid":        input.Amount,
			"sisa bayar cicilan": tenor - newRemaining,
		},
	})
}

func DeleteLimit(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var limit models.Limit

	limit.Id = uint(id)

	result := config.DB.Delete(&limit)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete limit",
		})
	}

	return c.JSON(fiber.Map{
		"Message": "Limit deleted successfully",
	})
}
