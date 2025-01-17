// Copyright (c) 2024  The Go-CoreUtils Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package update

import (
	"os"
	"strings"
)

func goWorkProxyEnviron(goProxy string) (environ []string) {
	for _, line := range os.Environ() {
		if strings.HasPrefix(line, "GOWORK=") || strings.HasPrefix(line, "GOPROXY=") {
			continue
		}
		environ = append(environ, line)
	}
	environ = append(environ, "GOWORK=off", "GOPROXY="+goProxy)
	return
}
