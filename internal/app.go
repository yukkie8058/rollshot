package internal

import (
	"errors"
	"image/jpeg"
	"image/png"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/yukkie8058/rollshot/data"
	ilayout "github.com/yukkie8058/rollshot/internal/internal/layout"
)

func ShowMainWindow(a fyne.App, images data.ImageList) {
	w := a.NewWindow(a.Metadata().Name)
	g := NewGlobalizer(nil)

	innerPadding := theme.InnerPadding()

	list := newImageList(w, g, images)
	listVScroll := container.NewVScroll(container.New(
		layout.NewCustomPaddedLayout(innerPadding, innerPadding, innerPadding, innerPadding),
		container.NewBorder(nil, nil, list, nil)),
	)

	reverse := newDynamicButton("Reverse", theme.ViewRefreshIcon(), func() {
		v, _ := images.Get()
		slices.Reverse(v)
		images.Set(v)
		list.Refresh()
	})
	preview := newDynamicButton("Preview", theme.VisibilityIcon(), func() {
		cont := canvas.NewImageFromImage(images.Merge())
		cont.FillMode = canvas.ImageFillContain
		cont.ScaleMode = canvas.ImageScaleFastest
		cont.SetMinSize(imageSizeByBounds(cont.Image.Bounds()))
		dialog.ShowCustom("Preview", "Close", cont, w)
	})
	save := newDynamicButton("Save", theme.DocumentSaveIcon(), func() {
		d := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if writer == nil {
				return
			}
			defer writer.Close()
			switch writer.URI().Extension() {
			case ".jpg", ".jpeg":
				jpeg.Encode(writer, images.Merge(), nil)
			case ".png":
				png.Encode(writer, images.Merge())
			default:
				dialog.ShowError(errors.New("unsupported format"), w)
				return
			}
			popUp := widget.NewPopUp(&widget.Label{
				Text:      "Saved successfully",
				Alignment: fyne.TextAlignCenter,
				TextStyle: fyne.TextStyle{Bold: true},
			}, w.Canvas())
			cs, ps := w.Canvas().Size(), popUp.MinSize()
			popUp.ShowAtPosition(fyne.NewPos(cs.Width/2, cs.Height/2).SubtractXY(ps.Width/2, ps.Height/2))
		}, w)
		d.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/jpeg", "image/png"}))
		d.SetFileName("image.png")
		d.Show()
	})
	save.Importance = widget.HighImportance
	bottomRight := container.New(ilayout.NewAlign(ilayout.AlignRight), reverse, preview, save)

	images.AddListener(binding.NewDataListener(func() {
		if images.Length() > 0 {
			listVScroll.SetMinSize(fyne.NewSquareSize(listVScroll.Content.MinSize().Width))
			bottomRight.Show()
		} else {
			listVScroll.SetMinSize(listVScroll.Content.MinSize())
			bottomRight.Hide()
		}
	}))

	w.SetOnDropped(func(pos fyne.Position, items []fyne.URI) {
		go func() {
			for _, v := range items {
				img, closed := tryLoadImage(w, v)
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
		&mainContentLayout{ilayout.NewCorner(nil, nil, nil, bottomRight), bottomRight, innerPadding},
		listVScroll, bottomRight,
	)
	w.SetContent(g)
	w.Show()
}

type mainContentLayout struct {
	corner *ilayout.Corner

	bottomRight *fyne.Container
	cornerPad   float32
}

func (l *mainContentLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	l.corner.Layout(objects, size)

	for _, obj := range objects {
		if obj == l.bottomRight {
			l.bottomRight.Move(l.bottomRight.Position().Subtract(fyne.NewSquareOffsetPos(l.cornerPad)))
		}
	}
}

func (l *mainContentLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
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

func tryLoadImage(w fyne.Window, uri fyne.URI) (image *data.Image, dialogClosed <-chan struct{}) {
	if img, err := data.LoadImage(uri); err == nil {
		return img, nil
	} else {
		d := dialog.NewError(err, w)
		closed := make(chan struct{})
		d.SetOnClosed(func() { close(closed) })
		d.Show()
		return nil, closed
	}
}
