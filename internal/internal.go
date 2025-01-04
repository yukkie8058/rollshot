package internal

import "fyne.io/fyne/v2"

func NewInfiniteAnimation(fn func()) *fyne.Animation {
	anim := fyne.NewAnimation(1, nil)
	anim.RepeatCount = -1
	anim.Tick = func(float32) { fn() }
	return anim
}
