package controllers

import (
	"be-kreditkita/src/config"
	"be-kreditkita/src/models"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AllTransactions(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))

	transactions := []models.Transaction{}
	result := config.DB.Find(&transactions)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch transactions",
		})
	}

	return c.JSON(models.Paginate(config.DB, &models.Transaction{}, page))
}

type TransactionRequest struct {
	ConsumerId uint
	Tenor      int
	OTR        float64
	Response   chan error
}

var transactionQueue = make(chan TransactionRequest)

func init() {
	go processTransactions()
}

func processTransactions() {
	for req := range transactionQueue {
		err := handleTransaction(req)
		req.Response <- err
	}
}

func handleTransaction(req TransactionRequest) error {
	var consumer models.Consumer
	if err := config.DB.First(&consumer, req.ConsumerId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(http.StatusNotFound, "Consumer not found")
		}
		return err
	}

	var limit models.Limit
	if err := config.DB.Where("consumer_id = ?", req.ConsumerId).First(&limit).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(http.StatusNotFound, "Limit not found for this consumer")
		}
		return err
	}

	var consumerLimit float64
	switch req.Tenor {
	case 1:
		consumerLimit = limit.RemainingTenor1
	case 2:
		consumerLimit = limit.RemainingTenor2
	case 3:
		consumerLimit = limit.RemainingTenor3
	case 6:
		consumerLimit = limit.RemainingTenor4
	default:
		return fiber.NewError(http.StatusBadRequest, "Invalid tenor value")
	}

	if req.OTR > consumerLimit {
		return fiber.NewError(http.StatusBadRequest, "Insufficient limit")
	}

	newLimit := consumerLimit - req.OTR

	switch req.Tenor {
	case 1:
		limit.RemainingTenor1 = newLimit
	case 2:
		limit.RemainingTenor2 = newLimit
	case 3:
		limit.RemainingTenor3 = newLimit
	case 6:
		limit.RemainingTenor4 = newLimit
	}

	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&limit).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func CreateTransactions(c *fiber.Ctx) error {

	type TransactionInput struct {
		ConsumerID uint    `json:"consumer_id"`
		Tenor      int     `json:"tenor"`
		OTR        float64 `json:"otr"`
		AdminFee   float64 `json:"admin_fee"`
		Interest   float64 `json:"interest"`
		AssetName  string  `json:"asset_name"`
	}

	var input TransactionInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(502).JSON(fiber.Map{
			"message": "Invalid input data",
			"error":   err.Error(),
		})
	}

	responseChan := make(chan error)

	transactionQueue <- TransactionRequest{
		ConsumerId: input.ConsumerID,
		Tenor:      input.Tenor,
		OTR:        input.OTR,
		Response:   responseChan,
	}

	if err := <-responseChan; err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Transaction failed",
			"error":   err.Error(),
		})
	}

	jumlahBunga := input.OTR * input.Interest / 100
	jumlahCicilan := (input.OTR + jumlahBunga) / float64(input.Tenor)

	transaction := models.Transaction{
		ConsumerId:    input.ConsumerID,
		NomorKontrak:  strconv.Itoa(int(input.ConsumerID)) + "-" + strconv.Itoa(int(input.Tenor)),
		OTR:           input.OTR,
		AdminFee:      input.AdminFee,
		JumlahCicilan: jumlahCicilan,
		JumlahBunga:   jumlahBunga,
		NamaAsset:     input.AssetName,
	}

	if err := config.DB.Create(&transaction).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create transaction",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Transaction created successfully",
		"transaction": transaction,
	})
}

func GetTransactions(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var transaction models.Transaction
	if err := config.DB.Preload("Consumer").First(&transaction, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Transaction not found",
		})
	}

	return c.JSON(transaction)
}

func GetTransactionsByConsumerId(c *fiber.Ctx) error {
	id := c.Params("id")

	var transactions []models.Transaction

	if err := config.DB.Where("consumer_id = ?", id).Preload("Consumer").Find(&transactions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve transactions",
		})
	}

	if len(transactions) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "No transactions found for this consumer",
		})
	}
	return c.JSON(transactions)
}
