package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type dynamicButton struct {
	widget.Button

	Text string

	shrinkSize fyne.Size
}

func newDynamicButton(label string, icon fyne.Resource, tapped func()) *dynamicButton {
	b := &dynamicButton{Text: label}
	b.ExtendBaseWidget(b)
	b.Icon = icon
	b.OnTapped = tapped
	return b
}

func (b *dynamicButton) MouseIn(e *desktop.MouseEvent) {
	b.Button.SetText(b.Text)
	b.Button.MouseIn(e)
}

func (b *dynamicButton) MouseMoved(e *desktop.MouseEvent) {
	b.Button.MouseMoved(e)
}

func (b *dynamicButton) MouseOut() {
	b.Button.SetText("")
	b.Button.MouseOut()
}

func (b *dynamicButton) MinSize() fyne.Size {
	size := b.Button.MinSize()
	if b.shrinkSize.IsZero() {
		b.shrinkSize = size
	} else {
		b.shrinkSize = b.shrinkSize.Min(size)
	}
	return size
}

func (b *dynamicButton) ShrinkSize() fyne.Size { return b.shrinkSize }
