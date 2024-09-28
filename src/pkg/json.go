package pkg

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func UnmarshalAndValidate(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("v must be a pointer to a struct")
	}
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return validateStruct(rv.Elem())
}

func validateStruct(v reflect.Value) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("json")
		if tag == "-" {
			continue
		}

		validate := fieldType.Tag.Get("validate")
		if validate == "" {
			continue
		}
		if field.Kind() == reflect.Struct {
			if err := validateStruct(field); err != nil {
				return fmt.Errorf("in field %s: %w", fieldType.Name, err)
			}
			continue
		}
		if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.Struct {
			for j := 0; j < field.Len(); j++ {
				if err := validateStruct(field.Index(j)); err != nil {
					return fmt.Errorf("in slice %s at index %d: %w", fieldType.Name, j, err)
				}
			}
			continue
		}
		if strings.Contains(validate, "required") && field.IsZero() {
			return fmt.Errorf("field %s is required", fieldType.Name)
		}
		if err := validateField(field, validate); err != nil {
			return fmt.Errorf("in field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

func validateField(field reflect.Value, rules string) error {
	for _, rule := range strings.Split(rules, ",") {
		parts := strings.Split(rule, "=")
		if len(parts) != 2 {
			continue
		}

		switch parts[0] {
		case "min":
			min, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("invalid min value: %s", parts[1])
			}
			switch field.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if field.Int() < int64(min) {
					return fmt.Errorf("value must be at least %d", min)
				}
			case reflect.String:
				if len(field.String()) < min {
					return fmt.Errorf("length must be at least %d", min)
				}
			}
		case "max":
			max, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("invalid max value: %s", parts[1])
			}
			switch field.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if field.Int() > int64(max) {
					return fmt.Errorf("value must be at most %d", max)
				}
			case reflect.String:
				if len(field.String()) > max {
					return fmt.Errorf("length must be at most %d", max)
				}
			}
		}
	}
	return nil
}
