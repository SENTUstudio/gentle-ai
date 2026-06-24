package etl

import (
	"strings"
	"testing"
)

func TestBuildExceptSQLBasic(t *testing.T) {
	got := BuildExceptSQL("encuestas_csi", "db_dl_dev_stg_encuestas", "db_dl_prd_stg_encuestas")
	// The query must compare dev against prd for the same target table via
	// EXCEPT, selecting all columns from both sides.
	want := "SELECT * FROM db_dl_dev_stg_encuestas.encuestas_csi EXCEPT SELECT * FROM db_dl_prd_stg_encuestas.encuestas_csi"
	if got != want {
		t.Fatalf("BuildExceptSQL() =\n %s\nwant\n %s", got, want)
	}
}

func TestBuildExceptSQLUsesBothDatabasesAndTargetTwice(t *testing.T) {
	// Triangulation with different names: the target table appears on both
	// sides of EXCEPT, each qualified by its own database.
	got := BuildExceptSQL("ventas", "db_dev", "db_prd")
	if !strings.Contains(got, "db_dev.ventas") {
		t.Fatalf("missing dev-qualified table: %q", got)
	}
	if !strings.Contains(got, "db_prd.ventas") {
		t.Fatalf("missing prd-qualified table: %q", got)
	}
	if strings.Count(got, "ventas") != 2 {
		t.Fatalf("target table must appear twice (once per side), got %d in: %q",
			strings.Count(got, "ventas"), got)
	}
}

func TestBuildExceptSQLContainsExceptKeyword(t *testing.T) {
	// The EXCEPT operator is the parity contract; assert it is present and
	// uppercase so the query is unambiguous to Athena/Spark SQL.
	got := BuildExceptSQL("t", "d", "p")
	if !strings.Contains(got, " EXCEPT ") {
		t.Fatalf("output missing ' EXCEPT ': %q", got)
	}
}

func TestBuildExceptSQLOrderIsDevFirstThenPrd(t *testing.T) {
	// The dev table MUST come first: rows present in dev but not in prd are
	// the signal of drift we care about. Assert lexical ordering.
	got := BuildExceptSQL("t", "AAADEV", "ZZZPRD")
	devIdx := strings.Index(got, "AAADEV")
	prdIdx := strings.Index(got, "ZZZPRD")
	if devIdx < 0 || prdIdx < 0 || devIdx > prdIdx {
		t.Fatalf("expected dev before prd; devIdx=%d prdIdx=%d in %q", devIdx, prdIdx, got)
	}
}
