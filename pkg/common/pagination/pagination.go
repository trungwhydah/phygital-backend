package pagination

import (
	"errors"
	"time"
)

const (
	PageSizeLimit  = 50
	DefaultOrderBy = "created_at"
)

var (
	ErrorInvalidLenCursor = errors.New("invalid length of cursor")
	ErrorEncode           = errors.New("encode cursor error")
)

type OrderDirectionType string

const (
	AscOrderDirection  OrderDirectionType = "asc"
	DescOrderDirection OrderDirectionType = "desc"
)

type DefaultCursor struct {
	CreatedAt *time.Time `json:"created_at"`
}

type Pagination struct {
	OrderBy        string `json:"order_by" form:"order_by" url:"order_by"`
	OrderDirection string `json:"order_direction" form:"order_direction" enums:"asc,desc" default:"desc" url:"order_direction"`
	Limit          int64  `json:"limit" form:"limit" default:"50" url:"limit"`
	Page           int64  `json:"page" form:"page" default:"1" url:"page"`
	Cursor         string `json:"cursor" form:"cursor" url:"cursor"`

	NextCursor string `json:"next_cursor,omitempty"`
	Total      int    `json:"total,omitempty"`
}

func (p *Pagination) IsAsc() bool {
	return p.OrderDirection == string(AscOrderDirection)
}

func (p *Pagination) Fulfill() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 || p.Limit > PageSizeLimit {
		p.Limit = PageSizeLimit
	}

	if p.OrderBy == "" {
		p.OrderBy = DefaultOrderBy
	}

	if p.OrderDirection == "" {
		p.OrderDirection = string(DescOrderDirection)
	}
}
