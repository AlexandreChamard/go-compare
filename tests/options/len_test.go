package test

import (
	"compare"
	"testing"
)

func TestCompareOptionLen(t *testing.T) {
	var testName string
	var errs []error

	type a1 struct {
		_1 string
	}

	testName = `" - "foo" - [len("", 3)]`
	t.Logf("Test %s", testName)
	errs = compare.Compare("", "foo", compare.CompareOptions{}.
		Field("").Len(3))
	if len(errs) > 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}
	t.Log("--------------------")

	testName = `a1{} - a1{"foo"} - [len("_1", 3)]`
	t.Logf("Test %s", testName)
	errs = compare.Compare(a1{}, a1{"foo"}, compare.CompareOptions{}.
		Field("_1").Len(3))
	if len(errs) > 0 {
		t.Errorf("%s: should have no error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}
	t.Log("--------------------")

	testName = `a1{} - b1{0, false, "bar"} - [equal("_3", "foo")]`
	t.Logf("Test %s", testName)
	errs = compare.Compare(a1{}, a1{"foobar"}, compare.CompareOptions{}.
		Field("_1").Len(3))
	if len(errs) != 3 {
		t.Errorf("%s: should have one error but got %d", testName, len(errs))
		for _, err := range errs {
			t.Error(err)
		}
	}
	t.Log("--------------------")
}
