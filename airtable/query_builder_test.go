package airtable

import (
	"fmt"
	"net/url"
	"testing"

	. "github.com/Antfood/airgo/testutils/testutils"
)

func TestQueryBuilder(t *testing.T) {

	var qb queryBuilder
	base := "base"
	table := "table"
	fields := []string{"field1", "field2"}
	offset := "offset_id"
	pageSize := 10
   maxRecords := 100
   filter := "{foo} = bar"
   sorts := Sorts{{"field1", "asc"}, {"field2", "desc"}}
	qb.New(base, table)

	want := buildUrl(base, table)
	have := qb.Flush()
	Assert(t, have == want, "Expected '%s', got '%s'", want, have)

	qb.New(base, table)
	qb.AddFields(fields)

	want = buildUrl(base, table) + "?fields%5B%5D=" + fields[0] + "&fields%5B%5D=" + fields[1]
	have = qb.Flush()
	Assert(t, have == want, "Expected '%s', got '%s'", want, have)

	qb.New(base, table)
	qb.AddOffset(offset)

	want = buildUrl(base, table) + "?offset=" + offset
	have = qb.Flush()
	Assert(t, have == want, "Expected '%s', got '%s'", want, have)

	qb.New(base, table)
	qb.AddPageSize(pageSize)

	want = buildUrl(base, table) + "?pageSize=" + fmt.Sprintf("%d", pageSize)
   have = qb.Flush()

   Assert(t, want == have, "Expected '%s', got '%s'", want, have)

   qb.New(base, table)
   qb.AddMaxRecords(maxRecords)

   want = buildUrl(base, table) + "?maxRecords=" + fmt.Sprintf("%d", maxRecords)
   have = qb.Flush()

   Assert(t, want == have, "Expected '%s', got '%s'", want, have)

	qb.New(base, table)
	qb.AddFilterByFormula(filter)

	want = buildUrl(base, table) + "?filterByFormula=%7Bfoo%7D+%3D+bar"
	have = qb.Flush()

	Assert(t, have == want, "Expected '%s', got '%s'", want, have)

   qb.New(base, table)
   qb.AddSort(sorts)

   want = buildUrl(base, table) + "?sort%5B0%5D%5Bfield%5D=field1&sort%5B0%5D%5Bdirection%5D=asc&sort%5B1%5D%5Bfield%5D=field2&sort%5B1%5D%5Bdirection%5D=desc"
   have = qb.Flush()

   Assert(t, want == have, "Expected '%s', got '%s'", want, have)

   qb.New(base, table)
	qb.AddFields(fields)
	qb.AddOffset(offset)
	qb.AddPageSize(pageSize)
   qb.AddMaxRecords(maxRecords)
	qb.AddFilterByFormula(filter)
   qb.AddSort(sorts)

   want = buildUrl(base, table) + "?fields%5B%5D=" + fields[0] + "&fields%5B%5D=" + fields[1] + "&offset=" + offset + "&pageSize=" + fmt.Sprintf("%d", pageSize) + "&maxRecords=" + fmt.Sprintf("%d", maxRecords) + "&filterByFormula=%7Bfoo%7D+%3D+bar&sort%5B0%5D%5Bfield%5D=field1&sort%5B0%5D%5Bdirection%5D=asc&sort%5B1%5D%5Bfield%5D=field2&sort%5B1%5D%5Bdirection%5D=desc"
   have = qb.Flush()

   Assert(t, want == have, "Expected '%s', got '%s'", want, have)
}

func buildUrl(base string, table string) string {
	s, _ := url.JoinPath(baseUrl, base, table)
	return s
}
