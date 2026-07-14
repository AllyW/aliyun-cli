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

package cli

var postExecuteHook func(ctx *Context, args []string, err error)

// SetPostExecuteHook registers a callback invoked once after the root command finishes.
func SetPostExecuteHook(fn func(ctx *Context, args []string, err error)) {
	postExecuteHook = fn
}

func invokePostExecuteHook(ctx *Context, args []string, err error) {
	if postExecuteHook != nil {
		postExecuteHook(ctx, args, err)
	}
}
