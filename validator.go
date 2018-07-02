package validator

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

const tagName string = "valid"

type Validator struct {
	Translator    func(structName string, tag string, attribute string, params map[string]string) error
	Attributes    map[string]string
	CustomMessage map[string]string
}

var loadValidatorOnce *Validator
var once sync.Once

// New returns a new instance of 'valid' with sane defaults.
func New() *Validator {
	once.Do(func() {
		loadValidatorOnce = &Validator{}
	})
	return loadValidatorOnce
}

// Struct as gin binding function needed
func (v *Validator) Struct(s interface{}) error {
	_, err := ValidateStruct(s)
	return err
}

// IsRequiredIf check value required when anotherfield str is a member of the set of strings params
func IsRequiredIf(v reflect.Value, anotherfield reflect.Value, params ...string) bool {
	switch anotherfield.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:

		if IsIn(anotherfield.String(), params...) {
			return isEmptyValue(v)
		}

		return false
	}

	return false
}

// IsIn check if string str is a member of the set of strings params
func IsIn(str string, params ...string) bool {
	for _, param := range params {
		if str == param {
			return true
		}
	}

	return false
}

// IsEmail check if the string is an email.
func IsEmail(str string) bool {
	// TODO uppercase letters are not supported
	return rxEmail.MatchString(str)
}

/*
// Between check string's length
func Between(str string, params ...string) bool {

	if len(params) == 2 {
		value, _ := ToFloat(str)
		min, _ := ToFloat(params[0])
		max, _ := ToFloat(params[1])
		return DigitsBetween(value, min, max)
	}

	return false
}
*/

// Between check The field under validation must have a size between the given min and max. Strings, numerics, arrays, and files are evaluated in the same fashion as the size rule.
func Between(v reflect.Value, params ...string) bool {
	if len(params) != 2 {
		return false
	}

	switch v.Kind() {
	case reflect.String:
		min, _ := ToInt(params[0])
		max, _ := ToInt(params[1])
		return BetweenString(v.String(), min, max)
	case reflect.Slice, reflect.Map, reflect.Array:
		min, _ := ToInt(params[0])
		max, _ := ToInt(params[1])
		return DigitsBetweenInt64(int64(v.Len()), min, max)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		min, _ := ToInt(params[0])
		max, _ := ToInt(params[1])
		return DigitsBetweenInt64(v.Int(), min, max)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		min, _ := ToUint(params[0])
		max, _ := ToUint(params[1])
		return DigitsBetweenUint64(v.Uint(), min, max)
	case reflect.Float32, reflect.Float64:
		min, _ := ToFloat(params[0])
		max, _ := ToFloat(params[1])
		return DigitsBetweenFloat64(v.Float(), min, max)
	}

	panic(fmt.Sprintf("validator: Between unsupport Type %T", v.Interface()))
}

// IsNumeric check if the string contains only numbers. Empty string is valid.
func IsNumeric(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxNumeric.MatchString(str)
}

// IsInt check if the string is an integer. Empty string is valid.
func IsInt(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxInt.MatchString(str)
}

// IsFloat check if the string is a float.
func IsFloat(str string) bool {
	return str != "" && rxFloat.MatchString(str)
}

// IsNull check if the string is null.
func IsNull(str string) bool {
	return len(str) == 0
}

// Max is the validation function for validating if the current field's value is less than or equal to the param's value.
func Max(v reflect.Value, param string) bool {
	return IsLte(v, param)
}

// Min is the validation function for validating if the current field's value is greater than or equal to the param's value.
func Min(v reflect.Value, param string) bool {
	return IsGte(v, param)
}

// IsLt is the validation function for validating if the current field's value is less than the param's value.
func IsLt(v reflect.Value, param string) bool {
	switch v.Kind() {
	case reflect.String:
		p, _ := ToInt(param)
		return IsLtString(v.String(), p)
	case reflect.Slice, reflect.Map, reflect.Array:
		p, _ := ToInt(param)
		return int64(v.Len()) < p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, _ := ToInt(param)
		return IsLtInt64(v.Int(), p)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, _ := ToUint(param)
		return IsLtUnit64(v.Uint(), p)
	case reflect.Float32, reflect.Float64:
		p, _ := ToFloat(param)
		return IsLtFloat64(v.Float(), p)
	}

	panic(fmt.Sprintf("validator: IsLt unsupport Type %T", v.Interface()))
}

// IsLte is the validation function for validating if the current field's value is less than or equal to the param's value.
func IsLte(v reflect.Value, param string) bool {
	switch v.Kind() {
	case reflect.String:
		p, _ := ToInt(param)
		return IsLteString(v.String(), p)
	case reflect.Slice, reflect.Map, reflect.Array:
		p, _ := ToInt(param)
		return int64(v.Len()) <= p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, _ := ToInt(param)
		return IsLteInt64(v.Int(), p)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, _ := ToUint(param)
		return IsLteUnit64(v.Uint(), p)
	case reflect.Float32, reflect.Float64:
		p, _ := ToFloat(param)
		return IsLteFloat64(v.Float(), p)
	}

	panic(fmt.Sprintf("validator: IsLte unsupport Type %T", v.Interface()))
}

// IsGte is the validation function for validating if the current field's value is greater than or equal to the param's value.
func IsGte(v reflect.Value, param string) bool {
	switch v.Kind() {
	case reflect.String:
		p, _ := ToInt(param)
		return IsGteString(v.String(), p)
	case reflect.Slice, reflect.Map, reflect.Array:
		p, _ := ToInt(param)
		return int64(v.Len()) >= p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, _ := ToInt(param)
		return IsGteInt64(v.Int(), p)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, _ := ToUint(param)
		return IsGteUnit64(v.Uint(), p)
	case reflect.Float32, reflect.Float64:
		p, _ := ToFloat(param)
		return IsGteFloat64(v.Float(), p)
	}

	panic(fmt.Sprintf("validator: IsGte unsupport Type %T", v.Interface()))
}

// IsGt is the validation function for validating if the current field's value is greater than to the param's value.
func IsGt(v reflect.Value, param string) bool {
	switch v.Kind() {
	case reflect.String:
		p, _ := ToInt(param)
		return IsGtString(v.String(), p)
	case reflect.Slice, reflect.Map, reflect.Array:
		p, _ := ToInt(param)
		return int64(v.Len()) > p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, _ := ToInt(param)
		return IsGtInt64(v.Int(), p)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, _ := ToUint(param)
		return IsGtUnit64(v.Uint(), p)
	case reflect.Float32, reflect.Float64:
		p, _ := ToFloat(param)
		return IsGtFloat64(v.Float(), p)
	}

	panic(fmt.Sprintf("validator: IsGt unsupport Type %T", v.Interface()))
}

// ValidateStruct use tags for fields.
// result will be equal to `false` if there are any errors.
func ValidateStruct(s interface{}) (bool, error) {
	if s == nil {
		return false, nil
	}

	//var err error

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	// we only accept structs
	if val.Kind() != reflect.Struct {
		return false, fmt.Errorf("function only accepts structs; got %s", val.Kind())
	}

	var errs Errors
	result := true
	fields := cachedTypeFields(val.Type())

	for _, f := range fields {
		valueField := val.Field(f.index[0])
		isValid, err := valid(valueField, &f, val)
		if err != nil {
			errs = append(errs, err)
			fmt.Println(err)
		}
		result = result && isValid
	}

	return result, errs
}

func valid(v reflect.Value, f *field, o reflect.Value) (isValid bool, resultErr error) {
	if !v.IsValid() {
		return false, nil
	}

	for _, tag := range f.validTags {
		if tag.name == "required" && isEmptyValue(v) {
			return false, &Error{
				Name:       f.name,
				StructName: f.structName,
				Attribute:  f.attribute,
				Err:        formatsMessages(f.structName, tag.name, f.attribute, nil),
				Tag:        tag.name,
			}
		}

		if tag.name == "requiredIf" && len(tag.params) >= 2 && IsRequiredIf(v, o.FieldByName(tag.params[0]), tag.params[1:]...) {
			return false, &Error{
				Name:       f.name,
				StructName: f.structName,
				Attribute:  f.attribute,
				Err: formatsMessages(f.structName, tag.name, f.attribute, map[string]string{
					"other": tag.params[0],
					"value": strings.Join(tag.params[1:], ","),
				}),
				Tag: tag.name,
			}
		}
	}

	switch v.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:

		for _, tag := range f.validTags {

			if validfunc, ok := RuleMap[tag.name]; ok {
				isValid := validfunc(v)
				if !isValid {
					return false, &Error{
						Name:       f.name,
						StructName: f.structName,
						Attribute:  f.attribute,
						Err:        formatsMessages(f.structName, tag.name, f.attribute, nil),
						Tag:        tag.name,
					}
				}
			}

			if validfunc, ok := ParamRuleMap[tag.name]; ok {
				isValid := validfunc(v, tag.params...)
				if !isValid {
					return false, &Error{
						Name:       f.name,
						StructName: f.structName,
						Attribute:  f.attribute,
						Err:        formatsMessages(f.structName, tag.name, f.attribute, nil),
						Tag:        tag.name,
					}
				}
			}

			switch v.Kind() {
			case reflect.String:
				if validfunc, ok := StringRulesMap[tag.name]; ok {
					isValid := validfunc(v.String())
					if !isValid {
						return false, Error{
							Name:       f.name,
							StructName: f.structName,
							Attribute:  f.attribute,
							Err:        formatsMessages(f.structName, tag.name, f.attribute, nil),
							Tag:        tag.name,
						}
					}
				}
			}
		}
		return true, nil
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return false, &UnsupportedTypeError{v.Type()}
		}

		for _, tag := range f.validTags {
			if validfunc, ok := ParamRuleMap[tag.name]; ok {
				isValid := validfunc(v, tag.params...)
				if !isValid {
					return false, &Error{
						Name: f.name,
						Tag:  tag.name,
					}
				}
			}
		}

		var sv stringValues
		sv = v.MapKeys()
		sort.Sort(sv)
		result := true
		for _, k := range sv {
			var resultItem bool
			var err error
			if v.MapIndex(k).Kind() != reflect.Struct {
				resultItem, err = valid(v.MapIndex(k), f, o)
				if err != nil {
					return false, err
				}
			} else {
				resultItem, err = ValidateStruct(v.MapIndex(k).Interface())
				if err != nil {
					return false, err
				}
			}
			result = result && resultItem
		}
		return result, nil
	case reflect.Slice, reflect.Array:
		for _, tag := range f.validTags {
			if validfunc, ok := ParamRuleMap[tag.name]; ok {
				isValid := validfunc(v, tag.params...)
				if !isValid {
					return false, &Error{
						Name: f.name,
						Tag:  tag.name,
					}
				}
			}
		}
		result := true
		for i := 0; i < v.Len(); i++ {
			var resultItem bool
			var err error
			if v.Index(i).Kind() != reflect.Struct {
				resultItem, err = valid(v.Index(i), f, o)
				if err != nil {
					return false, err
				}
			} else {
				resultItem, err = ValidateStruct(v.Index(i).Interface())
				if err != nil {
					return false, err
				}
			}
			result = result && resultItem
		}
		return result, nil
	case reflect.Interface:
		// If the value is an interface then encode its element
		if v.IsNil() {
			return true, nil
		}
		return ValidateStruct(v.Interface())
	case reflect.Ptr:
		// If the value is a pointer then check its element
		if v.IsNil() {
			return true, nil
		}

		return ValidateStruct(v.Interface())

	case reflect.Struct:
		return ValidateStruct(v.Interface())
	default:
		return false, &UnsupportedTypeError{v.Type()}
	}
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

// Error returns string equivalent for reflect.Type
func (e *UnsupportedTypeError) Error() string {
	return "validator: unsupported type: " + e.Type.String()
}

func (sv stringValues) Len() int           { return len(sv) }
func (sv stringValues) Swap(i, j int)      { sv[i], sv[j] = sv[j], sv[i] }
func (sv stringValues) Less(i, j int) bool { return sv.get(i) < sv.get(j) }
func (sv stringValues) get(i int) string   { return sv[i].String() }

func formatsMessages(structName string, tag string, attribute string, params map[string]string) error {
	validator := New()
	if validator.Translator != nil {
		return validator.Translator(structName, tag, attribute, params)
	}

	if m, ok := validator.CustomMessage[structName+"."+tag]; ok {
		return fmt.Errorf(m)
	}

	message, ok := MessageMap[tag]
	if ok {
		if a, ok := validator.Attributes[structName]; ok {
			attribute = a
		}

		message = strings.Replace(message, ":attribute", attribute, -1)
		for key, value := range params {
			message = strings.Replace(message, ":"+key, value, -1)
		}
		return fmt.Errorf(message)
	}

	return fmt.Errorf("validator: undefined message : %s", tag)
}