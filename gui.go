package main

import "github.com/go-gl/glfw/v3.1/glfw"

// newMenu creates a hidden menu, only when you first call show will it be
// visible.
func newMenu(margin int, elements ...guiElement) *menu {
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
	color := [4]float32{0.8, 0, 0, 0.7}
	if b.hot {
		color[3] = 1
	}
	g.rect(b.x, b.y, b.w, b.h, color)
	g.writeTextLineCenteredInRect(b.text, b.rect, [4]float32{1, 1, 1, 1})
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

func newTextField(bounds rect, font textSizer, id int) *textBox {
	return &textBox{
		rect: bounds,
		id:   id,
		font: font,
	}
}

type textSizer interface {
	TextSize(text string) (w, h int)
}

type textBox struct {
	rect
	id         int
	hot        bool
	text       string
	textLength int
	font       textSizer
}

func (t *textBox) bounds() rect { return t.rect }

func (t *textBox) click(x, y int) int {
	t.hot = t.contains(x, y)
	println(t.hot)
	return -1
}

func (t *textBox) draw(g *graphics) {
	color := [4]float32{0.8, 0, 0, 0.7}
	if t.hot {
		color[3] = 1
	}
	g.rect(t.x, t.y, t.w, t.h, color)
	g.writeTextLineCenteredInRect(t.text, t.rect, [4]float32{1, 1, 1, 1})
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
