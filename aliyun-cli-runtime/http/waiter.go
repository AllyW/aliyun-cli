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
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	jmespath "github.com/jmespath/go-jmespath"
)

type Waiter struct {
	Expr     string
	To       string
	Timeout  time.Duration
	Interval time.Duration
}

func NewWaiter(expr, to string, timeout, interval int) *Waiter {
	return &Waiter{
		Expr:     expr,
		To:       to,
		Timeout:  time.Duration(timeout) * time.Second,
		Interval: time.Duration(interval) * time.Second,
	}
}

func (w *Waiter) CallWith(op *Operation) (string, error) {
	begin := time.Now()

	for {
		response, err := op.MakeRequest()
		if err != nil {
			return "", err
		}

		bodyStr := response.GetBodyString()
		v, err := w.evaluateExpr(bodyStr)
		if err != nil {
			return "", fmt.Errorf("failed to evaluate expression: %w", err)
		}

		if v == w.To {
			return bodyStr, nil
		}

		duration := time.Since(begin)
		if duration > w.Timeout {
			return "", fmt.Errorf("wait '%s' to '%s' timeout (%d seconds), last='%s'",
				w.Expr, w.To, int(w.Timeout.Seconds()), v)
		}

		time.Sleep(w.Interval)
	}
}

func (w *Waiter) evaluateExpr(body string) (string, error) {
	var v any
	err := json.Unmarshal([]byte(body), &v)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	result, err := jmespath.Search(w.Expr, v)
	if err != nil {
		return "", fmt.Errorf("jmespath search failed: %w", err)
	}

	switch val := result.(type) {
	case string:
		return val, nil
	case json.Number:
		return val.String(), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	case int:
		return strconv.Itoa(val), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case bool:
		return strconv.FormatBool(val), nil
	case nil:
		return "", nil
	default:
		return fmt.Sprintf("%v", val), nil
	}
}
