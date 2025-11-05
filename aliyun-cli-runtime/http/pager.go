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
	"math"
	"strconv"
	"strings"

	jmespath "github.com/jmespath/go-jmespath"
)

// Pager handles pagination for pageable APIs
type Pager struct {
	// Configuration
	PageNumberExpr string
	PageSizeExpr   string
	TotalCountExpr string
	NextTokenExpr  string
	CollectionPath string

	// State
	PageSize          int
	totalCount        int
	currentPageNumber int
	nextTokenMode     bool
	nextToken         string
	results           []any
}

// NewPager creates a new pager with configuration
func NewPager(config map[string]string) *Pager {
	pager := &Pager{
		results: make([]any, 0),
	}

	if pageNumber, ok := config["PageNumber"]; ok {
		pager.PageNumberExpr = pageNumber
	}
	if pageSize, ok := config["PageSize"]; ok {
		pager.PageSizeExpr = pageSize
	}
	if totalCount, ok := config["TotalCount"]; ok {
		pager.TotalCountExpr = totalCount
	}
	if nextToken, ok := config["NextToken"]; ok {
		pager.NextTokenExpr = nextToken
	}
	if path, ok := config["path"]; ok {
		pager.CollectionPath = path
	}

	return pager
}

func (p *Pager) CallWith(op *Operation) (string, error) {
	for {
		response, err := op.MakeRequest()
		if err != nil {
			return "", err
		}

		err = p.FeedResponse(response.GetBodyString())
		if err != nil {
			return "", fmt.Errorf("call failed: %w", err)
		}

		if !p.HasMore() {
			break
		}

		if err := p.MoveNextPage(op); err != nil {
			return "", fmt.Errorf("failed to move to next page: %w", err)
		}
	}

	return p.GetResponseCollection(), nil
}

func (p *Pager) FeedResponse(body string) error {
	var j any
	err := json.Unmarshal([]byte(body), &j)
	if err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}

	if p.CollectionPath != "" {
		items, err := jmespath.Search(p.CollectionPath, j)
		if err == nil {
			if arr, ok := items.([]any); ok {
				p.results = append(p.results, arr...)
			}
		}
	} else {
		p.results = append(p.results, j)
	}

	if p.NextTokenExpr != "" {
		token, err := jmespath.Search(p.NextTokenExpr, j)
		if err == nil {
			if tokenStr, ok := token.(string); ok && tokenStr != "" {
				p.nextToken = tokenStr
				p.nextTokenMode = true
			} else {
				p.nextToken = ""
			}
		}
	}

	if p.TotalCountExpr != "" {
		count, err := jmespath.Search(p.TotalCountExpr, j)
		if err == nil {
			if countNum, ok := count.(json.Number); ok {
				if n, err := countNum.Int64(); err == nil {
					p.totalCount = int(n)
				}
			} else if countNum, ok := count.(float64); ok {
				p.totalCount = int(countNum)
			} else if countStr, ok := count.(string); ok {
				if n, err := strconv.Atoi(countStr); err == nil {
					p.totalCount = n
				}
			}
		}
	}

	if p.PageSizeExpr != "" {
		size, err := jmespath.Search(p.PageSizeExpr, j)
		if err == nil {
			if sizeNum, ok := size.(json.Number); ok {
				if n, err := sizeNum.Int64(); err == nil {
					p.PageSize = int(n)
				}
			} else if sizeNum, ok := size.(float64); ok {
				p.PageSize = int(sizeNum)
			} else if sizeStr, ok := size.(string); ok {
				if n, err := strconv.Atoi(sizeStr); err == nil {
					p.PageSize = n
				}
			}
		}
	}

	if p.PageNumberExpr != "" {
		pageNum, err := jmespath.Search(p.PageNumberExpr, j)
		if err == nil {
			if num, ok := pageNum.(json.Number); ok {
				if n, err := num.Int64(); err == nil {
					p.currentPageNumber = int(n)
				}
			} else if num, ok := pageNum.(float64); ok {
				p.currentPageNumber = int(num)
			} else if numStr, ok := pageNum.(string); ok {
				if n, err := strconv.Atoi(numStr); err == nil {
					p.currentPageNumber = n
				}
			}
		}
	}

	return nil
}

func (p *Pager) HasMore() bool {
	if p.nextTokenMode {
		return p.nextToken != ""
	}
	if p.PageSize > 0 && p.totalCount > 0 {
		pages := int(math.Ceil(float64(p.totalCount) / float64(p.PageSize)))
		return p.currentPageNumber < pages
	}
	return false
}

func (p *Pager) MoveNextPage(op *Operation) error {
	if p.nextTokenMode {
		if p.NextTokenExpr != "" {
			parts := strings.Split(p.NextTokenExpr, ".")
			if len(parts) > 0 {
				tokenKey := parts[len(parts)-1]
				op.AddQueryParam(tokenKey, p.nextToken)
			}
		}
	} else {
		p.currentPageNumber++
		if p.PageNumberExpr != "" {
			parts := strings.Split(p.PageNumberExpr, ".")
			if len(parts) > 0 {
				pageKey := parts[len(parts)-1]
				op.AddQueryParam(pageKey, strconv.Itoa(p.currentPageNumber))
			}
		}
	}
	return nil
}

func (p *Pager) GetResponseCollection() string {
	if p.CollectionPath == "" {
		result, err := json.Marshal(p.results)
		if err != nil {
			return "[]"
		}
		return string(result)
	}

	root := make(map[string]any)
	current := make(map[string]any)
	path := p.CollectionPath

	parts := strings.Split(path, ".")
	if len(parts) > 1 {
		// Nested path
		root[parts[0]] = current
		key := strings.TrimSuffix(parts[len(parts)-1], "[]")
		current[key] = p.results
	} else {
		// Simple path
		key := strings.TrimSuffix(path, "[]")
		root[key] = p.results
	}

	result, err := json.Marshal(root)
	if err != nil {
		return "{}"
	}
	return string(result)
}
