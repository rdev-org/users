package shared

type Filter struct {
	UserId   string `json:"user_id"`
	Page     int    `json:"page" query:"page"`
	PageSize int    `json:"page_size" query:"page_size"`
	Sort     string `json:"sort" query:"sort"`
	Status   string `json:"status" query:"status"`
	Search   string `json:"search" query:"search"`
}
