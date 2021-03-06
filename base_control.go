package clui

import (
	term "github.com/nsf/termbox-go"
)

// BaseControl is a base for all visible controls.
// Every new control must inherit it or implement
// the same set of methods
type BaseControl struct {
	title         string
	x, y          int
	width, height int
	minW, minH    int
	scale         int
	fg, bg        term.Attribute
	fgActive      term.Attribute
	bgActive      term.Attribute
	tabSkip       bool
	disabled      bool
	align         Align
	parent        Control
	inactive      bool
	modal         bool
	padX, padY    int
	gapX, gapY    int
	pack          PackType
	children      []Control
}

func (c *BaseControl) Title() string {
	return c.title
}

func (c *BaseControl) SetTitle(title string) {
	c.title = title
}

func (c *BaseControl) Size() (widht int, height int) {
	return c.width, c.height
}

func (c *BaseControl) SetSize(width, height int) {
	if width < c.minW {
		width = c.minW
	}
	if height < c.minH {
		height = c.minH
	}

	if height != c.height || width != c.width {
		c.height = height
		c.width = width
	}
}

func (c *BaseControl) Pos() (x int, y int) {
	return c.x, c.y
}

func (c *BaseControl) SetPos(x, y int) {
	c.x = x
	c.y = y
}

func (c *BaseControl) applyConstraints() {
	ww, hh := c.width, c.height
	if ww < c.minW {
		ww = c.minW
	}
	if hh < c.minH {
		hh = c.minH
	}
	if hh != c.height || ww != c.width {
		c.SetSize(ww, hh)
	}
}

func (c *BaseControl) Constraints() (minw int, minh int) {
	return c.minW, c.minH
}

func (c *BaseControl) SetConstraints(minw, minh int) {
	c.minW = minw
	c.minH = minh
	c.applyConstraints()
}

func (c *BaseControl) Active() bool {
	return !c.inactive
}

func (c *BaseControl) SetActive(active bool) {
	c.inactive = !active
}

func (c *BaseControl) TabStop() bool {
	return !c.tabSkip
}

func (c *BaseControl) SetTabStop(tabstop bool) {
	c.tabSkip = !tabstop
}

func (c *BaseControl) Enabled() bool {
	return !c.disabled
}

func (c *BaseControl) SetEnabled(enabled bool) {
	c.disabled = !enabled
}

func (c *BaseControl) Parent() Control {
	return c.parent
}

func (c *BaseControl) SetParent(parent Control) {
	if c.parent == nil {
		c.parent = parent
	}
}

func (c *BaseControl) Modal() bool {
	return c.modal
}

func (c *BaseControl) SetModal(modal bool) {
	c.modal = modal
}

func (c *BaseControl) Paddings() (px int, py int) {
	return c.padX, c.padY
}

func (c *BaseControl) SetPaddings(px, py int) {
	if px >= 0 {
		c.padX = px
	}
	if py >= 0 {
		c.padY = py
	}
}

func (c *BaseControl) Gaps() (dx int, dy int) {
	return c.gapX, c.gapY
}

func (c *BaseControl) SetGaps(dx, dy int) {
	if dx >= 0 {
		c.gapX = dx
	}
	if dy >= 0 {
		c.gapY = dy
	}
}

func (c *BaseControl) Pack() PackType {
	return c.pack
}

func (c *BaseControl) SetPack(pack PackType) {
	c.pack = pack
}

func (c *BaseControl) Scale() int {
	return c.scale
}

func (c *BaseControl) SetScale(scale int) {
	if scale >= 0 {
		c.scale = scale
	}
}

func (c *BaseControl) Align() Align {
	return c.align
}

func (c *BaseControl) SetAlign(align Align) {
	c.align = align
}

func (c *BaseControl) TextColor() term.Attribute {
	return c.fg
}

func (c *BaseControl) SetTextColor(clr term.Attribute) {
	c.fg = clr
}

func (c *BaseControl) BackColor() term.Attribute {
	return c.bg
}

func (c *BaseControl) SetBackColor(clr term.Attribute) {
	c.bg = clr
}

func (c *BaseControl) ResizeChildren() {
	if len(c.children) == 0 {
		return
	}

	fullWidth := c.width - 2*c.padX
	fullHeight := c.height - 2*c.padY
	if c.pack == Horizontal {
		fullWidth -= (len(c.children) - 1) * c.gapX
	} else {
		fullHeight -= (len(c.children) - 1) * c.gapY
	}

	totalSc := c.ChildrenScale()
	minWidth := 0
	minHeight := 0
	for _, child := range c.children {
		cw, ch := child.MinimalSize()
		if c.pack == Horizontal {
			minWidth += cw
		} else {
			minHeight += ch
		}
	}

	aStep := 0
	diff := fullWidth - minWidth
	if c.pack == Vertical {
		diff = fullHeight - minHeight
	}
	if totalSc > 0 {
		aStep = int(float32(diff) / float32(totalSc))
	}

	for _, ctrl := range c.children {
		tw, th := ctrl.MinimalSize()
		sc := ctrl.Scale()
		d := int(ctrl.Scale() * aStep)
		if c.pack == Horizontal {
			if sc != 0 {
				if sc == totalSc {
					tw += diff
					d = diff
				} else {
					tw += d
				}
			}
			th = fullHeight
		} else {
			if sc != 0 {
				if sc == totalSc {
					th += diff
					d = diff
				} else {
					th += d
				}
			}
			tw = fullWidth
		}
		diff -= d
		totalSc -= sc

		ctrl.SetSize(tw, th)
		ctrl.ResizeChildren()
	}
}

func (c *BaseControl) AddChild(control Control) {
	if c.children == nil {
		c.children = make([]Control, 1)
		c.children[0] = control
	} else {
		if c.ChildExists(control) {
			panic("Double adding a child")
		}

		c.children = append(c.children, control)
	}

	var ctrl Control
	var mainCtrl Control
	ctrl = c
	for ctrl != nil {
		ww, hh := ctrl.MinimalSize()
		cw, ch := ctrl.Size()
		if ww > cw || hh > ch {
			if ww > cw {
				cw = ww
			}
			if hh > ch {
				ch = hh
			}
			ctrl.SetConstraints(cw, ch)
		}

		if ctrl.Parent() == nil {
			mainCtrl = ctrl
		}
		ctrl = ctrl.Parent()
	}

	if mainCtrl != nil {
		mainCtrl.ResizeChildren()
		mainCtrl.PlaceChildren()
	}
}

func (c *BaseControl) Children() []Control {
	child := make([]Control, len(c.children))
	copy(child, c.children)
	return child
}

func (c *BaseControl) ChildExists(control Control) bool {
	if len(c.children) == 0 {
		return false
	}

	for _, ctrl := range c.children {
		if ctrl == control {
			return true
		}
	}

	return false
}

func (c *BaseControl) ChildrenScale() int {
	if len(c.children) == 0 {
		return c.scale
	}

	total := 0
	for _, ctrl := range c.children {
		total += ctrl.Scale()
	}

	return total
}

func (c *BaseControl) MinimalSize() (w int, h int) {
	if len(c.children) == 0 {
		return c.minW, c.minH
	}

	totalX := 2 * c.padX
	totalY := 2 * c.padY

	if c.pack == Vertical {
		totalY += (len(c.children) - 1) * c.gapY
	} else {
		totalX += (len(c.children) - 1) * c.gapX
	}

	for _, ctrl := range c.children {
		ww, hh := ctrl.MinimalSize()
		if c.pack == Vertical {
			totalY += hh
			if ww+2*c.padX > totalX {
				totalX = ww + 2*c.padX
			}
		} else {
			totalX += ww
			if hh+2*c.padY > totalY {
				totalY = hh + 2*c.padY
			}
		}
	}

	if totalX < c.minW {
		totalX = c.minW
	}
	if totalY < c.minH {
		totalY = c.minH
	}

	return totalX, totalY
}

func (c *BaseControl) Draw() {
	panic("BaseControl Draw Called")
}

func (c *BaseControl) DrawChildren() {
	PushClip()
	defer PopClip()

	SetClipRect(c.x+c.padX, c.y+c.padY, c.width-2*c.padX, c.height-2*c.padY)

	for _, child := range c.children {
		child.Draw()
	}
}

func (c *BaseControl) HitTest(x, y int) HitResult {
	if x > c.x && x < c.x+c.width-1 &&
		y > c.y && y < c.y+c.height-1 {
		return HitInside
	}

	if (x == c.x || x == c.x+c.width-1) &&
		y >= c.y && y < c.y+c.height {
		return HitBorder
	}

	if (y == c.y || y == c.y+c.height-1) &&
		x >= c.x && x < c.x+c.width {
		return HitBorder
	}

	return HitOutside
}

func (c *BaseControl) ProcessEvent(ev Event) bool {
	return SendEventToChild(c, ev)
}

func (c *BaseControl) PlaceChildren() {
	if c.children == nil || len(c.children) == 0 {
		return
	}

	xx, yy := c.x+c.padX, c.y+c.padY
	for _, ctrl := range c.children {
		ctrl.SetPos(xx, yy)

		ww, hh := ctrl.Size()
		if c.pack == Vertical {
			yy += c.gapY + hh
		} else {
			xx += c.gapX + ww
		}

		ctrl.PlaceChildren()
	}
}

// ActiveColors return the attrubutes for the controls when it
// is active: text and background colors
func (c *BaseControl) ActiveColors() (term.Attribute, term.Attribute) {
	return c.fgActive, c.bgActive
}

// SetActiveTextColor changes text color of the active control
func (c *BaseControl) SetActiveTextColor(clr term.Attribute) {
	c.fgActive = clr
}

// SetActiveBackColor changes background color of the active control
func (c *BaseControl) SetActiveBackColor(clr term.Attribute) {
	c.bgActive = clr
}
