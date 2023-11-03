package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Price       int64              `json:"price"`
	Images      []string           `json:"images"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}
