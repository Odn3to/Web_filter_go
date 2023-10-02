package models

import "gorm.io/gorm"

type WebFilter struct {
    gorm.Model
    Data string
}

var WebFilters []WebFilter