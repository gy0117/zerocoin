package pages

import "math"

type PageResult struct {
	Content       []any `json:"content"`       // 内容
	TotalElements int64 `json:"totalElements"` // 总数
	Number        int64 `json:"number"`        // 当前页
	TotalPages    int64 `json:"totalPages"`    // 总页数
	HasNext       bool  `json:"hasNext"`       // 是否有下一页
	IsLast        bool  `json:"isLast"`        // 是否最后
}

func New(content []any, page, pageSize, total int64) *PageResult {
	pd := &PageResult{}
	pd.Content = content
	pd.Number = page
	pd.TotalElements = total
	if pageSize <= 0 {
		pd.TotalPages = 1
	} else {
		pd.TotalPages = int64(math.Ceil(float64(total) / float64(pageSize)))
	}
	pd.HasNext = pd.Number+1 < pd.TotalPages
	pd.IsLast = !pd.HasNext
	return pd
}
