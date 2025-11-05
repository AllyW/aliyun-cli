// Copyright (c) 2009-present, Alibaba Cloud All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/tabwriter"

	jmespath "github.com/jmespath/go-jmespath"
)

// OutputFilter filters and formats output
type OutputFilter interface {
	FilterOutput(input string) (string, error)
}

// TableOutputFilter formats output as a table
type TableOutputFilter struct {
	Cols    []string
	Rows    string
	ShowNum bool
}

// NewTableOutputFilter creates a new table output filter
func NewTableOutputFilter(cols []string, rows string, showNum bool) OutputFilter {
	return &TableOutputFilter{
		Cols:    cols,
		Rows:    rows,
		ShowNum: showNum,
	}
}

// FilterOutput filters and formats the output as a table
func (f *TableOutputFilter) FilterOutput(s string) (string, error) {
	var v any
	// Wrap in RootFilter array for consistent parsing
	wrapped := fmt.Sprintf("{\"RootFilter\":[%s]}", s)
	decoder := json.NewDecoder(bytes.NewBufferString(wrapped))
	decoder.UseNumber()
	err := decoder.Decode(&v)
	if err != nil {
		return s, fmt.Errorf("unmarshal output failed: %w", err)
	}

	var rowPath string
	if f.Rows != "" {
		rowPath = "RootFilter[0]." + f.Rows
	} else {
		rowPath = "RootFilter"
	}

	if len(f.Cols) == 0 {
		return s, fmt.Errorf("you need to specify columns with --output cols=col1,col2,...")
	}

	return f.FormatTable(rowPath, f.Cols, v)
}

// FormatTable formats data as a table
func (f *TableOutputFilter) FormatTable(rowPath string, colNames []string, v any) (string, error) {
	// Add row number if requested
	if f.ShowNum {
		colNames = append([]string{"Num"}, colNames...)
	}

	rows, err := jmespath.Search(rowPath, v)
	if err != nil {
		return "", fmt.Errorf("jmespath '%s' failed: %w", rowPath, err)
	}

	rowsArray, ok := rows.([]any)
	if !ok {
		return "", fmt.Errorf("jmespath '%s' failed: need array expression", rowPath)
	}

	// Determine data type
	dataType := 1 // 1 = object, 2 = array
	if len(rowsArray) > 0 {
		_, ok := rowsArray[0].(map[string]any)
		if !ok {
			if isArrayOrSlice(rowsArray[0]) {
				dataType = 2
			}
		}
	}

	colNamesArray := make([]string, 0)
	colIndexArray := make([]int, 0)

	if dataType == 2 {
		// Array type: columns must be in "name:index" format
		for _, colName := range colNames {
			if colName == "Num" {
				colNamesArray = append(colNamesArray, colName)
				continue
			}
			if !strings.Contains(colName, ":") {
				return "", fmt.Errorf("colNames: %s must be in 'name:index' format, like 'name:0'", colName)
			}
			parts := strings.Split(colName, ":")
			if len(parts) != 2 {
				return "", fmt.Errorf("colNames: %s must be in 'name:index' format", colName)
			}
			idx, err := strconv.Atoi(parts[1])
			if err != nil {
				return "", fmt.Errorf("colNames: %s must be in 'name:index' format", colName)
			}
			colNamesArray = append(colNamesArray, parts[0])
			colIndexArray = append(colIndexArray, idx)
		}
	} else {
		colNamesArray = colNames
	}

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	w := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)

	// Write header
	header := strings.Join(colNamesArray, "\t")
	fmt.Fprintln(w, header)

	// Write rows
	for i, row := range rowsArray {
		values := make([]string, 0, len(colNamesArray))

		if f.ShowNum {
			values = append(values, strconv.Itoa(i+1))
		}

		if dataType == 2 {
			// Array type
			if arr, ok := row.([]any); ok {
				for _, idx := range colIndexArray {
					if idx < len(arr) {
						values = append(values, formatValue(arr[idx]))
					} else {
						values = append(values, "")
					}
				}
			}
		} else {
			// Object type
			if obj, ok := row.(map[string]any); ok {
				for _, colName := range colNamesArray {
					if colName == "Num" {
						continue
					}
					if val, ok := obj[colName]; ok {
						values = append(values, formatValue(val))
					} else {
						values = append(values, "")
					}
				}
			}
		}

		fmt.Fprintln(w, strings.Join(values, "\t"))
	}

	w.Flush()
	writer.Flush()
	return buf.String(), nil
}

func isArrayOrSlice(value any) bool {
	v := reflect.ValueOf(value)
	return v.Kind() == reflect.Array || v.Kind() == reflect.Slice
}

func formatValue(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case json.Number:
		return val.String()
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case bool:
		return strconv.FormatBool(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
