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
	f             interface{}
}

type compareOptions []compareOption

type optionType int

const (
	optionType_ERROR optionType = iota
	optionType_EQUAL
	optionType_LEN
)

const fieldFormat = "^([a-zA-Z_][a-zA-Z0-9_]*\\.)*[a-zA-Z_][a-zA-Z0-9_]*$"

func optionErrorMsg(id uint, field string) string {
	return fmt.Sprintf("Option %d field `%s'", id, field)
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

func (this compareOptions) sortOptions(optionTypes []optionType) map[optionType]compareOptions {
	m := make(map[optionType]compareOptions, len(optionTypes)+1)

	m[optionType_ERROR] = compareOptions{}
	for _, optionType := range optionTypes {
		m[optionType] = compareOptions{}
	}

	for _, option := range this {
		arr, ok := m[option.optionType]
		if ok {
			m[option.optionType] = append(arr, option)
		} else {
			m[optionType_ERROR] = append(m[optionType_ERROR], option)
		}
	}
	return m
}

func (this CompareOptions) Field(fields ...string) CompareOptions {
	this.tmpFields = fields
	this.optionId++
	return this
}

type equalFunctor func(v1 reflect.Value) []error

func (this CompareOptions) Equal(values ...interface{}) CompareOptions {
	this.checkEmptyTmpFields()
	for _, field := range this.tmpFields {
		var f equalFunctor = func(v1 reflect.Value) []error {
			for _, v2 := range values {
				errs := Compare1(v1, v2)
				if errs == nil {
					return nil
				}
				fmt.Println(errs)
			}
			return []error{fmt.Errorf("%s: value does not match expected values", optionErrorMsg(this.optionId, field))}
		}
		this.options = append(this.options, compareOption{
			id:            this.optionId,
			field:         field,
			completeField: field,
			optionType:    optionType_EQUAL,
			f:             f,
		})
	}
	return this
}
