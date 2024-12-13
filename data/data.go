package data

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/storage"
	"github.com/yukkie8058/rollshot/bindingx"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

type ImageList struct {
	bindingx.TypedList[*Image]
}

func NewImageList() ImageList {
	return ImageList{bindingx.NewTypedList[*Image]()}
}

var ErrUnsupportedExtension = errors.New("unsupported extension")

func (l ImageList) Save(writer fyne.URIWriteCloser) error {
	switch writer.URI().Extension() {
	case ".jpg", ".jpeg":
		return jpeg.Encode(writer, l.Merge(), nil)
	case ".png":
		return png.Encode(writer, l.Merge())
	default:
		return ErrUnsupportedExtension
	}
}

func (l ImageList) Merge() image.Image {
	images := make([]image.Image, l.Length())
	list, _ := l.Get()
	for i, v := range list {
		images[i] = v.Trim()
	}

	w, h := 0, 0
	for _, v := range images {
		if v.Bounds().Dx() > w {
			w = v.Bounds().Dx()
		}
		h += v.Bounds().Dy()
	}
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	y := 0
	for _, v := range images {
		draw.Draw(dst, dst.Bounds().Add(image.Pt(0, y)), v, v.Bounds().Min, draw.Src)
		y += v.Bounds().Dy()
	}
	return dst
}

type Image struct {
	URI   fyne.URI
	Image image.Image

	TrimLeading, TrimTrailing binding.Int
}

func LoadImage(uri fyne.URI) (*Image, error) {
	i := &Image{URI: uri}
	r, err := storage.Reader(uri)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	i.Image, _, err = image.Decode(r)
	if err != nil {
		return nil, err
	}

	i.TrimLeading = trim{binding.NewInt(), i, trimLeading}
	i.TrimTrailing = trim{binding.NewInt(), i, trimTrailing}
	return i, nil
}

func (i Image) Trim() image.Image {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	tl, _ := i.TrimLeading.Get()
	tt, _ := i.TrimTrailing.Get()
	return i.Image.(subImager).SubImage(image.Rectangle{
		i.Image.Bounds().Min.Add(image.Point{0, tl}),
		i.Image.Bounds().Max.Sub(image.Point{0, tt}),
	})
}

type trimDirection int

const (
	trimLeading trimDirection = iota
	trimTrailing
)

type trim struct {
	binding.Int
	image *Image

	direction trimDirection
}

func (t trim) Set(val int) error {
	h := t.image.Image.Bounds().Dy()

	if val < 0 {
		val = 0
	} else if val > h {
		val = h
	}

	switch t.direction {
	case trimLeading:
		if t.image.TrimTrailing != nil {
			tt, _ := t.image.TrimTrailing.Get()
			if val > h-tt {
				t.image.TrimTrailing.Set(h - val - 1)
			}
		}
	case trimTrailing:
		if t.image.TrimLeading != nil {
			tl, _ := t.image.TrimLeading.Get()
			if val > h-tl {
				t.image.TrimLeading.Set(h - val - 1)
			}
		}
	}

	return t.Int.Set(val)
}
