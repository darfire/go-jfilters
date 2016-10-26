package json_filters

import (
	"fmt"
	"strconv"
)

type StepFunc func(v interface{}) interface{}

type Step struct {
	Description string
	Apply       StepFunc
}

type Filter struct {
	Steps []Step
}

type BoundFilter struct {
	Filter *Filter
	Value  interface{}
}

func (f *Filter) Bind(v interface{}) *BoundFilter {
	return &BoundFilter{Filter: f, Value: v}
}

func New() *Filter {
	return &Filter{Steps: make([]Step, 0)}
}

func (f *Filter) appendStep(s Step) []Step {
	return append(f.Steps, s)
}

func (f *Filter) Map() *Filter {
	return &Filter{
		Steps: f.appendStep(Step{"Map", checkIsMap}),
	}
}

func (f *Filter) Array() *Filter {
	return &Filter{
		Steps: f.appendStep(Step{"Array", checkIsArray}),
	}
}

func (f *Filter) String() *Filter {
	return &Filter{
		Steps: f.appendStep(Step{"String", checkIsString}),
	}
}

func (f *Filter) Number() *Filter {
	return &Filter{
		Steps: f.appendStep(Step{"Number", checkIsNumber}),
	}
}

func (f *Filter) Key(key string) *Filter {
	return &Filter{
		Steps: f.appendStep(Step{"Key", func(v interface{}) interface{} {
			return selectMapKey(v, key)
		}}),
	}
}

func (f *Filter) Index(index int) *Filter {
	return &Filter{
		Steps: f.appendStep(Step{"Index", func(v interface{}) interface{} {
			return selectArrayIndex(v, index)
		}}),
	}
}

func (f *Filter) Apply(v interface{}) (result interface{}, err error) {
	for i := range f.Steps {
		v, err = f.ApplyStep(i, v)
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

func (f *Filter) ApplyStep(index int, v interface{}) (ret interface{}, err error) {
	s := &f.Steps[index]
	defer func() {
		r := recover()
		if r != nil {
			// reformat a little, show index
			err = fmt.Errorf(fmt.Sprintf("%v(%v): %v", s.Description, index, r))
		}
	}()

	ret = s.Apply(v)
	return ret, nil
}

// Step functions

func checkIsMap(v interface{}) interface{} {
	if m, ok := v.(map[string]interface{}); !ok {
		panic(fmt.Errorf("Expected map"))
	} else {
		return m
	}
}

func checkIsArray(v interface{}) interface{} {
	if a, ok := v.([]interface{}); !ok {
		panic(fmt.Errorf("Expected array"))
	} else {
		return a
	}
}

func checkIsString(v interface{}) interface{} {
	if s, err := convertToString(v); err != nil {
		panic(err)
	} else {
		return s
	}
}

func checkIsBool(v interface{}) interface{} {
	if b, err := convertToBool(v); err != nil {
		panic(err)
	} else {
		return b
	}
}

func checkIsNumber(v interface{}) interface{} {
	s, err := validateNumber(v)
	if err != nil {
		panic(err)
	}
	return s
}

func selectMapKey(v interface{}, k string) interface{} {
	m := checkIsMap(v).(map[string]interface{})
	if x, ok := m[k]; ok {
		return x
	} else {
		panic(fmt.Errorf(fmt.Sprintf("Key %v not found", k)))
	}
}

func selectArrayIndex(v interface{}, index int) interface{} {
	a := checkIsArray(v).([]interface{})
	if index < 0 || index >= len(a) {
		panic(fmt.Errorf(fmt.Sprintf("Index %v out of range", index)))
	}
	return a[index]
}

func convertToString(v interface{}) (string, error) {
	if s, ok := v.(string); ok {
		return s, nil
	} else {
		return "", fmt.Errorf("Not a string")
	}
}

func convertToBool(v interface{}) (bool, error) {
	if b, ok := v.(bool); ok {
		return b, nil
	} else {
		return false, fmt.Errorf("Not a bool")
	}
}

func validateNumber(v interface{}) (string, error) {
	s := fmt.Sprintf("%v", v)
	_, err := strconv.Atoi(s)
	if err == nil {
		return s, nil
	}
	_, err = strconv.ParseFloat(s, 64)
	if err == nil {
		return s, nil
	}
	return "", err
}

// BoundFilter functions

func (f *BoundFilter) IsValid() bool {
	_, err := f.Get()
	return err == nil
}

func (f *BoundFilter) Get() (interface{}, error) {
	return f.Filter.Apply(f.Value)
}

func (f *BoundFilter) GetString() (string, error) {
	v, err := f.Get()
	if err == nil {
		return convertToString(v)
	} else {
		return "", err
	}
}

func (f *BoundFilter) GetBool() (bool, error) {
	v, err := f.Get()
	if err == nil {
		return convertToBool(v)
	} else {
		return false, err
	}
}

func (f *BoundFilter) GetInt() (int, error) {
	v, err := f.Get()
	if err == nil {
		s, err := validateNumber(v)
		if err == nil {
			return strconv.Atoi(s)
		}
	}
	return 0, err
}

func (f *BoundFilter) GetFloat() (float64, error) {
	v, err := f.Get()
	if err == nil {
		s := fmt.Sprintf("%v", v)
		if err == nil {
			return strconv.ParseFloat(s, 64)
		}
	}
	return 0.0, err
}

func (f *BoundFilter) GetMap() (map[string]interface{}, error) {
	v, err := f.Get()
	if err == nil {
		if m, ok := v.(map[string]interface{}); ok {
			return m, nil
		} else {
			return nil, fmt.Errorf("Expected map")
		}
	}
	return nil, err
}

func (f *BoundFilter) GetArray() ([]interface{}, error) {
	v, err := f.Get()
	if err == nil {
		if a, ok := v.([]interface{}); ok {
			return a, nil
		} else {
			return nil, fmt.Errorf("Expected array")
		}
	}
	return nil, err
}
