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
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func SerializeURLParam(name string, value interface{}, required bool, skipQuote bool) (map[string]string, error) {
	if value == nil {
		if required {
			return nil, fmt.Errorf("URL parameter %s is required", name)
		}
		return map[string]string{}, nil
	}

	var strValue string
	switch v := value.(type) {
	case bool:
		strValue = strconv.FormatBool(v)
	case int:
		strValue = strconv.Itoa(v)
	case int64:
		strValue = strconv.FormatInt(v, 10)
	case float64:
		strValue = strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		strValue = v
	default:
		strValue = fmt.Sprintf("%v", v)
	}

	if !skipQuote {
		strValue = url.QueryEscape(strValue)
	}

	return map[string]string{name: strValue}, nil
}

func SerializeQueryParam(name string, value interface{}, required bool, skipQuote bool, div string) (map[string]string, error) {
	if value == nil {
		if required {
			return nil, fmt.Errorf("query parameter %s is required", name)
		}
		return map[string]string{}, nil
	}

	// Handle list/slice
	if reflect.TypeOf(value).Kind() == reflect.Slice || reflect.TypeOf(value).Kind() == reflect.Array {
		rv := reflect.ValueOf(value)
		if rv.Len() == 0 {
			return map[string]string{}, nil
		}

		// Convert slice to comma-separated string
		var parts []string
		for i := 0; i < rv.Len(); i++ {
			item := rv.Index(i).Interface()
			var strItem string
			switch v := item.(type) {
			case bool:
				strItem = strconv.FormatBool(v)
			case int:
				strItem = strconv.Itoa(v)
			case int64:
				strItem = strconv.FormatInt(v, 10)
			case float64:
				strItem = strconv.FormatFloat(v, 'f', -1, 64)
			case string:
				strItem = v
			default:
				strItem = fmt.Sprintf("%v", v)
			}
			if !skipQuote {
				strItem = url.QueryEscape(strItem)
			}
			parts = append(parts, strItem)
		}

		separator := div
		if separator == "" {
			separator = ","
		}
		return map[string]string{name: strings.Join(parts, separator)}, nil
	}

	var strValue string
	switch v := value.(type) {
	case bool:
		strValue = strconv.FormatBool(v)
	case int:
		strValue = strconv.Itoa(v)
	case int64:
		strValue = strconv.FormatInt(v, 10)
	case float64:
		strValue = strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		strValue = v
	default:
		strValue = fmt.Sprintf("%v", v)
	}

	if !skipQuote {
		strValue = url.QueryEscape(strValue)
	}

	return map[string]string{name: strValue}, nil
}

// ReplacePathPlaceholders replaces placeholders in path with values
// Example: "/api/{resourceId}/items/{itemId}" with {"resourceId": "123", "itemId": "456"}
// Returns: "/api/123/items/456"
func ReplacePathPlaceholders(path string, values map[string]string) string {
	result := path
	for key, value := range values {
		placeholder := "{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}
