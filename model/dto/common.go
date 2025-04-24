package dto

type PageParam[T any] struct {
	Params   T   `json:"params"`
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

type PageResult[T any] struct {
	List     []T `json:"list"`
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}
