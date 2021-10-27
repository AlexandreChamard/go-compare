package test

import (
	"compare"
	"testing"
)

/*
test:
- same struct same values
- same struct different values
- different struct same values
- different struct different values
- different struct same values + one extra empty
- sub structures
- sub structures through interface
- sub values through interface
*/

func TestCompareStructs(t *testing.T) {
	var testName string
	var errs []error

	type a struct {
		A int
		B string
	}

	type b struct {
		A int
		B string
	}

	type b_ struct {
		A int
		B string
		c string
	}
	type c struct {
		A a
	}

	type d struct {
		A interface{}
	}

	testName = "a{42, 'foo'} - a{42, 'foo'}"
	t.Logf("Test %s", testName)
	errs = compare.Compare(a{42, "foo"}, a{42, "foo"})
	if len(errs) > 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}

	testName = "a{42, 'foo'} - a{42, 'bar'}"
	t.Logf("Test %s", testName)
	errs = compare.Compare(a{42, "foo"}, a{42, "bar"})
	if len(errs) == 0 {
		t.Errorf("%s: should have errors but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}

	testName = "a{42, 'foo'} - b{42, 'foo'}"
	t.Logf("Test %s", testName)
	errs = compare.Compare(a{42, "foo"}, b{42, "foo"})
	if len(errs) != 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}

	testName = "a{42, 'foo'} - b{42, 'bar'}"
	t.Logf("Test %s", testName)
	errs = compare.Compare(a{42, "foo"}, b{42, "bar"})
	if len(errs) == 0 {
		t.Errorf("%s: should have errors but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}

	testName = "a{42, 'foo'} - b{42, 'bar', ''}"
	t.Logf("Test %s", testName)
	errs = compare.Compare(a{42, "foo"}, b_{42, "foo", ""})
	if len(errs) != 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}

	testName = "a{42, 'foo'} - b{42, 'bar', 'AH'}"
	t.Logf("Test %s", testName)
	errs = compare.Compare(a{42, "foo"}, b_{42, "foo", "AH"})
	if len(errs) == 0 {
		t.Errorf("%s: should have errors but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}

	testName = "c{a{}} - d{nil}"
	t.Logf("Test %s", testName)
	errs = compare.Compare(c{a{}}, d{nil})
	if len(errs) != 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}
}
