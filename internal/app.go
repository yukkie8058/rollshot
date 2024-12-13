package internal

import (
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/yukkie8058/rollshot/data"
	ilayout "github.com/yukkie8058/rollshot/internal/internal/layout"
)

type editor struct {
	fyne.Window

	Images data.ImageList

	scroll *container.Scroll
}

func ShowEditor(a fyne.App, images data.ImageList) {
	e := &editor{Window: a.NewWindow("Rollshot"), Images: images}
	g := NewGlobalizer(nil)

	innerPadding := theme.InnerPadding()

	e.scroll = container.NewVScroll(container.New(
		layout.NewCustomPaddedLayout(innerPadding, innerPadding, innerPadding, innerPadding),
		container.NewBorder(nil, nil, newImageList(e, g), nil)),
	)

	reverse := newDynamicButton("Reverse", theme.ViewRefreshIcon(), e.ReverseImages)
	preview := newDynamicButton("Preview", theme.VisibilityIcon(), e.ShowImagePreviewDialog)
	save := newDynamicButton("Save", theme.DocumentSaveIcon(), e.ShowImageSaveDialog)
	save.Importance = widget.HighImportance
	bottomRight := container.New(ilayout.NewAlign(ilayout.AlignRight), reverse, preview, save)

	images.AddListener(binding.NewDataListener(func() {
		if images.Length() > 0 {
			e.scroll.SetMinSize(fyne.NewSquareSize(e.scroll.Content.MinSize().Width))
			bottomRight.Show()
		} else {
			e.scroll.SetMinSize(e.scroll.Content.MinSize())
			bottomRight.Hide()
		}
	}))

	e.SetOnDropped(func(pos fyne.Position, items []fyne.URI) {
		go func() {
			for _, v := range items {
				img, closed := e.tryLoadImage(v)
				if img != nil {
					images.Append(img)
				}
				if closed != nil {
					<-closed
					time.Sleep(time.Second / 60)
				}
			}
		}()
	})

	g.Content = container.New(
		&editorContentLayout{ilayout.NewCorner(nil, nil, nil, bottomRight), bottomRight, innerPadding},
		e.scroll, bottomRight,
	)
	e.SetContent(g)
	e.Show()
}

func (e editor) ReverseImages() {
	v, _ := e.Images.Get()
	slices.Reverse(v)
	e.Images.Set(v)
	e.scroll.Content.Refresh()
}

func (e editor) ShowImageAddDialog() {
	e.ShowImageOpenDialog(func(img *data.Image) { e.Images.Append(img) })
}

func (e editor) ShowImagePreviewDialog() {
	img := canvas.NewImageFromImage(e.Images.Merge())
	img.FillMode = canvas.ImageFillContain
	img.ScaleMode = canvas.ImageScaleFastest
	img.SetMinSize(imageSizeByBounds(img.Image.Bounds()))
	d := dialog.NewCustom("Preview", "Close", img, e)
	d.Show()
}

func (e editor) ShowImageSaveDialog() {
	d := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, e)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()
		if err := e.Images.Save(writer); err != nil {
			dialog.ShowError(err, e)
			return
		}
		popUp := widget.NewPopUp(&widget.Label{
			Text:      "Saved successfully",
			Alignment: fyne.TextAlignCenter,
			TextStyle: fyne.TextStyle{Bold: true},
		}, e.Canvas())
		cs, ps := e.Canvas().Size(), popUp.MinSize()
		popUp.ShowAtPosition(fyne.NewPos(cs.Width/2, cs.Height/2).SubtractXY(ps.Width/2, ps.Height/2))
	}, e)
	d.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/jpeg", "image/png"}))
	d.SetFileName("image.png")
	d.Show()
}

func (e editor) ShowImageOpenDialog(callback func(img *data.Image)) {
	d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, e)
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()
		if img, _ := e.tryLoadImage(reader.URI()); img != nil {
			callback(img)
		}
	}, e)
	d.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/*"}))
	d.Show()
}

func (e editor) tryLoadImage(uri fyne.URI) (image *data.Image, dialogClosed <-chan struct{}) {
	if img, err := data.LoadImage(uri); err == nil {
		return img, nil
	} else {
		d := dialog.NewError(err, e)
		closed := make(chan struct{})
		d.SetOnClosed(func() { close(closed) })
		d.Show()
		return nil, closed
	}
}

type editorContentLayout struct {
	corner *ilayout.Corner

	bottomRight *fyne.Container
	cornerPad   float32
}

func (l *editorContentLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	l.corner.Layout(objects, size)

	for _, obj := range objects {
		if obj == l.bottomRight {
			l.bottomRight.Move(l.bottomRight.Position().Subtract(fyne.NewSquareOffsetPos(l.cornerPad)))
		}
	}
}

func (l *editorContentLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	size := l.corner.MinSize(objects)

	if !l.bottomRight.Visible() {
		return size
	}
	buttonWidth := float32(0)
	for _, button := range l.bottomRight.Objects {
		db, ok := button.(*dynamicButton)
		if ok && db.ShrinkSize().Width > buttonWidth {
			buttonWidth = db.ShrinkSize().Width
		}
	}
	return size.AddWidthHeight(buttonWidth, 0)
}
