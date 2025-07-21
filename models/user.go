package models

import (
	"strconv"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u User) GetIdToString() string {
	return strconv.FormatUint(uint64(u.ID), 10)
}
