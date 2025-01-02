package infragin

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/schema"
	"strconv"
)

type Page struct {
	PageNo   int `json:"page"`
	PageSize int `json:"pageSize"`
}

type PageResult[T schema.Tabler] struct {
	Page
	Records []T `json:"records"`
}

func (p *Page) Offset() int {
	var offset int
	if p.PageNo == 0 || p.PageNo == 1 {
		return 0
	}
	offset = (p.PageNo - 1) * p.PageSize
	return offset
}

func (p *Page) Limit() int {
	return p.PageSize
}

func ParsePageParams(c *gin.Context) Page {
	var page = 1
	var limit = 10
	if pageNo, success := c.GetQuery("pageNo"); success {
		if p, e := strconv.Atoi(pageNo); e == nil {
			page = p
		}
	}

	if pageSize, success := c.GetQuery("pageSize"); success {
		if p, e := strconv.Atoi(pageSize); e == nil {
			limit = p
		}
	}
	return Page{
		PageNo: page, PageSize: limit,
	}
}
