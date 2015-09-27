package main

import "github.com/go-gl/glfw/v3.1/glfw"

// newMenu creates a hidden menu, only when you first call show will it be
// visible.
func newMenu(elements ...guiElement) *menu {
	const margin = 50
	m := &menu{elements: elements, visible: false}
	l, r, t, b := 0, 0, 0, 0
	if len(elements) > 0 {
		bounds := elements[0].bounds()
		l, t = bounds.x, bounds.y
		r, b = l+bounds.w, t+bounds.h
		for _, elem := range elements {
			bounds = elem.bounds()
			if bounds.x < l {
				l = bounds.x
			}
			if bounds.y < t {
				t = bounds.y
			}
			if right := bounds.x + bounds.w; right > r {
				r = right
			}
			if bottom := bounds.y + bounds.h; bottom > b {
				b = bottom
			}
		}
	}
	m.rect = rect{l - margin, t - margin, r - l + 2*margin, b - t + 2*margin}
	return m
}

type guiElement interface {
	bounds() rect
	draw(g *graphics)
	mouseMovedTo(x, y int)
	click(x, y int) (actionID int)
	runeTyped(r rune)
	keyPressed(key glfw.Key)
}

type menu struct {
	rect
	elements []guiElement
	visible  bool
}

func (m *menu) draw(g *graphics) {
	if !m.visible {
		return
	}

	color := [4]float32{1, 0.2, 0.2, 0.5}
	g.rect(m.x, m.y, m.w, m.h, color)
	for _, elem := range m.elements {
		elem.draw(g)
	}
}

func (m *menu) mouseMovedTo(x, y int) {
	if !m.visible {
		return
	}

	for _, elem := range m.elements {
		elem.mouseMovedTo(x, y)
	}
}

func (m *menu) click(x, y int) (actionID int) {
	if !m.visible {
		return -1
	}

	for _, elem := range m.elements {
		if action := elem.click(x, y); action != -1 {
			return action
		}
	}
	return -1
}

func (m *menu) runeTyped(r rune) {
	if !m.visible {
		return
	}

	for _, elem := range m.elements {
		elem.runeTyped(r)
	}
}

func (m *menu) keyPressed(key glfw.Key) {
	for _, elem := range m.elements {
		elem.keyPressed(key)
	}
}

func (m *menu) hide() { m.visible = false }
func (m *menu) show() { m.visible = true }

// some constants that are equal for all menu elements

var (
	menuColdBackColor      = [4]float32{0.8, 0, 0, 0.7}
	menuHotBackColor       = [4]float32{0.8, 0, 0, 1}
	menuFontColor          = [4]float32{1, 1, 1, 1}
	checkBoxCheckedColor   = [4]float32{0.5, 1, 0.5, 1}
	checkBoxUncheckedColor = [4]float32{1, 0.5, 0.5, 1}
)

// button

func newButton(text string, bounds rect, actionID int) *button {
	return &button{
		rect:   bounds,
		text:   text,
		action: actionID,
	}
}

type button struct {
	rect
	text   string
	hot    bool
	action int
}

func (b *button) bounds() rect { return b.rect }

func (b *button) draw(g *graphics) {
	color := menuColdBackColor
	if b.hot {
		color = menuHotBackColor
	}
	g.rect(b.x, b.y, b.w, b.h, color)
	g.writeTextLineCenteredInRect(b.text, b.rect, menuFontColor)
}

func (b *button) mouseMovedTo(x, y int) {
	b.hot = b.contains(x, y)
}

func (b *button) click(x, y int) int {
	if b.contains(x, y) {
		return b.action
	}
	return -1
}

func (b *button) runeTyped(r rune)        {}
func (b *button) keyPressed(key glfw.Key) {}

// text field

func newTextField(bounds rect, font textSizer) *textBox {
	return &textBox{
		rect: bounds,
		font: font,
	}
}

type textSizer interface {
	TextSize(text string) (w, h int)
}

type textBox struct {
	rect
	hot        bool
	text       string
	textLength int
	font       textSizer
}

func (t *textBox) bounds() rect { return t.rect }

func (t *textBox) click(x, y int) int {
	t.hot = t.contains(x, y)
	return -1
}

func (t *textBox) draw(g *graphics) {
	color := menuColdBackColor
	if t.hot {
		color = menuHotBackColor
	}
	g.rect(t.x, t.y, t.w, t.h, color)
	g.writeTextLineCenteredInRect(t.text, t.rect, menuFontColor)
}

func (t *textBox) mouseMovedTo(x, y int) {}

func (t *textBox) runeTyped(r rune) {
	if t.hot {
		newText := t.text + string(r)
		w, _ := t.font.TextSize(newText)
		if w < t.rect.w {
			t.text += string(r)
			t.textLength++
		}
	}
}

func (t *textBox) keyPressed(key glfw.Key) {
	if t.hot {
		if key == glfw.KeyBackspace && len(t.text) > 0 {
			var last int
			for i := range t.text {
				last = i
			}
			t.text = t.text[:last]
			t.textLength--
		}
	}
}

// check box

func newCheckBox(text string, bounds rect, font textSizer, id int) *checkBox {
	const margin = 10
	size := bounds.h - 2*margin
	checkRect := rect{bounds.x + margin, bounds.y + margin, size, size}
	textW, textH := font.TextSize(text)
	textRect := rect{
		checkRect.x + checkRect.w + 3*margin,
		checkRect.y + (checkRect.h-textH)/2,
		textW,
		textH,
	}
	return &checkBox{
		rect:      bounds,
		checkRect: checkRect,
		textRect:  textRect,
		text:      text,
		id:        id,
	}
}

type checkBox struct {
	rect
	checkRect rect
	textRect  rect
	text      string
	id        int
	checked   bool
}

func (c *checkBox) bounds() rect {
	return c.rect
}

func (c *checkBox) draw(g *graphics) {
	g.rect(c.x, c.y, c.w, c.h, menuColdBackColor)
	checkColor := checkBoxUncheckedColor
	if c.checked {
		checkColor = checkBoxCheckedColor
	}
	b := c.checkRect
	g.rect(b.x, b.y, b.w, b.h, checkColor)
	g.writeTextLineCenteredInRect(c.text, c.textRect, menuFontColor)
}

func (c *checkBox) mouseMovedTo(x, y int) {}

func (c *checkBox) click(x, y int) (actionID int) {
	if c.checkRect.contains(x, y) {
		c.checked = !c.checked
		if c.checked {
			return c.id
		}
	}
	return -1
}

func (c *checkBox) runeTyped(r rune)        {}
func (c *checkBox) keyPressed(key glfw.Key) {}

// check box group

func newCheckBoxGroup(boxes ...*checkBox) *checkBoxGroup {
	for _, box := range boxes {
		box.checked = false
	}
	if len(boxes) > 0 {
		boxes[0].checked = true
	}
	return &checkBoxGroup{boxes: boxes}
}

type checkBoxGroup struct {
	boxes        []*checkBox
	checkedIndex int
}

func (group *checkBoxGroup) bounds() rect {
	if len(group.boxes) == 0 {
		return rect{}
	}
	box := group.boxes[0]
	x, y, r, b := box.x, box.y, box.w, box.h
	for i := 1; i < len(group.boxes); i++ {
		box = group.boxes[i]
		if box.x < x {
			x = box.x
		}
		if box.y < y {
			y = box.y
		}
		if right := box.x + box.w; right > r {
			r = right
		}
		if bottom := box.y + box.h; bottom > b {
			b = bottom
		}
	}
	return rect{x, y, r - x, b - y}
}

func (group *checkBoxGroup) draw(g *graphics) {
	for _, box := range group.boxes {
		box.draw(g)
	}
}

func (group *checkBoxGroup) mouseMovedTo(x, y int) {
	for _, box := range group.boxes {
		box.mouseMovedTo(x, y)
	}
}

func (group *checkBoxGroup) click(x, y int) (actionID int) {
	for index, box := range group.boxes {
		wasChecked := box.checked
		id := box.click(x, y)
		if wasChecked != box.checked {
			// check changed
			if !box.checked {
				// check boxes in a group can not be unchecked by clicking on
				// them, so if this was the case, check it again
				box.checked = true
			} else {
				// new box was checked, uncheck the last one
				group.boxes[group.checkedIndex].checked = false
				group.checkedIndex = index
				return id
			}
		}
	}
	return -1
}

func (group *checkBoxGroup) runeTyped(r rune) {
	for _, box := range group.boxes {
		box.runeTyped(r)
	}
}

func (group *checkBoxGroup) keyPressed(key glfw.Key) {
	for _, box := range group.boxes {
		box.keyPressed(key)
	}
}
