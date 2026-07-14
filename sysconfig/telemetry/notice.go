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

package telemetry

import (
	"fmt"
	"io"
)

const innerNoticeText = `[notice] Aliyun CLI detected an inner plugin command. Per internal telemetry policy, usage metrics (command names, duration, result; no parameter values or resource identifiers) are reported to the internal network only.
To disable: aliyun configure telemetry inner-auto disable  or  export ALIYUN_CLI_TELEMETRY_INNER_OPTOUT=1
Details: docs/telemetry-inner-v1.md
This notice is shown once.
`

func maybeShowInnerNotice(configDir string, cfg *Config, stderr io.Writer, quiet bool) {
	if quiet || stderr == nil {
		return
	}
	if cfg != nil && cfg.InnerNoticeShownAt != "" {
		return
	}
	_, _ = fmt.Fprint(stderr, innerNoticeText)
	_ = MarkInnerNoticeShown(configDir, cfg)
}
