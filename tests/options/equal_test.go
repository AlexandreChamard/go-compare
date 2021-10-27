package test

import (
	"compare"
	"testing"
)

func TestCompareBools(t *testing.T) {
	var testName string
	var errs []error

	testName = "true - true"
	t.Logf("Test %s", testName)
	errs = compare.Compare(true, true)
	if len(errs) > 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}

	testName = "true - (*bool)(nil)"
	t.Logf("Test %s", testName)
	var v3 *bool
	errs = compare.Compare(true, v3)
	if len(errs) != 1 {
		t.Errorf("%s: should have one error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}
}
