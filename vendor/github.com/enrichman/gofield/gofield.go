package gofield

import (
	"reflect"
	"strings"
)

//Reduce will reduce the object following his json specification and selecting only the specified fields.
// The field are selected following the "Facebook convention", they're comma separated values,
// for inner objects use the curly brackets.
// For example: name,address{post_code}
func Reduce(obj interface{}, fields string) interface{} {
	// if no fields are specified just return the full object
	if fields == "" {
		return obj
	}

	valueObj := reflect.ValueOf(obj)

	// if the obj is a Slice iterate and Reduce each item
	if valueObj.Kind() == reflect.Slice {
		slice := make([]interface{}, valueObj.Len())
		for i := 0; i < valueObj.Len(); i++ {
			slice[i] = valueObj.Index(i).Interface()
		}

		sliceUnwinded := make([]interface{}, 0)
		for _, sliceElem := range slice {
			sliceUnwinded = append(sliceUnwinded, Reduce(sliceElem, fields))
		}
		return sliceUnwinded
	}

	// iterate through each field to select it
	result := make(map[string]interface{})

	for _, field := range Split(fields, ",") {

		// check if we need to filter the inner object
		index := strings.Index(field, "{")
		var innerObjectFields string
		if index > -1 && string(field[len(field)-1]) == "}" {
			innerObjectFields = field[index+1 : len(field)-1]
			field = field[:index]
		}

		// extract the element if the the value is a pointer
		if valueObj.Kind() == reflect.Ptr {
			valueObj = valueObj.Elem()
		}

		// extract the element if the the value is a pointer

		if valueObj.Kind() == reflect.Map {
			obj := obj.(map[string]interface{})
			for k := range obj {
				if field == k {

					// check if we have to Reduce more the object
					if innerObjectFields != "" {
						result[k] = Reduce(obj[k], innerObjectFields)
					} else {
						result[k] = obj[k]
					}
				}
			}

		} else {
			for i := 0; i < valueObj.NumField(); i++ {
				fieldObj := valueObj.Field(i)

				typeOfField := valueObj.Type().Field(i)
				tag := typeOfField.Tag.Get("json")
				tag = strings.Split(tag, ",")[0]

				// tag matched! Reduce it
				if field == tag {

					// check if we have to Reduce more the object
					if innerObjectFields != "" {
						result[tag] = Reduce(fieldObj.Interface(), innerObjectFields)
					} else {
						result[tag] = fieldObj.Interface()
					}
				}
			}
		}
	}

	return result
}

//Split the fields string ignoring the elements inside the curly brackets
func Split(fields, sep string) []string {
	fieldsArr := make([]string, 0)

	var count int
	var isInside bool
	var field string

	for _, k := range fields {
		sk := string(k)

		if sk == sep && !isInside {
			fieldsArr = append(fieldsArr, field)
			field = ""
		} else {
			field += sk
			if sk == "{" {
				count++
				isInside = true
			} else if sk == "}" {
				count--
				if count == 0 {
					isInside = false
				}
			}
		}
	}

	if field != "" {
		fieldsArr = append(fieldsArr, field)
	}

	return fieldsArr
}
