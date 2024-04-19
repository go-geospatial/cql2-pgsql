package cql2_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/go-geospatial/cql2-pgsql"
)

var _ = Describe("CQL2", func() {
	DescribeTable("predicates",
		func(cqlStr string, sql string) {
			actual, err := cql2.TranspileToSQL(cqlStr, 4326, 4326)
			Expect(err).To(BeNil())

			actual = strings.TrimSpace(actual)
			Expect(actual).To(Equal(sql))
		},
		Entry("between", "1990-01-01 BETWEEN time_start AND time_end", "timestamp '1990-01-01' BETWEEN \"time_start\" AND \"time_end\""),
		Entry("empty", "", ""),
		Entry("greater than property", "id > tt", "\"id\" > \"tt\""),
		Entry("greater than constant", "id > 1", "\"id\" > 1"),
		Entry("greater than or equal", "id >= 1", "\"id\" >= 1"),
		Entry("less than", "id < 1", "\"id\" < 1"),
		Entry("less than or equal", "id <= 1", "\"id\" <= 1"),
		Entry("equal", "id = 1", "\"id\" = 1"),
		Entry("not equal", "id <> 1", "\"id\" <> 1"),
		Entry("equal negative number", "id = -1.2345", "\"id\" = -1.2345"),
		Entry("equal property", "id = id2", "\"id\" = \"id2\""),
		Entry("equal string", "id = 'foo'", "\"id\" = 'foo'"),
		Entry("like string", "id LIKE 'foo'", "\"id\" LIKE 'foo'"),
		Entry("case-insensitive like string", "id ILIKE 'foo'", "\"id\" ILIKE 'foo'"),
		Entry("case-insensitvie like wildcard", "id ILIKE '%Ca%'", "\"id\" ILIKE '%Ca%'"),
		Entry("property bwetween constants", "id BETWEEN 1 and 2", "\"id\" BETWEEN 1 AND 2"),
		Entry("property not between constants", "id NOT BETWEEN 1 and 2", "\"id\" NOT BETWEEN 1 AND 2"),
		Entry("in list", "id IN (1,2,3)", "\"id\" IN (1,2,3)"),
		Entry("not in list", "id NOT IN (1,2,3)", "\"id\" NOT IN (1,2,3)"),
		Entry("in list of strings", "id IN ('a','b','c')", "\"id\" IN ('a','b','c')"),
		Entry("is null", "id IS NULL", "\"id\" IS NULL"),
		Entry("is not null", "id IS NOT NULL", "\"id\" IS NOT NULL"),
		Entry("spatial crosses function", "crosses(geom, POINT(0 0))", "ST_Crosses(\"geom\",'SRID=4326;POINT(0 0)'::geometry)"),
		Entry("spatial contains function", "Contains(geom, POINT(0 0))", "ST_Contains(\"geom\",'SRID=4326;POINT(0 0)'::geometry)"),
		Entry("spatial disjoint function", "DISJOINT(geom, POINT(0 0))", "ST_Disjoint(\"geom\",'SRID=4326;POINT(0 0)'::geometry)"),
		Entry("spatial equals function", "EQUALS(geom, POINT(0 0))", "ST_Equals(\"geom\",'SRID=4326;POINT(0 0)'::geometry)"),
		Entry("spatial intersects function", "INTERSECTS(geom, POINT(0 0))", "ST_Intersects(\"geom\",'SRID=4326;POINT(0 0)'::geometry)"),
		Entry("spatial overlaps function", "OVERLAPS(geom, POINT(0 0))", "ST_Overlaps(\"geom\",'SRID=4326;POINT(0 0)'::geometry)"),
		Entry("spatial touches function", "TOUCHES(geom, POINT(0 0))", "ST_Touches(\"geom\",'SRID=4326;POINT(0 0)'::geometry)"),
		Entry("spatial within function", "within(geom, POINT(0 0))", "ST_Within(\"geom\",'SRID=4326;POINT(0 0)'::geometry)"),
		Entry("spatial dwithin function", "Dwithin(geom, POINT(0 0), 100)", "ST_DWithin(\"geom\",'SRID=4326;POINT(0 0)'::geometry,100)"),
		Entry("greater than 1 and less than 9", "x > 1 AND x < 9", "\"x\" > 1 AND \"x\" < 9"),
		Entry("equal 1 or 2", "x = 1 OR x = 2", "\"x\" = 1 OR \"x\" = 2"),
		Entry("x = 1 or 2 and y less than 4", "(x = 1 OR x = 2) AND y < 4", "(\"x\" = 1 OR \"x\" = 2) AND \"y\" < 4"),
		Entry("x equal 1 or (x = 2 and y < 4)", "x = 1 OR (x = 2 AND y < 4)", "\"x\" = 1 OR (\"x\" = 2 AND \"y\" < 4)"),
		Entry("multiple boolean expressions", "x = 1 AND y = 2 AND z = 3 OR a = 4", "\"x\" = 1 AND \"y\" = 2 AND \"z\" = 3 OR \"a\" = 4"),
		Entry("not is not null", "NOT x IS NOT NULL", "NOT  \"x\" IS NOT NULL"),
		Entry("not boolean", "NOT TRUE OR FALSE", "NOT TRUE OR FALSE"),
		Entry("not boolean lower-case", "NOT true OR false", "NOT true OR false"),
		Entry("complex expression", "x = 1 OR NOT (x = 2 AND y < 4)", "\"x\" = 1 OR NOT (\"x\" = 2 AND \"y\" < 4)"),
	)

	DescribeTable("arithematic",
		func(cqlStr string, sql string) {
			actual, err := cql2.TranspileToSQL(cqlStr, 4326, 4326)
			Expect(err).To(BeNil())

			actual = strings.TrimSpace(actual)
			Expect(actual).To(Equal(sql))
		},
		Entry("addition", "p > 1 + x", "\"p\" > 1 + \"x\""),
		Entry("addition and multiply", "p > 2 * 3 + x", "\"p\" > 2 * 3 + \"x\""),
		Entry("order of operations", "p > 2 * (3 + x)", "\"p\" > 2 * (3 + \"x\")"),
		Entry("division order of operations", "p > (y + 5) / (3 - x)", "\"p\" > (\"y\" + 5) / (3 - \"x\")"),
		Entry("modulus", "p = x % 10", "\"p\" = \"x\" % 10"),
		Entry("power", "p = x ^ (i + 2)", "\"p\" = \"x\" ^ (\"i\" + 2)"),
		Entry("between arithematic", "p BETWEEN x + 10 AND x * 2", "\"p\" BETWEEN \"x\" + 10 AND \"x\" * 2"),
		Entry("between arithematic expression and constant", "p BETWEEN 2 * (1 + 1000000) AND 900000", "\"p\" BETWEEN 2 * (1 + 1000000) AND 900000"),
		Entry("xor", "p = 'a' || x || 'b'", "\"p\" = 'a' || \"x\" || 'b'"),
	)

	DescribeTable("arithematic",
		func(cqlStr string, sql string) {
			actual, err := cql2.TranspileToSQL(cqlStr, 4326, 4326)
			Expect(err).To(BeNil())

			actual = strings.TrimSpace(actual)
			Expect(actual).To(Equal(sql))
		},
		Entry("exponential notation", "p > 1.0E+1", "\"p\" > 1.0E+1"),
		Entry("exponential notation lower-case", "p > 1.0e+1", "\"p\" > 1.0e+1"),
		Entry("equals point", "equals(geom, POINT(0 0))",
			"ST_Equals(\"geom\",'SRID=4326;POINT(0 0)'::geometry)"),
		Entry("equals linestring", "equals(geom, LINESTRING(0 0, 1 1))",
			"ST_Equals(\"geom\",'SRID=4326;LINESTRING(0 0,1 1)'::geometry)"),
		Entry("equals polygon", "equals(geom, POLYGON((0 0, 0 9, 9 0, 0 0)))",
			"ST_Equals(\"geom\",'SRID=4326;POLYGON((0 0,0 9,9 0,0 0))'::geometry)"),
		Entry("equals polygon", "equals(geom, POLYGON((0 0, 0 9, 9 0, 0 0),(1 1, 1 8, 8 1, 1 1)))",
			"ST_Equals(\"geom\",'SRID=4326;POLYGON((0 0,0 9,9 0,0 0),(1 1,1 8,8 1,1 1))'::geometry)"),
		Entry("equals multipoint", "equals(geom, MULTIPOINT((0 0), (0 9)))",
			"ST_Equals(\"geom\",'SRID=4326;MULTIPOINT((0 0),(0 9))'::geometry)"),
		Entry("equals multilinestring", "equals(geom, MULTILINESTRING((0 0, 1 1),(1 1, 2 2)))",
			"ST_Equals(\"geom\",'SRID=4326;MULTILINESTRING((0 0,1 1),(1 1,2 2))'::geometry)"),
		Entry("equals multipolygon", "equals(geom, MULTIPOLYGON(((1 4, 4 1, 1 1, 1 4)), ((1 9, 4 9, 1 6, 1 9))))",
			"ST_Equals(\"geom\",'SRID=4326;MULTIPOLYGON(((1 4,4 1,1 1,1 4)),((1 9,4 9,1 6,1 9)))'::geometry)"),
		Entry("equals geometrycollection", "equals(geom, GEOMETRYCOLLECTION(POLYGON((1 4, 4 1, 1 1, 1 4)),LINESTRING (3 3, 5 5), POINT (1 5)))",
			"ST_Equals(\"geom\",'SRID=4326;GEOMETRYCOLLECTION(POLYGON((1 4,4 1,1 1,1 4)),LINESTRING(3 3,5 5),POINT(1 5))'::geometry)"),
		Entry("equals envelope", "equals(geom, ENVELOPE(1,2,3,4))",
			"ST_Equals(\"geom\",ST_MakeEnvelope(1,2,3,4,4326))"),
	)

	DescribeTable("geometries with SRID",
		func(cqlStr string, filterSRID int, sourceSRID int, sql string) {
			actual, err := cql2.TranspileToSQL(cqlStr, filterSRID, sourceSRID)
			Expect(err).To(BeNil())
			actual = strings.TrimSpace(actual)
			Expect(actual).To(Equal(sql))
		},

		Entry("srid on point", "equals(geom, POINT(0 0))", 1111, 2222,
			"ST_Equals(\"geom\",ST_Transform('SRID=1111;POINT(0 0)'::geometry,2222))"),
		Entry("srid on envelope", "equals(geom, ENVELOPE(1,2,3,4))", 1111, 2222,
			"ST_Equals(\"geom\",ST_Transform(ST_MakeEnvelope(1,2,3,4,1111),2222))"),
	)

	DescribeTable("throws syntax errors",
		func(cqlStr string) {
			_, err := cql2.TranspileToSQL(cqlStr, 4326, 4326)
			Expect(err).ToNot(BeNil())
		},
		Entry("no operator", "x y"),
		Entry("double equal", "x == y"),
		Entry("mix constant and property", "x > 10y"),
		Entry("invalid expression", "NOT x IS > 3"),
		Entry("extra paren", "equals(geom, ENVELOPE(1,2,3,4)))"),
		Entry("comma between ordinates", "equals(geom, POINT(0,0))"),
		Entry("bad temporal value year", "p > 200-01"),
		Entry("bad temporal value no day", "p > 2000-01"),
		Entry("bad temporal values time missing minutes and seconds", "p > 2000-01-01T01"),
	)
})
