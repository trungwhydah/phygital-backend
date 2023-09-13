package common

type Paging struct {
	Page  int  `form:"page"`
	Limit int  `form:"limit"`
	Total uint `form:"total"`
}

func (p *Paging) Fullfill() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 {
		p.Limit = 20
	}
}
