package internal

// Company represents the company entity
type Company struct {
	ID                string `json:"id" validate:"required" bson:"_id"`
	Name              string `json:"name" validate:"required,max=15"`
	Description       string `json:"description" validate:"max=3000"`
	AmountOfEmployees int    `json:"amount_of_employees" validate:"required,min=1"`
	Registered        bool   `json:"registered" validate:"required"`
	Type              string `json:"type" validate:"required,oneof=Corporations NonProfit Cooperative SoleProprietorship"`
}
