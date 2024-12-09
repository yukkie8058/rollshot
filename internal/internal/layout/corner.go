package layout

import (
	"fyne.io/fyne/v2"
)

type Corner struct {
	TopLeft, TopRight, BottomLeft, BottomRight fyne.CanvasObject
}

func NewCorner(topLeft, topRight, bottomLeft, bottomRight fyne.CanvasObject) *Corner {
	return &Corner{topLeft, topRight, bottomLeft, bottomRight}
}

func (c *Corner) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if tl := c.TopLeft; tl != nil && tl.Visible() {
		tl.Move(fyne.NewPos(0, 0))
		tl.Resize(tl.MinSize())
	}
	if tr := c.TopRight; tr != nil && tr.Visible() {
		tr.Move(fyne.NewPos(size.Width-tr.MinSize().Width, 0))
		tr.Resize(tr.MinSize())
	}
	if bl := c.BottomLeft; bl != nil && bl.Visible() {
		bl.Move(fyne.NewPos(0, size.Height-bl.MinSize().Height))
		bl.Resize(bl.MinSize())
	}
	if br := c.BottomRight; br != nil && br.Visible() {
		br.Move(fyne.NewPos(size.Width-br.MinSize().Width, size.Height-br.MinSize().Height))
		br.Resize(br.MinSize())
	}

	for _, obj := range objects {
		if !obj.Visible() {
			continue
		}
		if obj == c.TopLeft || obj == c.TopRight || obj == c.BottomLeft || obj == c.BottomRight {
			continue
		}
		obj.Move(fyne.NewPos(0, 0))
		obj.Resize(size)
	}
}

func (c *Corner) MinSize(objects []fyne.CanvasObject) fyne.Size {
	size := fyne.NewSize(0, 0)
	for _, obj := range objects {
		if !obj.Visible() {
			continue
		}
		if obj == c.TopLeft || obj == c.TopRight || obj == c.BottomLeft || obj == c.BottomRight {
			continue
		}
		size = size.Max(obj.MinSize())
	}
	return size
}
