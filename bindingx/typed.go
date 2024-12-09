package bindingx

import "fyne.io/fyne/v2/data/binding"

type Typed[T comparable] interface {
	binding.DataItem

	Get() (T, error)
	Set(T) error
}

type ExternalTyped[T comparable] interface {
	Typed[T]

	Reload() error
}

func NewTyped[T comparable]() Typed[T]      { return &boundTyped[T]{binding.NewUntyped()} }
func BindTyped[T comparable](v *T) Typed[T] { return &boundExternalTyped[T]{binding.BindUntyped(v)} }

type boundTyped[T comparable] struct {
	binding.Untyped
}

func (t *boundTyped[T]) Get() (T, error) {
	v, err := t.Untyped.Get()
	return v.(T), err
}
func (t *boundTyped[T]) Set(v T) error { return t.Untyped.Set(v) }

type boundExternalTyped[T comparable] struct {
	binding.ExternalUntyped
}

func (t *boundExternalTyped[T]) Get() (T, error) {
	v, err := t.ExternalUntyped.Get()
	return v.(T), err
}
func (t *boundExternalTyped[T]) Set(v T) error { return t.ExternalUntyped.Set(v) }

func (t *boundExternalTyped[T]) Reload() error { return t.ExternalUntyped.Reload() }

type TypedList[T any] interface {
	binding.DataList

	Get() ([]T, error)
	Set(list []T) error
	GetValue(index int) (T, error)
	SetValue(index int, value T) error
	Append(value T) error
	Prepend(value T) error
	Remove(value T) error
}

type ExternalTypedList[T any] interface {
	TypedList[T]

	Reload() error
}

func NewTypedList[T any]() TypedList[T] { return &boundTypedList[T]{binding.NewUntypedList()} }
func BindTypedList[T any](v *[]T) ExternalTypedList[T] {
	s := makeAnySlice(*v)
	return &boundTypedList[T]{binding.BindUntypedList(&s)}
}

type boundTypedList[T any] struct {
	binding.UntypedList
}

func (l *boundTypedList[T]) Get() ([]T, error) {
	ul, err := l.UntypedList.Get()
	tl := make([]T, len(ul))
	for i, v := range ul {
		tl[i] = v.(T)
	}
	return tl, err
}
func (l *boundTypedList[T]) Set(tl []T) error {
	return l.UntypedList.Set(makeAnySlice(tl))
}

func (l *boundTypedList[T]) GetValue(i int) (T, error) {
	v, err := l.UntypedList.GetValue(i)
	if v == nil {
		var v T
		return v, err
	}
	return v.(T), err
}
func (l *boundTypedList[T]) SetValue(i int, v T) error { return l.UntypedList.SetValue(i, v) }

func (l *boundTypedList[T]) Append(v T) error  { return l.UntypedList.Append(v) }
func (l *boundTypedList[T]) Prepend(v T) error { return l.UntypedList.Prepend(v) }
func (l *boundTypedList[T]) Remove(v T) error  { return l.UntypedList.Remove(v) }

func (l *boundTypedList[T]) Reload() error {
	return l.UntypedList.(binding.ExternalUntypedList).Reload()
}

func makeAnySlice[T any](typed []T) []any {
	ul := make([]any, len(typed))
	for i, v := range typed {
		ul[i] = v
	}
	return ul
}
