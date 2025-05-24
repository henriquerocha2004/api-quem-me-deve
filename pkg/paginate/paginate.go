package paginate

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Result struct {
	TotalRecords int `json:"total_records"`
	Data         any `json:"data"`
}

type ColumnSearch struct {
	ColumnName  string
	ColumnValue string
}

type PaginateRequest struct {
	Page          int                 `json:"page"`
	Limit         int                 `json:"limit"`
	SearchTerm    string              `json:"search_term"`
	SortField     string              `json:"sort_field"`
	SortDirection string              `json:"sort_direction"`
	ColumnSearch  []map[string]string `json:"column_search"`
}

type SearchDto struct {
	Limit         int
	TermSearch    string
	SortField     string
	SortDirection string
	ColumnSearch  []ColumnSearch
	offset        int
}

func (p *SearchDto) SetPage(page int) {
	if page == 0 {
		page = 1
	}

	p.offset = (page * p.Limit) - p.Limit
}

func (p *SearchDto) Offset() int {
	return p.offset
}

func (p *SearchDto) AddColumnSearch(columns []map[string]string) {
	for _, column := range columns {
		p.ColumnSearch = append(p.ColumnSearch, ColumnSearch{
			ColumnName:  column["name"],
			ColumnValue: column["value"],
		})
	}
}

func GetPaginateParams(r *http.Request) (*PaginateRequest, error) {
	var params PaginateRequest

	params.Page = 1
	params.Limit = 10
	params.SearchTerm = ""
	params.SortField = "created_at"
	params.SortDirection = "desc"

	if r.URL.Query().Get("page") != "" {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			log.Println("Error converting page to int:", err)
			return nil, err
		}
		params.Page = page
	}

	if r.URL.Query().Get("limit") != "" {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			log.Println("Error converting limit to int:", err)
			return nil, err
		}
		params.Limit = limit
	}

	if r.URL.Query().Get("search_term") != "" {
		params.SearchTerm = r.URL.Query().Get("search_term")
	}
	if r.URL.Query().Get("sort_field") != "" {
		params.SortField = r.URL.Query().Get("sort_field")
	}
	if r.URL.Query().Get("sort_direction") != "" {
		params.SortDirection = r.URL.Query().Get("sort_direction")
	}

	if r.URL.Query().Get("column_search[0][name]") != "" {
		for i := 0; ; i++ {
			name := r.URL.Query().Get(fmt.Sprintf("column_search[%d][name]", i))
			value := r.URL.Query().Get(fmt.Sprintf("column_search[%d][value]", i))
			if name == "" && value == "" {
				break
			}

			params.ColumnSearch = append(params.ColumnSearch, map[string]string{
				"name":  name,
				"value": value,
			})
		}
	}

	return &params, nil
}
