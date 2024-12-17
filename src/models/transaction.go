package models

import "gorm.io/gorm"

type Transaction struct {
	Id            uint     `json:"id"`
	NomorKontrak  string   `json:"nomor_kontrak"`
	OTR           float64  `json:"otr"`
	AdminFee      float64  `json:"admin_fee"`
	JumlahCicilan float64  `json:"jumlah_cicilan"`
	JumlahBunga   float64  `json:"jumlah_bunga"`
	NamaAsset     string   `json:"nama_asset"`
	ConsumerId    uint     `json:"consumer_id"`
	Consumer      Consumer `json:"consumer" gorm:"foreignKey:ConsumerId"`
}

func (transaction *Transaction) Count(db *gorm.DB) int64 {
	var total int64

	db.Model(&Transaction{}).Count(&total)

	return total
}

func (transaction *Transaction) Take(db *gorm.DB, limit int, offset int) interface{} {
	var transactions []Transaction

	db.Offset(offset).Limit(limit).Find(&transactions)

	return transactions
}
