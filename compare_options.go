package compare

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

/*
fumc Model(i interface{}) {

}

créer une liste de fields temporaire sur laquelle va s'appliquer les prochains checks
Field(...string)

les différents check:
Len(int): check la longueur des champs selectionés (fail on non slice fields)
Lenf(func(interface{})int}): check la longueur des champs selectionés (fail on non slice fields)
LenRange(struct{N, Min, Max int}): check la longueur des champs selectionés (fail on non slice fields)
LenRangef(struct{N, Min, Max func(interface{})int}): check la longueur des champs selectionés (fail on non slice fields)

Equal(...interface{})
Equalf(func(interface{})[]interface{})
EqualOrNull(...interface{})
EqualOrNullf(...interface{})
EqualCase(...interface{}): ignore la case (fail on non string fields)
EqualCasef(func(interface{})[]interface{}): ignore la case (fail on non string fields)
EqualRange(struct{N, Min, Max interface{}}): check la valeur du nombre (fail on non same number field types)
EqualRangef(func(interface{})struct{N, Min, Max interface{}}): check la valeur du nombre (fail on non same number field types)
EqualAlmost(interface{}): check la valeur du float à un delta près
EqualAlmostf(func(interface{})interface{})

Regex(string)
Regexf(func(interface{})string)

Null(): (fail on non pointer or non slice fields)
NotNull(): (fail on non pointer or non slice fields)

Empty(): (fail on non slice fields)
NotEmpty(): (fail on non slice fields)

Assert(item, field interface{}) []error: do the assertion youreslf on the selected fields
AssertGroup(item, fields ...interface{}) []error: do the assertion youreslf on the selected fields

In(...interface{}): check that every slice element is in the element given (fail on non slice fields)
Inf(func(interface{}[]interface{}))
Unique(): check that every slice element is unique (fail on non slice fields)

Present(): can be combine with Null or Empty to check Present or Null or Empty
Ignore(): ignore the selected fields. Cancel every check that can be set on every selected fields or sub-fields.
IgnoreFields(): ignore les fields de ce champs mais n'ignore pas ce champ. Cancel every check that can be set on every selected sub-fields.

Rename(leftFieldName, rightFieldName string): compare the left field not to his mirror but to the right selected field

*/

type CompareOptions struct {
	tmpFields []string
	optionId  uint
	options   compareOptions
}

type compareOption struct {
	id            uint
	field         string
	completeField string
	optionType    optionType
	f             optionFunctor
}

type compareOptions []compareOption

type optionType int

const (
	optionType_ERROR optionType = iota
	optionType_EQUAL
	optionType_LEN
	optionType_LEN_RANGE
)

type valueGetter func() (reflect.Value, bool, error)
type valueGetters map[optionType]valueGetter

type optionCheckResponse struct {
	FilteredOptions map[string]compareOptions
	DoDefaultCheck  bool
}

type optionFunctor func(reflect.Value) []error

const fieldFormat = "^([a-zA-Z_][a-zA-Z0-9_]*\\.)*[a-zA-Z_][a-zA-Z0-9_]*$"

func optionErrorMsg(id uint, field string) string {
	return fmt.Sprintf("Option %d field %q", id, field)
}

// if this.tmpFields is empty, add an empty string
func (this CompareOptions) checkEmptyTmpFields() {
	if len(this.tmpFields) == 0 {
		this.tmpFields = []string{""}
	}
}

// check the validity of the fields format and remove them if malformed
func (this compareOptions) checkFieldsFormation() (errors []error) {
	if this == nil {
		return
	}
	var idxToRemove []int
	for idx, field := range this {
		if field.completeField == "" {
			continue // empty field name is valid by default
		}
		matched, err := regexp.MatchString(fieldFormat, field.completeField)
		if err != nil {
			errors = append(errors, fmt.Errorf("%s: invalid field format: %v", optionErrorMsg(field.id, field.completeField), err))
			idxToRemove = append(idxToRemove, idx)
		}
		if !matched {
			errors = append(errors, fmt.Errorf("%s: field format does not match", optionErrorMsg(field.id, field.completeField)))
			idxToRemove = append(idxToRemove, idx)
		}
	}

	// remove all invalid options
	for i := len(idxToRemove) - 1; i >= 0; i-- {
		this = append((this)[idxToRemove[i]:], (this)[:idxToRemove[i]+1]...)
	}
	return
}

func (this compareOptions) filterOptions() map[string]compareOptions {
	m := make(map[string]compareOptions)
	for _, option := range this {
		strs := strings.SplitN(option.field, ".", 1)
		var nextField string
		if len(strs) > 1 {
			nextField = strs[1]
		}

		m[strs[0]] = append(m[strs[0]], compareOption{
			id:            option.id,
			optionType:    option.optionType,
			field:         nextField,
			completeField: option.completeField,
			f:             option.f,
		})
	}
	return m
}

func (this compareOptions) execOptions(valueGetters valueGetters, hasChildren bool) (optionCheckResponse, []error) {
	var errors []error
	doDefaultCheck := true

	filteredOptions := this.filterOptions()

	if !hasChildren {
		for k, option := range filteredOptions {
			if k == "" {
				continue
			}
			for _, field := range option {
				errors = append(errors, fmt.Errorf("%s: does not exist", optionErrorMsg(field.id, field.completeField)))
			}
		}
	}

	for _, option := range filteredOptions[""] {
		getter, ok := valueGetters[option.optionType]
		if ok {
			value, doDefault, err := getter()
			if err == nil {
				doDefaultCheck = doDefaultCheck && doDefault
				errs := option.f(value)
				if len(errs) > 0 {
					errors = append(errors, errs...)
				}
			} else {
				errors = append(errors, err)
				doDefaultCheck = false
			}
		} else {
			errors = append(errors, fmt.Errorf("%s: unexpected option on this type", optionErrorMsg(option.id, option.completeField)))
		}
	}
	return optionCheckResponse{
		FilteredOptions: filteredOptions,
		DoDefaultCheck:  doDefaultCheck,
	}, errors
}

func (this CompareOptions) Field(fields ...string) CompareOptions {
	this.tmpFields = fields
	this.optionId++
	return this
}

func (this CompareOptions) Equal(values ...interface{}) CompareOptions {
	this.checkEmptyTmpFields()
	for _, field := range this.tmpFields {
		this.options = append(this.options, compareOption{
			id:            this.optionId,
			field:         field,
			completeField: field,
			optionType:    optionType_EQUAL,
			f: func(v1 reflect.Value) []error {
				for _, v2 := range values {
					errs := Compare1(v1, v2)
					if errs == nil {
						return nil
					}
				}
				return []error{fmt.Errorf("%s: value does not match expected values", optionErrorMsg(this.optionId, field))}
			},
		})
	}
	return this
}

func (this CompareOptions) Len(len int) CompareOptions {
	this.checkEmptyTmpFields()
	for _, field := range this.tmpFields {
		this.options = append(this.options, compareOption{
			id:            this.optionId,
			field:         field,
			completeField: field,
			optionType:    optionType_LEN,
			f: func(v reflect.Value) []error {
				l := v.Len()
				if l != len {
					return []error{fmt.Errorf("%s: value does not have the expected length (expected %d got %d)", optionErrorMsg(this.optionId, field), len, l)}
				}
				return nil
			},
		})
	}
	return this
}

func (this CompareOptions) LenRange(min, max int) CompareOptions {
	this.checkEmptyTmpFields()
	for _, field := range this.tmpFields {
		this.options = append(this.options, compareOption{
			id:            this.optionId,
			field:         field,
			completeField: field,
			optionType:    optionType_LEN_RANGE,
			f: func(v reflect.Value) []error {
				l := v.Len()
				if l < min || l > max {
					return []error{fmt.Errorf("%s: value does not have the expected length (expected between %d and %d got %d)", optionErrorMsg(this.optionId, field), min, max, l)}
				}
				return nil
			},
		})
	}
	return this
}
