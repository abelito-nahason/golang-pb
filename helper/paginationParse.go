package helper

import (
	"strconv"

	"github.com/pocketbase/pocketbase/core"
)

func GetPagesLimit(e *core.RequestEvent) (int, int) {
	var page int = 1
	var limit int = 1

	pageParse := e.Request.URL.Query().Get("page")
	page, err := strconv.Atoi(pageParse)
	if err != nil {
		page = 1
	}

	limitParse := e.Request.URL.Query().Get("limit")
	limit, err = strconv.Atoi(limitParse)
	if err != nil {
		limit = 1
	}

	return page, limit
}
