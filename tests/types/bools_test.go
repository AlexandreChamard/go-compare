package test

import (
	"compare"
	"testing"
)

func TestCompareBools(t *testing.T) {
	var testName string
	var errs []error

	testName = `V: true - true`
	t.Logf("Test %s", testName)
	errs = compare.Compare(true, true)
	if len(errs) > 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}
	t.Log("--------------------")

	testName = `F: true - false`
	t.Logf("Test %s", testName)
	errs = compare.Compare(true, false)
	if len(errs) != 1 {
		t.Errorf("%s: should have one error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}
	t.Log("--------------------")

	testName = `V: true - &true`
	t.Logf("Test %s", testName)
	v1 := true
	errs = compare.Compare(true, &v1)
	if len(errs) > 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}
	t.Log("--------------------")

	testName = `V:&true - true`
	t.Logf("Test %s", testName)
	errs = compare.Compare(&v1, true)
	if len(errs) > 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}
	t.Log("--------------------")

	testName = `V: false - &false`
	t.Logf("Test %s", testName)
	v2 := false
	errs = compare.Compare(false, v2)
	if len(errs) > 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}
	t.Log("--------------------")

	testName = `V: &false - false`
	t.Logf("Test %s", testName)
	errs = compare.Compare(v2, false)
	if len(errs) > 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}
	t.Log("--------------------")

	testName = `F: true - (*bool)(nil)`
	t.Logf("Test %s", testName)
	var v3 *bool
	errs = compare.Compare(true, v3)
	if len(errs) != 1 {
		t.Errorf("%s: should have one error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}
	t.Log("--------------------")

	testName = `V: false - (*bool)(nil)`
	t.Logf("Test %s", testName)
	errs = compare.Compare(false, v3)
	if len(errs) > 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}
	t.Log("--------------------")
}
