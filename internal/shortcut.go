package internal

import "fyne.io/fyne/v2"

type ShortcutNew struct{}

func (ShortcutNew) ShortcutName() string  { return "New" }
func (ShortcutNew) Key() fyne.KeyName     { return fyne.KeyN }
func (ShortcutNew) Mod() fyne.KeyModifier { return fyne.KeyModifierShortcutDefault }

type ShortcutOpen struct{}

func (ShortcutOpen) ShortcutName() string  { return "Open" }
func (ShortcutOpen) Key() fyne.KeyName     { return fyne.KeyO }
func (ShortcutOpen) Mod() fyne.KeyModifier { return fyne.KeyModifierShortcutDefault }

type ShortcutSave struct{}

func (ShortcutSave) ShortcutName() string  { return "Save" }
func (ShortcutSave) Key() fyne.KeyName     { return fyne.KeyS }
func (ShortcutSave) Mod() fyne.KeyModifier { return fyne.KeyModifierShortcutDefault }

type ShortcutClose struct{}

func (ShortcutClose) ShortcutName() string  { return "Close" }
func (ShortcutClose) Key() fyne.KeyName     { return fyne.KeyW }
func (ShortcutClose) Mod() fyne.KeyModifier { return fyne.KeyModifierShortcutDefault }
