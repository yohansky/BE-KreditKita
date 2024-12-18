package controllers

import (
	"be-kreditkita/src/config"
	"be-kreditkita/src/helpers"
	"be-kreditkita/src/middlewares"
	"be-kreditkita/src/models"
	"be-kreditkita/src/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AllConsumers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))

	consumers := []models.Consumer{}
	result := config.DB.Find(&consumers)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch consumers",
		})
	}

	return c.JSON(models.Paginate(config.DB, &models.Consumer{}, page))
}

func CreateConsumer(c *fiber.Ctx) error {
	validTypes := []string{"image/png", "image/jpeg", "application/pdf"}
	maxSize := int64(2 << 20)

	fotoKTP, err := c.FormFile("foto_ktp")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to upload file: " + err.Error())
	}
	_, err = helpers.ValidateAndReadFile(fotoKTP, maxSize, validTypes)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	ktpURL, err := services.UploadCLoudinary(c, fotoKTP, "kredit/foto_ktp")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	fotoSelfie, err := c.FormFile("foto_selfie")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to upload file: " + err.Error())
	}
	_, err = helpers.ValidateAndReadFile(fotoSelfie, maxSize, validTypes)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	selfieURL, err := services.UploadCLoudinary(c, fotoSelfie, "kredit/foto_selfie")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
	}
	values := form.Value

	sanitizedValues := map[string]interface{}{}
	for key, val := range values {
		if len(val) > 0 {
			sanitizedValues[key] = val[0]
		}
	}

	sanitizedValues = middlewares.XSSMiddleware(sanitizedValues)

	nik, err := strconv.ParseUint(sanitizedValues["nik"].(string), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid NIK")
	}
	var existConsumer models.Consumer
	if err := config.DB.Where("nik = ?", nik).First(&existConsumer).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "NIK sudah terdaftar",
		})
	}
	salary, err := strconv.ParseFloat(sanitizedValues["salary"].(string), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid salary")
	}

	consumer := models.Consumer{
		NIK:          uint(nik),
		FullName:     sanitizedValues["full_name"].(string),
		LegalName:    sanitizedValues["legal_name"].(string),
		PlaceOfBirth: sanitizedValues["place_of_birth"].(string),
		DateOfBirth:  sanitizedValues["date_of_birth"].(string),
		Salary:       salary,
		PhotoKTP:     ktpURL.SecureURL,
		PhotoSelfie:  selfieURL.SecureURL,
	}

	if err := config.DB.Create(&consumer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create consumer")
	}

	principal := consumer.Salary * 0.3
	annualRate := 12.0
	monthlyRate := annualRate / 12 / 100
	tenors := 3

	var limits [3]float64
	tempPrincipal := principal
	for i := 0; i < tenors; i++ {
		interest := tempPrincipal * monthlyRate
		installment := (tempPrincipal / float64(tenors)) + interest
		limits[i] = installment
		tempPrincipal -= (principal / float64(tenors))
	}

	tempPrincipal = principal
	tenor6 := 6
	limitTenor6 := 0.0
	for i := 0; i < tenor6; i++ {
		interest := tempPrincipal * monthlyRate
		installment := (tempPrincipal / float64(tenor6)) + interest
		limitTenor6 += installment
		tempPrincipal -= (principal / float64(tenor6))
	}

	defaultLimit := models.Limit{
		ConsumerId:      consumer.Id,
		Tenor1:          limits[0],
		Tenor2:          limits[1],
		Tenor3:          limits[2],
		Tenor4:          limitTenor6,
		RemainingTenor1: limits[0],
		RemainingTenor2: limits[1],
		RemainingTenor3: limits[2],
		RemainingTenor4: limitTenor6,
	}

	if err := config.DB.Create(&defaultLimit).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to initialize limit",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Consumer created",
		"Consumer": consumer,
		"Limit":    defaultLimit,
	})
}

func GetConsumer(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var consumer models.Consumer

	consumer.Id = uint(id)

	result := config.DB.First(&consumer, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "ID tidak ditemukan",
		})
	}

	return c.JSON(consumer)
}

func UpdateConsumer(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var consumer models.Consumer

	consumer.Id = uint(id)

	niktostr, err := strconv.ParseUint(c.FormValue("nik"), 10, 64)
	if err != nil {
		return err
	}
	var existConsumer models.Consumer
	if err := config.DB.Where("nik = ? AND != ?", niktostr, id).First(&existConsumer).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "NIK sudah terdaftar oleh consumer lain",
		})
	}

	slrytostr, err := strconv.ParseFloat(c.FormValue("salary"), 64)
	if err != nil {
		return err
	}

	consumer.NIK = uint(niktostr)
	consumer.FullName = c.FormValue("full_name")
	consumer.LegalName = c.FormValue("legal_name")
	consumer.PlaceOfBirth = c.FormValue("place_of_birth")
	consumer.DateOfBirth = c.FormValue("date_of_birth")
	consumer.Salary = slrytostr

	fotoKTP, err := c.FormFile("foto_ktp")
	if err == nil {
		fileKtp, err := fotoKTP.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Gagal membaca file foto KTP: " + err.Error())
		}
		defer fileKtp.Close()

		fotoKTPURL, err := services.UploadCLoudinary(c, fotoKTP, "kredit/foto_ktp")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Gagal mengunggah foto KTP: " + err.Error())
		}
		consumer.PhotoKTP = fotoKTPURL.SecureURL
	}

	fotoSelfie, err := c.FormFile("foto_selfie")
	if err == nil {
		fileSelfie, err := fotoSelfie.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Gagal membaca file foto Selfie: " + err.Error())
		}
		defer fileSelfie.Close()

		fotoSelfieURL, err := services.UploadCLoudinary(c, fotoSelfie, "kredit/foto_selfie")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Gagal mengunggah foto Selfie: " + err.Error())
		}
		consumer.PhotoSelfie = fotoSelfieURL.SecureURL
	}

	config.DB.Model(&consumer).Updates(consumer)

	return c.JSON(consumer)
}

func DeleteConsumer(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	var consumer models.Consumer
	var limit models.Limit
	var transaction models.Transaction

	config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("consumer_id = ?", id).Delete(&transaction).Error; err != nil {
			return err
		}
		if err := tx.Where("consumer_id = ?", id).Delete(&limit).Error; err != nil {
			return err
		}
		if err := tx.Delete(&consumer, id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to delete consumer and related limits",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"Message": "Delete Complete",
	})
}
