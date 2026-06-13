package database

import (
	"testing"
)

func TestInstrumentedMySQlDriverReturnsDriverName(t *testing.T) {
	t.Parallel()

	driverName, err := instrumentedMySQLDriver()
	if err != nil {
		t.Fatalf("instrumentedMySQLDriver() error = %v, wnat nil", err)
	}

	if driverName == "" {
		t.Fatal("driverName = empty string, want non-empty")
	}

	if driverName == baseMySQLDriverName {
		t.Fatalf("driverName = %q, want instrumented driver name", driverName)
	}
}

func TestInstrumentedMySQLDriverCanBeCalledRepeatedly(t *testing.T) {
	t.Parallel()

	firstDriverName, err := instrumentedMySQLDriver()
	if err != nil {
		t.Fatalf("first instrumentedMySQLDriver() error = %v, want nil", err)
	}

	secondDriverName, err := instrumentedMySQLDriver()

}
