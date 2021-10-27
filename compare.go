package compare

import (
	"fmt"
	"reflect"
	"strings"
)

type compareFunctor func(reflect.Value, reflect.Value, compareOptions) []error

func Compare(i1, i2 interface{}, options ...CompareOptions) (errors []error) {
	return Compare2(reflect.ValueOf(i1), reflect.ValueOf(i2), options...)
}

func Compare1(i1 reflect.Value, i2 interface{}, options ...CompareOptions) (errors []error) {
	return Compare2(i1, reflect.ValueOf(i2), options...)
}

func Compare2(i1, i2 reflect.Value, options ...CompareOptions) (errors []error) {
	if len(options) > 1 {
		panic("Compare tooks maximum one options structure")
	}
	var optionErrors []error
	if len(options) == 1 {
		optionErrors = options[0].options.checkFieldsFormation()
		errors = compareValues(i1, i2, options[0].options)
	} else {
		errors = compareValues(i1, i2, nil)
	}
	if errors != nil {
		for i, j := 0, len(errors)-1; i < j; i, j = i+1, j-1 {
			errors[i], errors[j] = errors[j], errors[i]
		}
	}
	if len(optionErrors) > 0 {
		errors = append(optionErrors, errors...)
	}
	return errors
}

func compareValues(v1, v2 reflect.Value, options compareOptions) (errors []error) {
	// remove Ptr and Interface to get what's inside
	isElem := true
	for isElem {
		isElem = false
		switch v1.Type().Kind() {
		case reflect.Interface, reflect.Ptr:
			if v1.Elem().IsValid() {
				v1 = v1.Elem()
				isElem = true
			}
		}
	}

	isElem = true
	for isElem {
		isElem = false
		switch v2.Type().Kind() {
		case reflect.Interface, reflect.Ptr:
			if v2.Elem().IsValid() {
				v2 = v2.Elem()
				isElem = true
			}
		}
	}
	return compare(v1, v2, options)
}

func compare(v1, v2 reflect.Value, options compareOptions) (errors []error) {
	k1 := getSimplifiedType(v1.Kind())
	k2 := getSimplifiedType(v2.Kind())
	if k1 != k2 {
		// check empty values
		if (k1 == SimplifiedType_Empty && nullValue(v2)) ||
			(k2 == SimplifiedType_Empty && nullValue(v1)) {
			return
		}
		errors = append(errors, fmt.Errorf(`not the same type "%v" and "%v"`, v1.Kind(), v2.Kind()))
		return
	}
	var f compareFunctor
	switch k1 {
	case SimplifiedType_Bool:
		f = compareBools
	case SimplifiedType_Float:
		f = compareFloats
	case SimplifiedType_Int:
		f = compareInts
	case SimplifiedType_Complex:
		f = compareComplexs
	case SimplifiedType_String:
		f = compareStrings
	case SimplifiedType_Array:
		f = compareSlices
	case SimplifiedType_Map:
		f = compareMaps
	case SimplifiedType_Struct:
		f = compareStructs
	default:
		// other types are ignored
	}
	if f != nil {
		errs := f(v1, v2, options)
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}
	return
}

type SimplifiedType int

const (
	SimplifiedType_None  SimplifiedType = iota
	SimplifiedType_Empty                // Empty pointer or empty interface
	SimplifiedType_Bool
	SimplifiedType_Int
	SimplifiedType_Uint
	SimplifiedType_Float
	SimplifiedType_Complex
	SimplifiedType_String
	SimplifiedType_Array
	SimplifiedType_Map
	SimplifiedType_Struct
)

func getSimplifiedType(kind reflect.Kind) SimplifiedType {
	t, ok := map[reflect.Kind]SimplifiedType{
		reflect.Interface:  SimplifiedType_Empty,
		reflect.Ptr:        SimplifiedType_Empty,
		reflect.Bool:       SimplifiedType_Bool,
		reflect.Float32:    SimplifiedType_Float,
		reflect.Float64:    SimplifiedType_Float,
		reflect.Int:        SimplifiedType_Int,
		reflect.Int8:       SimplifiedType_Int,
		reflect.Int16:      SimplifiedType_Int,
		reflect.Int32:      SimplifiedType_Int,
		reflect.Int64:      SimplifiedType_Int,
		reflect.Uintptr:    SimplifiedType_Int,
		reflect.Uint:       SimplifiedType_Int,
		reflect.Uint8:      SimplifiedType_Int,
		reflect.Uint16:     SimplifiedType_Int,
		reflect.Uint32:     SimplifiedType_Int,
		reflect.Uint64:     SimplifiedType_Int,
		reflect.Complex64:  SimplifiedType_Complex,
		reflect.Complex128: SimplifiedType_Complex,
		reflect.String:     SimplifiedType_String,
		reflect.Array:      SimplifiedType_Array,
		reflect.Slice:      SimplifiedType_Array,
		reflect.Map:        SimplifiedType_Map,
		reflect.Struct:     SimplifiedType_Struct,
	}[kind]
	if ok {
		return t
	}
	return SimplifiedType_None
}

func getSimplifiedTypeStrict(kind reflect.Kind) SimplifiedType {
	t, ok := map[reflect.Kind]SimplifiedType{
		reflect.Interface:  SimplifiedType_Empty,
		reflect.Ptr:        SimplifiedType_Empty,
		reflect.Bool:       SimplifiedType_Bool,
		reflect.Float32:    SimplifiedType_Float,
		reflect.Float64:    SimplifiedType_Float,
		reflect.Int:        SimplifiedType_Int,
		reflect.Int8:       SimplifiedType_Int,
		reflect.Int16:      SimplifiedType_Int,
		reflect.Int32:      SimplifiedType_Int,
		reflect.Int64:      SimplifiedType_Int,
		reflect.Uintptr:    SimplifiedType_Uint,
		reflect.Uint:       SimplifiedType_Uint,
		reflect.Uint8:      SimplifiedType_Uint,
		reflect.Uint16:     SimplifiedType_Uint,
		reflect.Uint32:     SimplifiedType_Uint,
		reflect.Uint64:     SimplifiedType_Uint,
		reflect.Complex64:  SimplifiedType_Complex,
		reflect.Complex128: SimplifiedType_Complex,
		reflect.String:     SimplifiedType_String,
		reflect.Array:      SimplifiedType_Array,
		reflect.Slice:      SimplifiedType_Array,
		reflect.Map:        SimplifiedType_Map,
		reflect.Struct:     SimplifiedType_Struct,
	}[kind]
	if ok {
		return t
	}
	return SimplifiedType_None
}

func compareBools(v1, v2 reflect.Value, options compareOptions) []error {
	// TODO implements options

	if v1.Bool() != v2.Bool() {
		return []error{fmt.Errorf("expected %v got %v", v1.Bool(), v2.Bool())}
	}
	return nil
}

func compareFloats(v1, v2 reflect.Value, options compareOptions) []error {
	// TODO implements options

	if v1.Float() != v2.Float() {
		return []error{fmt.Errorf("expected %f got %f", v1.Float(), v2.Float())}
	}
	return nil
}

func compareInts(v1, v2 reflect.Value, options compareOptions) []error {
	// TODO implements options

	k1 := getSimplifiedTypeStrict(v1.Kind())
	k2 := getSimplifiedTypeStrict(v2.Kind())
	switch k1 {
	case SimplifiedType_Int:
		switch k2 {
		case SimplifiedType_Int:
			if v1.Int() != v2.Int() {
				return []error{fmt.Errorf("expected %d got %d", v1.Int(), v2.Int())}
			}
		case SimplifiedType_Uint:
			v1_ := v1.Int()
			if v1_ < 0 || uint64(v1_) != v2.Uint() {
				return []error{fmt.Errorf("expected %d got %d", v1.Int(), v2.Uint())}
			}
		}
	case SimplifiedType_Uint:
		switch k2 {
		case SimplifiedType_Int:
			v2_ := v2.Int()
			if v2_ < 0 || v1.Uint() != uint64(v2_) {
				return []error{fmt.Errorf("expected %d got %d", v1.Uint(), v2.Int())}
			}
		case SimplifiedType_Uint:
			if v1.Uint() != v2.Uint() {
				return []error{fmt.Errorf("expected %d got %d", v1.Uint(), v2.Uint())}
			}
		}
	}
	return nil
}

func compareComplexs(v1, v2 reflect.Value, options compareOptions) []error {
	if v1.Complex() != v2.Complex() {
		return []error{fmt.Errorf("expected %f+%fi got %f+%fi", real(v1.Complex()), imag(v1.Complex()), real(v2.Complex()), imag(v2.Complex()))}
	}
	return nil
}

func compareStrings(v1, v2 reflect.Value, options compareOptions) (errs []error) {
	opts := options.filterOptions()
	currentOpts, ok := opts[""]

	// TODO do this check in a global function
	if len(opts) > 1 || (len(opts) == 1 && !ok) {
		for k, opt := range opts {
			if k == "" {
				continue
			}
			for _, field := range opt {
				errs = append(errs, fmt.Errorf("%s: does not exist", optionErrorMsg(field.id, field.completeField)))
			}
		}
	}

	doDefaultCheck := true
	if ok && len(currentOpts) > 0 {
		a := currentOpts.sortOptions([]optionType{optionType_EQUAL, optionType_LEN})

		if len(a[optionType_EQUAL]) > 0 {
			doDefaultCheck = false
			for _, opt := range a[optionType_EQUAL] {
				errs_ := opt.f.(equalFunctor)(v2)
				if len(errs_) > 0 {
					errs = append(errs, errs_...)
				}
			}
		}
		// if len(a[optionType_LEN]) > 0 {
		// 	doDefaultCheck = false
		// 	for _, opt := range a[optionType_LEN] {
		// 		errs_ := opt.f.(lenFunctor)(v2)
		// 		if len(errs) > 0 {
		// 			errs = append(errs, errs_...)
		// 		}
		// 	}
		// }
	}
	if doDefaultCheck {
		if v1.String() != v2.String() {
			errs = append(errs, []error{fmt.Errorf("expected %q got %q", v1.String(), v2.String())}...)
		}
	}
	return
}

func compareSlices(v1, v2 reflect.Value, options compareOptions) (errors []error) {
	// TODO implements options
	if v1.Len() != v2.Len() {
		return []error{fmt.Errorf("Array length differ (expected %d got %d)", v1.Len(), v2.Len())}
	} else {
		for i := 0; i < v1.Len(); i++ {
			errs := compareValues(v1.Index(i), v2.Index(i), nil) // TODO implements options
			if len(errs) > 0 {
				errs[len(errs)-1] = fmt.Errorf("error in element %d: %v", i, errs[len(errs)-1])
				errors = append(errors, errs...)
			}
		}
	}
	if len(errors) > 0 {
		errors = append(errors, fmt.Errorf("there is %d errors", len(errors)))
	}
	return
}

func compareMaps(v1, v2 reflect.Value, options compareOptions) (errors []error) {
	// TODO check the missing fields found in options
	opts := options.filterOptions()

	type keyValueTuple struct {
		key   reflect.Value
		value reflect.Value
	}
	keyValueTuples := make([]keyValueTuple, 0, len(v1.MapKeys()))

	// remove tuple from slice when found
	foundTupleFromKey := func(key reflect.Value) *keyValueTuple {
		for i, tuple := range keyValueTuples {
			if compareValues(tuple.key, key, nil) == nil { // there is no options to compare key
				keyValueTuples = append(keyValueTuples[:i], keyValueTuples[i+1:]...)
				return &tuple
			}
		}
		return nil
	}

	iter := v1.MapRange()
	for iter.Next() {
		keyValueTuples = append(keyValueTuples, keyValueTuple{
			key:   iter.Key(),
			value: iter.Value(),
		})
	}
	iter = v2.MapRange()
	for iter.Next() {
		tuple := foundTupleFromKey(iter.Key())
		if tuple != nil {
			errs := compareValues(tuple.value, iter.Value(), opts[iter.Key().String()])
			if len(errs) > 0 {
				// TODO there is maybe a problem with the key printing
				errs[len(errs)-1] = fmt.Errorf("error in key %q: %v", tuple.key.String(), errs[len(errs)-1])
				errors = append(errors, errs...)
			}
		} else if !nullValue(iter.Value()) {
			errors = append(errors, fmt.Errorf("(_, %q)", iter.Key()))
		}
	}
	for _, tuple := range keyValueTuples {
		if !nullValue(tuple.value) {
			errors = append(errors, fmt.Errorf("(%q, _)", tuple.key))
		}
	}
	if len(errors) > 0 {
		errors = append(errors, fmt.Errorf("there is %d errors", len(errors)))
	}
	return
}

func compareStructs(v1, v2 reflect.Value, options compareOptions) (errors []error) {
	// TODO check the missing fields found in options
	opts := options.filterOptions()

	fieldNames := make(map[string]bool)
	// remove field from slice when found
	findField := func(field string) bool {
		if _, ok := fieldNames[field]; ok {
			delete(fieldNames, field)
			return true
		}
		return false
	}

	var errorFields []string

	t1 := v1.Type()
	for i := 0; i < t1.NumField(); i++ {
		fieldNames[t1.Field(i).Name] = true
	}

	t2 := v2.Type()
	for i := 0; i < t2.NumField(); i++ {
		fieldName := t2.Field(i).Name
		if findField(fieldName) {
			errs := compareValues(v1.FieldByName(fieldName), v2.FieldByName(fieldName), opts[fieldName])
			if len(errs) > 0 {
				errs[len(errs)-1] = fmt.Errorf("error in field %s: %v", fieldName, errs[len(errs)-1])
				errors = append(errors, errs...)
				errorFields = append(errorFields, fieldName)
			}
		} else if !nullValue(v2.FieldByName(fieldName)) {
			errors = append(errors, fmt.Errorf("(_, %s)", fieldName))
			errorFields = append(errorFields, fieldName)
		}
	}

	for fieldName := range fieldNames {
		if !nullValue(v1.FieldByName(fieldName)) {
			errors = append(errors, fmt.Errorf("(%s, _)", fieldName))
			errorFields = append(errorFields, fieldName)
		}
	}

	if len(errorFields) > 0 {
		errors = append(errors, fmt.Errorf("struct fields [%s] differ", strings.Join(errorFields, ", ")))
	}

	if len(errors) > 0 {
		errors = append(errors, fmt.Errorf("there is %d errors", len(errors)))
	}
	return
}

func nullValue(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Bool:
		return value.Bool() == false
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return value.Int() == 0
	case reflect.Uintptr,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return value.Uint() == 0
	case reflect.Complex64, reflect.Complex128:
		return value.Complex() == 0
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return value.Len() == 0
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			if !nullValue(value.Field(i)) {
				return false
			}
		}
		return true
	case reflect.Ptr:
		return value.Pointer() == 0
	case reflect.Interface:
		return true
	default:
		return false
	}
}
