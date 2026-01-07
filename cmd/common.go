/*
Copyright © (2020-2021) Bizfly Cloud

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"

	"github.com/jedib0t/go-pretty/table"
)

// ProcessDataTables - Process data to generate tables
func ProcessDataTables(data []table.Row, parsingData map[string]interface{}) []table.Row {
	for parentKey, parentValue := range parsingData {
		value := "\n"
		switch parentValue := parentValue.(type) {
		case map[string]interface{}:
			for childKey, childValue := range parentValue {
				value = value + fmt.Sprintf("%v: %v", childKey, childValue)
				if value[len(value)-1:] != "\n" {
					value = value + "\n"
				}
			}
			row := table.Row{parentKey, value}
			data = append(data, row)
		case []interface{}:
			for _, interfaceValue := range parentValue {
				for childKey, childValue := range interfaceValue.(map[string]interface{}) {
					value = value + fmt.Sprintf("%v: %v\n", childKey, childValue)
				}
				value = value + "\n"
			}
			if value[len(value)-1:] == "\n" {
				value = value[:len(value)-1]
			}
			row := table.Row{parentKey, value}
			data = append(data, row)
		default:
			row := table.Row{parentKey, parentValue}
			data = append(data, row)
		}
	}
	return data
}

// SliceContains - Check data in slice
func SliceContains(slice interface{}, val interface{}) (int, bool) {
	switch v := slice.(type) {
	case string:
		if slice == val {
			return 1, true
		}
		return -1, false
	default:
		for i, item := range v.([]string) {
			if item == val {
				return i, true
			}
		}
		return -1, false
	}
}
