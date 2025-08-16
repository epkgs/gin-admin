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

func NewList[T any](items []T, pager *Pager) *List[T] {

	var pg Pager
	if pager != nil {
		pg = *pager
	} else {
		pg = Pager{}
	}

	count := len(items)

	if pg.Total == 0 && count > 0 {
		pg.Total = int64(count)
	}

	if pg.Page <= 0 {
		pg.Page = 1
	}

	if pg.Limit <= 0 {
		pg.Limit = 20
	}

	if pg.Limit < count {
		pg.Limit = count
	}

	return &List[T]{
		Pager: pg,
		Items: items,
	}
}

type Pager struct {
	Total int64 `json:"total,omitempty"` // total number of items
	Page  int   `json:"page,omitempty"`  // pagination index. default(1)
	Limit int   `json:"limit,omitempty"` // pagination size, less than 0 is considered as unlimited quantity. default(20)
}
