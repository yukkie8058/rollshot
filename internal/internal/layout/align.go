package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
)

type AlignDirection int

const (
	AlignLeft AlignDirection = iota
	AlignRight
	AlignTop
	AlignBottom
)

type Align struct {
	box       fyne.Layout
	direction AlignDirection
}

func NewAlign(direction AlignDirection) *Align {
	if direction == AlignLeft || direction == AlignRight {
		return &Align{layout.NewVBoxLayout(), direction}
	} else {
		return &Align{layout.NewHBoxLayout(), direction}
	}
}

func (a *Align) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	a.box.Layout(objects, size)

	for _, obj := range objects {
		if !obj.Visible() {
			continue
		}
		switch a.direction {
		case AlignRight:
			obj.Move(fyne.NewPos(size.Width-obj.MinSize().Width, obj.Position().Y))
		case AlignBottom:
			obj.Move(fyne.NewPos(obj.Position().X, size.Height-obj.MinSize().Height))
		}
		if a.direction == AlignLeft || a.direction == AlignRight {
			obj.Resize(fyne.NewSize(obj.MinSize().Width, obj.Size().Height))
		} else {
			obj.Resize(fyne.NewSize(obj.Size().Width, obj.MinSize().Height))
		}
	}
}

func (a *Align) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return a.box.MinSize(objects)
}
