package request

type PaginationRequest struct {
	Page     int     `json:"page"`
	Size     int     `json:"size"`
	SortBy   *string `json:"sortBy"`
	SortType *string `json:"sortType"`
}

type GetStockOfficeByVariantIDRequest struct {
	OfficeID  int64 `json:"office_id"`
	ProductID int64 `json:"product_id"`
	VariantID int64 `json:"variant_id"`
}
