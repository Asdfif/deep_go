package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	tPerson := reflect.TypeOf(person)
	vPerson := reflect.ValueOf(person)
	fieldsCount := tPerson.NumField()

	var result []string
	fmt.Printf("%+v\n", person)

	for i := 0; i < fieldsCount; i++ {
		vField := vPerson.Field(i)
		tField := tPerson.Field(i)

		fieldTag := tField.Tag
		serializedName, present := fieldTag.Lookup("properties")
		if !present {
			continue
		}

		fmt.Printf("%v %v %v '%v'\n", serializedName, present, tField.Name, vField)

		if strings.Contains(serializedName, "omitempty") {
			if vField.IsZero() {
				continue
			}

			serializedName, _ = strings.CutSuffix(serializedName, ",omitempty")
		}

		result = append(result, fmt.Sprintf("%v=%v", serializedName, vField))
	}

	return strings.Join(result, "\n")
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
