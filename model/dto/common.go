package dto

type PageParam[T any] struct {
	Params   T   `json:"params" form:"params"`
	PageNum  int `json:"pageNum" form:"pageNum"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

type PageResult[T any] struct {
	List     []T   `json:"list"`
	PageNum  int   `json:"pageNum"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}
