package models

import "gorm.io/gorm"

type Limit struct {
	Id              uint     `json:"id"`
	Tenor1          float64  `json:"tenor_1"`
	Tenor2          float64  `json:"tenor_2"`
	Tenor3          float64  `json:"tenor_3"`
	Tenor4          float64  `json:"tenor_4"`
	RemainingTenor1 float64  `json:"remaining_tenor_1"`
	RemainingTenor2 float64  `json:"remaining_tenor_2"`
	RemainingTenor3 float64  `json:"remaining_tenor_3"`
	RemainingTenor4 float64  `json:"remaining_tenor_4"`
	ConsumerId      uint     `json:"consumer_id"`
	Consumer        Consumer `json:"consumer" gorm:"foreignKey:ConsumerId"`
}

func (limitt *Limit) Count(db *gorm.DB) int64 {
	var total int64

	db.Model(&Limit{}).Count(&total)

	return total
}

func (limitt *Limit) Take(db *gorm.DB, limit int, offset int) interface{} {
	var limits []Limit

	db.Offset(offset).Limit(limit).Find(&limits)

	return limits
}
