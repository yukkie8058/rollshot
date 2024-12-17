package internal

import (
	"image"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fyne-io/glfw-js"
	"github.com/srwiley/rasterx"
	"github.com/yukkie8058/rollshot/data"
	"golang.org/x/image/math/fixed"
)

type imageList struct {
	widget.BaseWidget

	Editor *editor

	container *fyne.Container
}

func newImageList(e *editor, g *Globalizer) *imageList {
	l := &imageList{
		Editor:    e,
		container: container.New(imageListLayout{layout.NewVBoxLayout()}),
	}
	l.ExtendBaseWidget(l)

	listener := binding.NewDataListener(l.refreshItemSliders)
	g.MouseOver.AddListener(listener)
	g.MousePos.AddListener(listener)

	e.Images.AddListener(binding.NewDataListener(l.Refresh))
	return l
}

func (l *imageList) Refresh() {
	l.container.RemoveAll()
	val, _ := l.Editor.Images.Get()
	for i, v := range val {
		l.container.Add(newImageItem(l, i, v))
	}
	l.container.Add(newImageAddButton(l.Editor.ShowImageAddDialog))

	l.BaseWidget.Refresh()
}

func (l *imageList) refreshItemSliders() {
	for _, obj := range l.container.Objects {
		if item, ok := obj.(*imageItem); ok {
			item.RefreshSliders()
		}
	}
}

func (l *imageList) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.New(
		layout.NewCustomPaddedLayout((&imageSliderThumb{}).MinSize().Height/2, 0, 0, 0),
		l.container,
	))
}

type imageListLayout struct {
	vBox fyne.Layout
}

func (l imageListLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	l.vBox.Layout(objects, size)

	for _, obj := range objects {
		obj.Resize(fyne.NewSize(obj.MinSize().Width, obj.Size().Height))
		obj.Move(fyne.NewPos(size.Width/2-obj.Size().Width/2, obj.Position().Y))
	}
	if len(objects) == 1 {
		objects[0].Move(fyne.NewPos(
			objects[0].Position().X,
			size.Height/2-objects[0].MinSize().Height/2,
		))
	}
}

func (l imageListLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	size := l.vBox.MinSize(objects)
	return fyne.NewSize(
		float32(math.Max(float64(size.Width), float64(imageBaseSize().Width))),
		float32(math.Max(float64(size.Height), float64(imageBaseSize().Height))),
	)
}

func imageBaseSize() fyne.Size {
	w := float32(400)
	m := glfw.GetPrimaryMonitor().GetVideoMode()
	return fyne.NewSize(w, w*float32(m.Height)/float32(m.Width))
}

func imageSizeByBounds(b image.Rectangle) fyne.Size {
	w := imageBaseSize().Width
	return fyne.NewSize(w, w*float32(b.Dy())/float32(b.Dx()))
}

type imageItem struct {
	widget.BaseWidget

	List  *imageList
	Index int
	Data  *data.Image

	sliderContainer *fyne.Container
}

func newImageItem(list *imageList, index int, data *data.Image) *imageItem {
	i := &imageItem{List: list, Index: index, Data: data}
	i.ExtendBaseWidget(i)
	return i
}

func (i *imageItem) TappedSecondary(e *fyne.PointEvent) {
	canMoveUp := i.Index > 0
	canMoveDown := i.Index < i.List.Editor.Images.Length()-1
	widget.ShowPopUpMenuAtPosition(fyne.NewMenu(i.Data.URI.Name(),
		&fyne.MenuItem{Icon: theme.MoveUpIcon(), Label: "Move Up", Action: func() {
			if canMoveUp {
				left, _ := i.List.Editor.Images.GetValue(i.Index - 1)
				i.List.Editor.Images.SetValue(i.Index, left)
				i.List.Editor.Images.SetValue(i.Index-1, i.Data)
				i.List.Refresh()
			}
		}, Disabled: !canMoveUp},
		&fyne.MenuItem{Icon: theme.MoveDownIcon(), Label: "Move Down", Action: func() {
			if canMoveDown {
				right, _ := i.List.Editor.Images.GetValue(i.Index + 1)
				i.List.Editor.Images.SetValue(i.Index, right)
				i.List.Editor.Images.SetValue(i.Index+1, i.Data)
				i.List.Refresh()
			}
		}, Disabled: !canMoveDown},
		fyne.NewMenuItemSeparator(),
		&fyne.MenuItem{Icon: theme.DeleteIcon(), Label: "Remove", Action: func() {
			i.List.Editor.Images.Remove(i.Data)
		}},
	), fyne.CurrentApp().Driver().CanvasForObject(i), e.AbsolutePosition)
}

func (i *imageItem) MinSize() fyne.Size {
	return imageSizeByBounds(i.Data.Image.Bounds()).
		AddWidthHeight((&imageSliderThumb{}).MinSize().Width*2, 0)
}

func (i *imageItem) RefreshSliders() {
	i.sliderContainer.Refresh()
}

func (i *imageItem) CreateRenderer() fyne.WidgetRenderer {
	image := canvas.NewImageFromImage(i.Data.Image)
	image.FillMode = canvas.ImageFillContain
	image.ScaleMode = canvas.ImageScaleFastest
	image.SetMinSize(imageSizeByBounds(i.Data.Image.Bounds()))

	i.sliderContainer = container.NewWithoutLayout(
		newImageSlider(i, sliderDirectionDown),
		newImageSlider(i, sliderDirectionUp),
	)

	return widget.NewSimpleRenderer(container.NewStack(
		container.NewCenter(image),
		i.sliderContainer,
	))
}

type sliderDirection int

const (
	sliderDirectionDown sliderDirection = iota
	sliderDirectionUp
)

type imageSlider struct {
	widget.BaseWidget
	Image *imageItem

	Direction sliderDirection
}

func newImageSlider(image *imageItem, direction sliderDirection) *imageSlider {
	s := &imageSlider{
		Image:     image,
		Direction: direction,
	}
	s.ExtendBaseWidget(s)
	switch direction {
	case sliderDirectionDown:
		image.Data.TrimLeading.AddListener(binding.NewDataListener(s.Refresh))
	case sliderDirectionUp:
		image.Data.TrimTrailing.AddListener(binding.NewDataListener(s.Refresh))
	}
	return s
}

func (s *imageSlider) CreateRenderer() fyne.WidgetRenderer {
	th := s.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	return &imageSliderRenderer{
		slider:     s,
		overlay:    canvas.NewRectangle(th.Color(theme.ColorNameShadow, v)),
		leftThumb:  newImageSliderThumb(s),
		rightThumb: newImageSliderThumb(s),
	}
}

type imageSliderRenderer struct {
	slider *imageSlider

	overlay               *canvas.Rectangle
	leftThumb, rightThumb *imageSliderThumb
}

func (r *imageSliderRenderer) Refresh() {
	img := r.slider.Image

	scale := img.Size().Height / float32(img.Data.Image.Bounds().Dy())
	offset := (&imageSliderThumb{}).MinSize().Height / 2
	switch r.slider.Direction {
	case sliderDirectionDown:
		v, _ := img.Data.TrimLeading.Get()
		r.slider.Move(fyne.NewPos(0, 0))
		r.slider.Resize(fyne.NewSize(img.Size().Width, float32(v)*scale+offset))
	case sliderDirectionUp:
		v, _ := img.Data.TrimTrailing.Get()
		r.slider.Move(fyne.NewPos(0, float32(img.Data.Image.Bounds().Dy()-v)*scale-offset))
		r.slider.Resize(fyne.NewSize(img.Size().Width, float32(v)*scale+offset))
	}

	const margin = 40
	g := GlobalizerForObject(r.slider)
	if g == nil {
		return
	}
	mouseOver, _ := g.MouseOver.Get()
	mousePos, _ := g.MousePos.Get()
	mousePress, _ := g.MousePress.Get()
	pos1 := fyne.CurrentApp().Driver().AbsolutePositionForObject(img).SubtractXY(margin, margin)
	pos2 := pos1.Add(img.Size()).AddXY(margin*2, margin*2)
	in := pos1.X <= mousePos.X && mousePos.X <= pos2.X &&
		pos1.Y <= mousePos.Y && mousePos.Y <= pos2.Y
	if mouseOver && in {
		r.leftThumb.Show()
		r.rightThumb.Show()
	} else if !mousePress {
		r.leftThumb.Hide()
		r.rightThumb.Hide()
	}

	r.Layout(r.slider.Size())
}

func (r *imageSliderRenderer) Layout(size fyne.Size) {
	ts := r.leftThumb.MinSize()

	oy := float32(0)
	switch r.slider.Direction {
	case sliderDirectionDown:
		oy = 0
	case sliderDirectionUp:
		oy = ts.Height / 2
	}
	r.overlay.Move(fyne.NewPos(ts.Width, oy))
	r.overlay.Resize(size.SubtractWidthHeight(ts.Width*2, ts.Height/2))

	ty := float32(0)
	switch r.slider.Direction {
	case sliderDirectionDown:
		ty = r.slider.Size().Height - ts.Height
	case sliderDirectionUp:
		ty = 0
	}
	r.leftThumb.Move(fyne.NewPos(0, ty))
	r.leftThumb.Resize(r.leftThumb.MinSize())
	r.rightThumb.Move(fyne.NewPos(r.slider.Size().Width-ts.Width, ty))
	r.rightThumb.Resize(r.rightThumb.MinSize())
}

func (r *imageSliderRenderer) MinSize() fyne.Size {
	return fyne.NewSquareSize(1)
}

func (r *imageSliderRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.overlay, r.leftThumb, r.rightThumb}
}

func (r *imageSliderRenderer) Destroy() {}

type imageSliderThumb struct {
	widget.BaseWidget
	slider *imageSlider
}

func newImageSliderThumb(slider *imageSlider) *imageSliderThumb {
	t := &imageSliderThumb{slider: slider}
	t.ExtendBaseWidget(t)
	return t
}

func (t *imageSliderThumb) Cursor() desktop.Cursor {
	return desktop.VResizeCursor
}

func (t *imageSliderThumb) Dragged(e *fyne.DragEvent) {
	s := t.slider

	scaled := int(e.Dragged.DY * float32(s.Image.Data.Image.Bounds().Dy()) / s.Image.Size().Height)
	switch s.Direction {
	case sliderDirectionDown:
		val, _ := s.Image.Data.TrimLeading.Get()
		s.Image.Data.TrimLeading.Set(val + scaled)
	case sliderDirectionUp:
		val, _ := s.Image.Data.TrimTrailing.Get()
		s.Image.Data.TrimTrailing.Set(val - scaled)
	}

	s.Refresh()
}

func (t *imageSliderThumb) DragEnd() {}

func (t *imageSliderThumb) MinSize() fyne.Size {
	th := t.Theme()
	return fyne.NewSquareSize(th.Size(theme.SizeNameInlineIcon) + th.Size(theme.SizeNamePadding))
}

func (t *imageSliderThumb) CreateRenderer() fyne.WidgetRenderer {
	th := t.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	var icon fyne.Resource
	switch t.slider.Direction {
	case sliderDirectionDown:
		icon = th.Icon(theme.IconNameMoveDown)
	case sliderDirectionUp:
		icon = th.Icon(theme.IconNameMoveUp)
	}

	return widget.NewSimpleRenderer(container.NewStack(
		&canvas.Circle{FillColor: th.Color(theme.ColorNameForeground, v)},
		widget.NewIcon(theme.NewInvertedThemedResource(icon)),
	))
}

type imageAddButton struct {
	widget.BaseWidget

	OnTapped func()
}

func newImageAddButton(tapped func()) *imageAddButton {
	b := &imageAddButton{
		OnTapped: tapped,
	}
	b.ExtendBaseWidget(b)
	return b
}

func (b *imageAddButton) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (b *imageAddButton) Tapped(*fyne.PointEvent) {
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *imageAddButton) MinSize() fyne.Size {
	const scale = 0.8
	size := imageBaseSize()
	return fyne.NewSize(size.Width*scale, size.Height*scale)
}

func (b *imageAddButton) CreateRenderer() fyne.WidgetRenderer {
	const clr = theme.ColorNameDisabled
	th := b.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	const outlineWidth = 16
	const outlineDashes = 32
	outline := canvas.NewRaster(func(width, height int) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		scanner := rasterx.NewScannerGV(width, height, img, img.Bounds())
		dasher := rasterx.NewDasher(width, height, scanner)
		dasher.SetColor(th.Color(clr, v))
		dasher.SetStroke(fixed.Int26_6(outlineWidth*64), 0, nil, nil, nil, rasterx.Round, []float64{outlineDashes}, 0)

		hw := outlineWidth / 2.0
		dasher.Start(rasterx.ToFixedP(hw, hw))
		dasher.Line(rasterx.ToFixedP(float64(width)-hw, hw))
		dasher.Line(rasterx.ToFixedP(float64(width)-hw, float64(height)-hw))
		dasher.Line(rasterx.ToFixedP(hw, float64(height)-hw))
		dasher.Stop(true)

		dasher.Draw()
		return img
	})

	icon := widget.NewIcon(theme.NewColoredResource(th.Icon(theme.IconNameContentAdd), clr))

	text := canvas.NewText("Drop an image here or click to add", th.Color(clr, v))
	text.Alignment = fyne.TextAlignCenter
	text.TextStyle = fyne.TextStyle{Bold: true}

	return widget.NewSimpleRenderer(container.NewStack(
		outline,
		container.New(
			layout.NewCustomPaddedLayout(outlineWidth, outlineWidth, outlineWidth, outlineWidth),
			container.NewPadded(
				icon,
				container.NewBorder(nil, text, nil, nil),
			),
		),
	))
}
