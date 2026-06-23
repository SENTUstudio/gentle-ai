package etl

import "fmt"

// BuildExceptSQL generates the dev-vs-prd parity query used by the
// data-engineering verify path (Camino B) to detect drift between the
// development and production copies of a target table. The query returns
// the rows present in dev but absent in prd, which a human or downstream
// check interprets as pending or divergent data.
//
// The generated SQL is a two-sided EXCEPT over the fully-qualified target
// table name (<database>.<table>) on each side:
//
//	SELECT * FROM <devDB>.<target> EXCEPT SELECT * FROM <prdDB>.<target>
//
// The dev table is placed on the left because the verify contract cares
// about rows that exist in dev but not yet in prd (the direction of
// promotion). EXCEPT and SELECT are uppercase so the statement is
// unambiguous to both Athena and Spark SQL parsers.
func BuildExceptSQL(target, devDB, prdDB string) string {
	return fmt.Sprintf(
		"SELECT * FROM %s.%s EXCEPT SELECT * FROM %s.%s",
		devDB, target, prdDB, target,
	)
}
