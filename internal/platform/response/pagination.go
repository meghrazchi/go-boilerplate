package response

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

func NewPaginationMeta(page, limit int, total int64) PaginationMeta {
	var totalPages int64
	if limit > 0 {
		totalPages = (total + int64(limit) - 1) / int64(limit)
	}
	return PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
