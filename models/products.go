package models

type Products struct{
	ID				uint		`json:"id"`
	Img_1			*string		`json:"img_1"`
	Img_2			*string		`json:"img_2"`
	Title			*string		`json:"title"`
	Description		*string		`json:"description"`
	Status			bool		`json:"status"`
	MRP_price		float64		`json:"mrp_price"`
	Discount_price	float64		`json:"discount_price"`
}
