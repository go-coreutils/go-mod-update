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

package ui

import (
	"fmt"
	"time"

	"github.com/go-corelibs/spinner"

	update "github.com/go-coreutils/go-mod-update"
)

func (u *cUI) requestUpdatesStatusUpdater(idx, moduleCount int, project, previous *cProject, symbol string, overrideMessage string) {
	if project != nil {
		project.Frame.SetLabel(symbol + " " + project.Name)
		project.Frame.Resize()
	}
	if previous != nil && previous.Path != project.Path {
		previous.UpdateTitle()
		previous = project
	}
	if overrideMessage == "" {
		if moduleCount > 1 {
			u.StatusLabel.SetLabel(fmt.Sprintf("updating... (%d of %d)", idx+1, moduleCount))
		} else {
			u.StatusLabel.SetLabel("updating...")
		}
	} else {
		u.StatusLabel.SetLabel(overrideMessage)
	}
	u.Display.RequestDraw()
	u.Display.RequestShow()
	return
}

func (u *cUI) requestUpdates() {
	u.modLock.Lock()
	defer u.modLock.Unlock()

	u.SetState(UpdatingState)

	var s spinner.Spinner
	var idx, moduleCount int
	var overrideMessage string
	var project, previous *cProject

	for _, p := range u.Projects {
		for _, pkg := range p.Packages {
			if pkg.Module.Pick && !pkg.Module.Done {
				moduleCount += 1
			}
		}
	}

	s = spinner.New(spinner.DefaultSymbols, func(symbol string) {
		u.requestUpdatesStatusUpdater(idx, moduleCount, project, previous, symbol, overrideMessage)
	})

	for _, project = range u.Projects {
		var updated bool
		for _, pkg := range project.Packages {
			if updated = pkg.Module.Pick && !pkg.Module.Done; updated {
				u.requestUpdatesStatusUpdater(idx, moduleCount, project, previous, s.String(), overrideMessage)
				pkg.GoModUpdate()
				pkg.Module.Pick = false
				u.requestUpdatesStatusUpdater(idx, moduleCount, project, previous, s.String(), overrideMessage)
				idx += 1
			}
		}
		if u.tidy && updated {
			overrideMessage = "go mod tidy: " + project.Name
			u.requestUpdatesStatusUpdater(idx, moduleCount, project, previous, s.String(), overrideMessage)
			time.Sleep(time.Millisecond * 500)
			if err := update.Tidy(project.Path, u.goProxy); err != nil {
				err = fmt.Errorf("%q error: %v", project.Name, err)
				u.ErrorList.PackStart(u.makeError(err), true, true, 0)
			}
			time.Sleep(time.Millisecond * 500)
			overrideMessage = ""
			u.requestUpdatesStatusUpdater(idx, moduleCount, project, previous, s.String(), overrideMessage)
		}
		previous = project
	}

	s.Stop()

	u.SetState(IdleState)
}
