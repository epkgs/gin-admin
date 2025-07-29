package dtos

type Result[T any] struct {
	HttpStatus int    `json:"-"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
}

func NewResult[T any](code int, message string, data T) *Result[T] {
	return &Result[T]{
		HttpStatus: 200,
		Code:       code,
		Message:    message,
		Data:       data,
	}
}

func NewResultData[T any](data T) *Result[T] {
	return NewResult(0, "ok", data)
}

type List[T any] struct {
	Pager
	Items []T `json:"items"`
}

type ResultList[T any] Result[List[T]]

func NewList[T any](items []T, page, limit int, total int64) *List[T] {

	pager := Pager{
		Total: total,
		Page:  page,
		Limit: limit,
	}

	if pager.Page <= 0 {
		pager.Page = 1
	}
	if pager.Limit <= 0 {
		pager.Limit = 20
	}
	if count := len(items); pager.Limit < count {
		pager.Limit = count
	}

	return &List[T]{
		Pager: pager,
		Items: items,
	}
}

type Pager struct {
	Total int64 `json:"total,omitempty"` // total number of items
	Page  int   `json:"page,omitempty"`  // pagination index. default(1)
	Limit int   `json:"limit,omitempty"` // pagination size, less than 0 is considered as unlimited quantity. default(20)
}
