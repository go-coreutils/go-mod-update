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
	"path/filepath"

	"github.com/go-curses/ctk"

	update "github.com/go-coreutils/go-mod-update"
)

type cProject struct {
	u *cUI

	Path string
	Name string

	Frame ctk.Frame
	VBox  ctk.VBox

	Packages []*cPackage
}

func (u *cUI) newProject(path string) (p *cProject) {
	p = &cProject{
		Path: path,
		Name: filepath.Base(path),
		u:    u,
	}
	p.Frame = ctk.NewFrame(p.Path)
	p.Frame.Show()
	p.Frame.SetName("project-entry")
	p.Frame.SetLabelAlign(0.0, 0.5)
	//_ = p.Frame.SetBoolProperty(ctk.PropertyDebug, true)
	p.VBox = ctk.NewVBox(false, 0)
	p.VBox.Show()
	p.VBox.SetName("project-packages")
	//_ = p.VBox.SetBoolProperty(ctk.PropertyDebugChildren, true)
	p.Frame.Add(p.VBox)
	return
}

func (p *cProject) setTitle() {
	if count := p.Pending(); count > 0 {
		p.Frame.SetLabel(fmt.Sprintf("%s (%d pending)", p.Name, count))
	} else {
		p.Frame.SetLabel(fmt.Sprintf("%s (none pending)", p.Name))
	}
}

func (p *cProject) UpdateTitle() {
	p.setTitle()
	p.Frame.Resize()
	p.u.Display.RequestDraw()
	p.u.Display.RequestShow()
}

func (p *cProject) Add(modules ...*update.Module) {
	for _, module := range modules {
		pkg := p.u.newPackage(p, module)
		p.Packages = append(p.Packages, pkg)
		p.VBox.PackStart(pkg.HBox, false, false, 0)
		if pkg.Error != nil {
			p.VBox.PackStart(pkg.Error, true, true, 0)
		}
	}
	p.setTitle()
	p.Resize()
	return
}

func (p *cProject) Height() (h int) {
	if h = 2; len(p.Packages) > 0 {
		for _, pkg := range p.Packages {
			h += 1 // the module update itself
			if pkg.Error != nil {
				h += 1 // the module updated with error
			}
		}
	}
	return
}

func (p *cProject) Resize() {
	w := p.u.getContentWidth()
	h := p.Height()

	p.Frame.SetSizeRequest(w, h)
	p.Frame.Resize()
}

func (p *cProject) Pending() (count int) {
	for _, pkg := range p.Packages {
		if !pkg.Module.Done {
			count += 1
		}
	}
	return
}

func (p *cProject) Refresh() {
	p.Frame.Freeze()
	defer p.Frame.Thaw()

	var modules update.Modules
	for _, pkg := range p.Packages {
		modules = append(modules, pkg.Module)
	}
	p.Packages = make([]*cPackage, 0)

	for _, child := range p.VBox.GetChildren() {
		p.VBox.Remove(child)
		child.Destroy()
	}

	p.Add(modules...)
}
