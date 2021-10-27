package test

import (
	"compare"
	"testing"
)

func TestCompareInts(t *testing.T) {
	var testName string
	var errs []error

	testName = "42 - 42"
	t.Logf("Test %s", testName)
	errs = compare.Compare(42, 42)
	if len(errs) > 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}

	testName = "42 - 0"
	t.Logf("Test %s", testName)
	errs = compare.Compare(42, 0)
	if len(errs) != 1 {
		t.Errorf("%s: should have one error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}

	testName = "42 - &42"
	t.Logf("Test %s", testName)
	v1 := 42
	errs = compare.Compare(42, &v1)
	if len(errs) > 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}

	testName = "&42 - 42"
	t.Logf("Test %s", testName)
	errs = compare.Compare(&v1, 42)
	if len(errs) > 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}

	testName = "0 - &0"
	t.Logf("Test %s", testName)
	v2 := 0
	errs = compare.Compare(0, v2)
	if len(errs) > 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}

	testName = "&0 - 0"
	t.Logf("Test %s", testName)
	errs = compare.Compare(v2, 0)
	if len(errs) > 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}

	testName = "42 - (*int)(nil)"
	t.Logf("Test %s", testName)
	var v3 *int
	errs = compare.Compare(42, v3)
	if len(errs) != 1 {
		t.Errorf("%s: should have one error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}

	testName = "0 - (*int)(nil)"
	t.Logf("Test %s", testName)
	errs = compare.Compare(0, v3)
	if len(errs) != 1 {
		t.Errorf("%s: should have one error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}

	testName = "int8(42) - int32(42)"
	t.Logf("Test %s", testName)
	v4 := int8(42)
	v5 := int32(42)
	errs = compare.Compare(v4, v5)
	if len(errs) != 0 {
		t.Errorf("%s: should not have error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}

	testName = "uint8(42) - int32(42)"
	t.Logf("Test %s", testName)
	v6 := uint8(42)
	v7 := int32(42)
	errs = compare.Compare(v6, v7)
	if len(errs) != 0 {
		t.Errorf("%s: should not have error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}

	testName = "uint8(42) - int32(-42)"
	t.Logf("Test %s", testName)
	v8 := uint8(42)
	v9 := int32(-42)
	errs = compare.Compare(v8, v9)
	if len(errs) != 1 {
		t.Errorf("%s: should not have error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}
}
