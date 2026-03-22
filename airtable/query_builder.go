package airtable

import (
	"fmt"
	"net/url"
	"strings"
)


/* When is time to write, it will remove the last &
query builder append & at every iteration
*/

type queryBuilder struct {
	Builder strings.Builder
}

func (q *queryBuilder) New(baseId string, tableId string) {
	q.Builder.Reset()

	endpoint, _ := url.JoinPath(config.EndpointUrl, baseId)
	q.Builder.WriteString(endpoint)
	q.Builder.WriteString("/")
	q.Builder.WriteString(tableId)
	q.Builder.WriteString("?")
}

func (q *queryBuilder)NewWithUrl(url string) {
   q.Builder.Reset()
   q.Builder.WriteString(url)
   q.Builder.WriteString("?")
}

func (q* queryBuilder) AddRecordIds(ids ...string){
   for _, id := range ids{
      q.Builder.WriteString("records%5B%5D=")
      q.Builder.WriteString(url.QueryEscape(id))
      q.Builder.WriteString("&")
   }
}

func (q *queryBuilder) AddFields(fields []string) {
	for _, field := range fields {
		q.Builder.WriteString("fields%5B%5D=")
		q.Builder.WriteString(url.QueryEscape(field))
		q.Builder.WriteString("&")
	}
}

func (q *queryBuilder) AddOffset(offset string) {
	q.Builder.WriteString("offset=")
	q.Builder.WriteString(url.QueryEscape(offset))
	q.Builder.WriteString("&")
}

func (q *queryBuilder) AddPageSize(pageSize int) {
	q.Builder.WriteString("pageSize=")
	q.Builder.WriteString(url.QueryEscape(fmt.Sprint(pageSize)))
	q.Builder.WriteString("&")
}

func (q *queryBuilder) AddMaxRecords(maxRecords int) {
	q.Builder.WriteString("maxRecords=")
	q.Builder.WriteString(url.QueryEscape(fmt.Sprint(maxRecords)))
	q.Builder.WriteString("&")
}

func (q *queryBuilder) AddFilterByFormula(filter string) {
	q.Builder.WriteString("filterByFormula=")
	q.Builder.WriteString(url.QueryEscape(filter))
	q.Builder.WriteString("&")

}
func (q *queryBuilder) AddSort(sorts Sorts) {

	for i, sort := range sorts {

		if i > 0 {
			q.Builder.WriteString("&")
		}

		field := fmt.Sprintf("sort%%5B%d%%5D%%5Bfield%%5D=%s&", i, url.QueryEscape(sort.Field))
		q.Builder.WriteString(field)

		dir := fmt.Sprintf("sort%%5B%d%%5D%%5Bdirection%%5D=%s", i, sort.Direction)
		q.Builder.WriteString(dir)
	}

	q.Builder.WriteString("&")
}

func (q *queryBuilder) AddView(view string) {
	if view == "" {
		return
	}

	q.Builder.WriteString("view=")
	q.Builder.WriteString(url.QueryEscape(view))
	q.Builder.WriteString("&")
}

func (q *queryBuilder) AddCellFormat(format string) {
	if format == "" {
		return
	}

	q.Builder.WriteString("cellFormat=")
	q.Builder.WriteString(url.QueryEscape(format))
	q.Builder.WriteString("&")
}

func (q *queryBuilder) AddTimeZone(tz string) {
	if tz == "" {
		return
	}

	q.Builder.WriteString("timeZone=")
	q.Builder.WriteString(url.QueryEscape(tz))
	q.Builder.WriteString("&")
}

func (q *queryBuilder) AddUserLocale(locale string) {
	if locale == "" {
		return
	}

	q.Builder.WriteString("userLocale=")
	q.Builder.WriteString(url.QueryEscape(locale))
	q.Builder.WriteString("&")
}

func (q *queryBuilder) AddRecordMetadata(metadata []string) {
	for _, m := range metadata {
		q.Builder.WriteString("recordMetadata%5B%5D=")
		q.Builder.WriteString(url.QueryEscape(m))
		q.Builder.WriteString("&")
	}
}

func (q *queryBuilder) Flush() string {
	s := q.Builder.String()
	return s[:len(s)-1] // remove last &
}
