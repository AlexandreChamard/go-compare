package test

import (
	"compare"
	"testing"
)

func TestCompareOptionEqual(t *testing.T) {
	var testName string
	var errs []error

	type a1 struct {
		_1 int
		_2 bool
		_3 string
	}
	type b1 struct {
		_1 int
		_2 bool
		_3 string
	}

	testName = `a1{} - b1{0, false, "foo"} - [equal("_3", "foo")]`
	t.Logf("Test %s", testName)
	options := compare.CompareOptions{}.
		Field("_3").Equal("foo")
	errs = compare.Compare(a1{}, b1{_3: "foo"}, options)
	if len(errs) > 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}
	t.Log("--------------------")

	testName = `a1{} - b1{0, false, "bar"} - [equal("_3", "foo")]`
	t.Logf("Test %s", testName)
	errs = compare.Compare(a1{}, b1{_3: "bar"}, options)
	if len(errs) != 3 {
		t.Errorf("%s: should have one error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}
	t.Log("--------------------")
}
