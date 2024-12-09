package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/yukkie8058/rollshot/bindingx"
)

type Globalizer struct {
	widget.BaseWidget

	Content fyne.CanvasObject

	MouseOver  binding.Bool
	MousePos   bindingx.Typed[fyne.Position]
	MousePress binding.Bool
}

func NewGlobalizer(content fyne.CanvasObject) *Globalizer {
	g := &Globalizer{
		Content:    content,
		MouseOver:  binding.NewBool(),
		MousePos:   bindingx.NewTyped[fyne.Position](),
		MousePress: binding.NewBool(),
	}
	g.ExtendBaseWidget(g)
	return g
}

func GlobalizerForObject(object fyne.CanvasObject) *Globalizer {
	c := fyne.CurrentApp().Driver().CanvasForObject(object)
	if c == nil {
		return nil
	}
	if g, ok := c.Content().(*Globalizer); ok {
		return g
	}
	return nil
}

func (g *Globalizer) MouseIn(e *desktop.MouseEvent) {
	g.MouseOver.Set(true)
	g.MousePos.Set(e.AbsolutePosition)
}

func (g *Globalizer) MouseMoved(e *desktop.MouseEvent) {
	g.MousePos.Set(e.AbsolutePosition)
}

func (g *Globalizer) MouseOut() {
	g.MouseOver.Set(false)
}

func (g *Globalizer) MouseDown(e *desktop.MouseEvent) {
	g.MousePress.Set(true)
}

func (g *Globalizer) MouseUp(e *desktop.MouseEvent) {
	g.MousePress.Set(false)
}

func (g *Globalizer) CreateRenderer() fyne.WidgetRenderer {
	r := &globalizerRenderer{globalizer: g}
	r.Refresh()
	return r
}

type globalizerRenderer struct {
	fyne.WidgetRenderer
	globalizer *Globalizer
}

func (r *globalizerRenderer) Refresh() {
	r.WidgetRenderer = widget.NewSimpleRenderer(r.globalizer.Content)
}
