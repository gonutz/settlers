package main

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/gonutz/settlers/lang"
)

var (
	menuHotBackColor       = [4]float32{0.71, 0.4, 0.31, 1}
	menuColdBackColor      = [4]float32{0.7, 0.27, 0.137, 1}
	menuFontColor          = [4]float32{1, 1, 1, 1}
	menuDisabledFontColor  = [4]float32{0.6, 0.6, 0.6, 1}
	menuColdFontColor      = [4]float32{0.7, 0.7, 0.7, 1}
	checkBoxCheckedColor   = [4]float32{0.5, 1, 0.5, 1}
	checkBoxUncheckedColor = [4]float32{1, 0.5, 0.5, 1}
)

//var (
//	menuColdBackColor      = [4]float32{0.8, 0, 0, 1}
//	menuHotBackColor       = [4]float32{0.6, 0, 0, 1}
//	menuFontColor          = [4]float32{1, 1, 1, 1}
//	menuDisabledFontColor  = [4]float32{0.6, 0.6, 0.6, 1}
//	menuColdFontColor      = [4]float32{0.7, 0.7, 0.7, 1}
//	checkBoxCheckedColor   = [4]float32{0.5, 1, 0.5, 1}
//	checkBoxUncheckedColor = [4]float32{1, 0.5, 0.5, 1}
//)

type guiElement interface {
	bounds() rect
	setBounds(rect)
	draw(*graphics)
	mouseMovedTo(x, y int)
	// click returns the ID of the activated action or -1 if none was activated.
	click(x, y int) (actionID int)
	runeTyped(rune)
	keyPressed(glfw.Key)
}

func boundingBox(bucket bounderBucket) rect {
	if bucket.Len() == 0 {
		return rect{}
	}

	bounds := bucket.At(0).bounds()
	l, t, r, b := bounds.x, bounds.y, bounds.x, bounds.y
	for i := 1; i < bucket.Len(); i++ {
		bounds = bucket.At(i).bounds()
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
	return rect{l, t, r - l, b - t}
}

type bounderBucket interface {
	Len() int
	At(index int) bounder
}

type bounder interface {
	bounds() rect
}

// windows are visible by default

func newWindow(bounds rect, layout layout, elems ...guiElement) *window {
	w := &window{newComposite(elems...), layout, true}
	for _, e := range elems {
		layout.addElement(e)
	}
	layout.relayout(bounds)
	return w
}

type window struct {
	*composite
	layout  layout
	visible bool
}

func (w *window) setBounds(bounds rect) {
	w.composite.setBounds(bounds)
	w.layout.relayout(bounds)
}

func (w *window) setVisible(v bool) { w.visible = v }

func (w *window) draw(g *graphics) {
	if !w.visible {
		return
	}
	w.composite.draw(g)
}

func (w *window) mouseMovedTo(x, y int) {
	if !w.visible {
		return
	}
	w.composite.mouseMovedTo(x, y)
}

func (w *window) click(x, y int) (actionID int) {
	if !w.visible {
		return -1
	}
	return w.composite.click(x, y)
}

func (w *window) runeTyped(r rune) {
	if !w.visible {
		return
	}
	w.composite.runeTyped(r)
}

func (w *window) keyPressed(key glfw.Key) {
	if !w.visible {
		return
	}
	w.composite.keyPressed(key)
}

// composite

func newComposite(elems ...guiElement) *composite {
	return &composite{elems}
}

type composite struct {
	elems []guiElement
}

func (c *composite) bounds() rect {
	return boundingBox(guiElementsToBounderBucket(c.elems))
}

func (c *composite) setBounds(bounds rect) {
	b := c.bounds()
	dx, dy := bounds.x-b.x, bounds.y-b.y
	for _, e := range c.elems {
		e.setBounds(e.bounds().moveBy(dx, dy))
	}
}

func (c *composite) draw(g *graphics) {
	for _, e := range c.elems {
		e.draw(g)
	}
}

func (c *composite) mouseMovedTo(x, y int) {
	for _, e := range c.elems {
		e.mouseMovedTo(x, y)
	}
}

func (c *composite) click(x, y int) (actionID int) {
	for _, e := range c.elems {
		if id := e.click(x, y); id != -1 {
			return id
		}
	}
	return -1
}

func (c *composite) runeTyped(r rune) {
	for _, e := range c.elems {
		e.runeTyped(r)
	}
}

func (c *composite) keyPressed(key glfw.Key) {
	for _, e := range c.elems {
		e.keyPressed(key)
	}
}

type guiElementsToBounderBucket []guiElement

func (slice guiElementsToBounderBucket) Len() int         { return len(slice) }
func (slice guiElementsToBounderBucket) At(i int) bounder { return slice[i] }

// button

func newButton(textID lang.Item, bounds rect, actionID int) *button {
	return &button{
		rect:   bounds,
		textID: textID,
		action: actionID,
	}
}

type button struct {
	rect
	textID lang.Item
	hot    bool
	action int
}

func (b *button) bounds() rect          { return b.rect }
func (b *button) setBounds(bounds rect) { b.rect = bounds }

func (b *button) draw(g *graphics) {
	color := menuColdBackColor
	if b.hot {
		color = menuHotBackColor
	}
	g.rect(b.x, b.y, b.w, b.h, color)
	g.writeTextLineCenteredInRect(lang.Get(b.textID), b.rect, menuFontColor)
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

func (*button) runeTyped(rune)      {}
func (*button) keyPressed(glfw.Key) {}

// text box

func newTextBox(captionID lang.Item, bounds rect, font textSizer) *textBox {
	return &textBox{
		rect:      bounds,
		captionID: captionID,
		font:      font,
	}
}

type textBox struct {
	rect
	hot               bool
	captionID         lang.Item
	text              string
	textLengthInRunes int
	font              textSizer
	captionRect       rect
	textRect          rect
	disabled          bool
}

func (t *textBox) bounds() rect          { return t.rect }
func (t *textBox) setBounds(bounds rect) { t.rect = bounds }

func (t *textBox) click(x, y int) int {
	if t.disabled {
		return -1
	}
	t.hot = t.contains(x, y)
	return -1
}

func (t *textBox) draw(g *graphics) {
	t.recalcRects()
	fontColor := menuFontColor
	if t.disabled {
		fontColor = menuDisabledFontColor
	}

	g.rect(t.x, t.y, t.w, t.h, menuColdBackColor)
	g.writeTextLineCenteredInRect(lang.Get(t.captionID), t.captionRect, fontColor)
	if t.hot {
		g.rect(t.textRect.x, t.textRect.y, t.textRect.w, t.textRect.h, menuHotBackColor)
	}
	g.writeTextLineCenteredInRect(t.text, t.textRect, fontColor)
}

func (t *textBox) recalcRects() {
	const margin = 20
	captionW, _ := t.font.TextSize(lang.Get(t.captionID))
	t.captionRect = rect{t.x + margin, t.y + margin, captionW, t.h - 2*margin}
	t.textRect = rect{t.x + 2*margin + captionW, t.y, t.w - 2*margin - captionW, t.h}
}

func (t *textBox) mouseMovedTo(x, y int) {}

func (t *textBox) runeTyped(r rune) {
	if t.disabled {
		return
	}
	if t.hot {
		newText := t.text + string(r)
		w, _ := t.font.TextSize(newText)
		if w < t.textRect.w {
			t.text += string(r)
			t.textLengthInRunes++
		}
	}
}

func (t *textBox) keyPressed(key glfw.Key) {
	if t.disabled {
		return
	}
	if t.hot {
		if key == glfw.KeyBackspace && len(t.text) > 0 {
			var last int
			for i := range t.text {
				last = i
			}
			t.text = t.text[:last]
			t.textLengthInRunes--
		}
		if key == glfw.KeyEnter || key == glfw.KeyKPEnter {
			t.hot = false
		}
	}
}

func (t *textBox) setEnabled(enabled bool) {
	t.disabled = !enabled
	if t.disabled {
		t.hot = false
	}
}

type textSizer interface {
	TextSize(text string) (w, h int)
}

// check box

func newCheckBox(textID lang.Item, bounds rect, id int) *checkBox {
	cb := &checkBox{
		rect:   bounds,
		textID: textID,
		id:     id,
	}
	cb.layout()
	return cb
}

func (c *checkBox) layout() {
	const margin = 10
	size := c.rect.h - 2*margin
	c.checkRect = rect{c.rect.x + margin, c.rect.y + margin, size, size}
	c.textX, c.textY = c.checkRect.x+c.checkRect.w+3*margin, c.rect.y+c.rect.h/2
}

type checkBox struct {
	rect
	checkRect        rect
	textX, textY     int
	textID           lang.Item
	id               int
	checked          bool
	checkChangeEvent func(bool)
}

func (c *checkBox) bounds() rect { return c.rect }

func (c *checkBox) setBounds(bounds rect) {
	c.rect = bounds
	c.layout()
}

func (c *checkBox) draw(g *graphics) {
	g.rect(c.x, c.y, c.w, c.h, menuColdBackColor)
	checkColor := checkBoxUncheckedColor
	if c.checked {
		checkColor = checkBoxCheckedColor
	}
	b := c.checkRect
	g.rect(b.x, b.y, b.w, b.h, checkColor)
	g.writeLeftAlignedVerticallyCenteredAt(lang.Get(c.textID), c.textX, c.textY, menuFontColor)
}

func (c *checkBox) mouseMovedTo(x, y int) {}

func (c *checkBox) click(x, y int) (actionID int) {
	if c.contains(x, y) {
		c.setChecked(!c.checked)
		if c.checked {
			return c.id
		}
	}
	return -1
}

func (c *checkBox) setChecked(checked bool) {
	if c.checked != checked {
		c.checked = checked
		if c.checkChangeEvent != nil {
			c.checkChangeEvent(checked)
		}
	}
}

func (c *checkBox) runeTyped(r rune)        {}
func (c *checkBox) keyPressed(key glfw.Key) {}

func (c *checkBox) onCheckChange(action func(bool)) {
	c.checkChangeEvent = action
	if action != nil {
		action(c.checked)
	}
}

// check box group

func newCheckBoxGroup(boxes ...*checkBox) *checkBoxGroup {
	for _, box := range boxes {
		box.setChecked(false)
	}
	if len(boxes) > 0 {
		boxes[0].setChecked(true)
	}
	layout := newTopLeftLayout()
	elems := make([]guiElement, len(boxes))
	for i, b := range boxes {
		elems[i] = b
		layout.addElement(b)
	}
	layout.relayout(rect{})
	return &checkBoxGroup{
		composite: newComposite(elems...),
		boxes:     boxes,
	}
}

type checkBoxGroup struct {
	*composite
	boxes        []*checkBox
	checkedIndex int
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
				box.setChecked(true)
			} else {
				// new box was checked, uncheck the last one
				group.boxes[group.checkedIndex].setChecked(false)
				group.checkedIndex = index
				return id
			}
		}
	}
	return -1
}

// tab sheet

func newTabSheet(font textSizer, tabs ...*tab) *tabSheet {
	sheet := &tabSheet{tabs: tabs, font: font}
	sheet.relayout()
	return sheet
}

func (s *tabSheet) relayout() {
	s.visibleTabs = make([]*tab, 0, len(s.tabs))
	elems := make([]guiElement, 0, len(s.tabs))
	lastVisibleIndex := 0
	for i, tab := range s.tabs {
		if tab.visible {
			s.visibleTabs = append(s.visibleTabs, tab)
			elems = append(elems, tab.content)
			lastVisibleIndex = i
		}
	}
	if !s.tabs[s.activeIndex].visible {
		s.activeIndex = lastVisibleIndex
	}
	bounds := boundingBox(guiElementsToBounderBucket(elems))

	captionBounds := make([]rect, len(s.visibleTabs))
	captionW := bounds.w / len(s.visibleTabs)

	x := bounds.x
	captionH := 0
	for i, tab := range s.visibleTabs {
		_, h := s.font.TextSize(tab.title)
		h += 25
		captionBounds[i] = rect{x, bounds.y - h, captionW, h}
		x += captionW
		if h > captionH {
			captionH = h
		}
	}

	const margin = 0
	s.rect = rect{
		bounds.x - margin,
		bounds.y - captionH - margin,
		bounds.w + 2*margin,
		bounds.h + captionH + 2*margin,
	}
	s.captionBounds = captionBounds
	s.captionH = captionH
}

type tabSheet struct {
	rect
	font          textSizer
	tabs          []*tab
	visibleTabs   []*tab
	activeIndex   int
	captionBounds []rect
	captionH      int
}

func (s *tabSheet) bounds() rect { return s.rect }

func (s *tabSheet) setBounds(bounds rect) {
	dx, dy := bounds.x-s.x, bounds.y-s.y
	for _, tab := range s.tabs {
		tab.content.setBounds(tab.content.bounds().moveBy(dx, dy))
	}
	s.rect = bounds
	s.relayout()
}

func (s *tabSheet) draw(g *graphics) {
	activeColor := s.visibleTabs[s.activeIndex].color
	activeColor[3] *= 0.85
	g.rect(s.x, s.y+s.captionH, s.w, s.h-s.captionH, activeColor)
	for i, tab := range s.visibleTabs {
		b := s.captionBounds[i]
		color := tab.color
		if i == s.activeIndex {
			color = activeColor
		}
		g.rect(b.x, b.y, b.w, b.h, color)
		fontColor := menuColdFontColor
		if i == s.activeIndex {
			fontColor = menuFontColor
		}
		g.writeTextLineCenteredInRect(tab.title, b, fontColor)
	}
	s.visibleTabs[s.activeIndex].content.draw(g)
}

func (s *tabSheet) mouseMovedTo(x, y int) {
	s.visibleTabs[s.activeIndex].content.mouseMovedTo(x, y)
}

func (s *tabSheet) click(x, y int) (actionID int) {
	for i, b := range s.captionBounds {
		if b.contains(x, y) {
			s.activeIndex = i
			return -1
		}
	}
	return s.visibleTabs[s.activeIndex].content.click(x, y)
}

func (s *tabSheet) runeTyped(r rune) {
	s.visibleTabs[s.activeIndex].content.runeTyped(r)
}

func (s *tabSheet) keyPressed(key glfw.Key) {
	s.visibleTabs[s.activeIndex].content.keyPressed(key)
}

// tab

func newTab(title string, color [4]float32, content guiElement, visible bool) *tab {
	return &tab{title, content, color, visible}
}

type tab struct {
	title   string
	content guiElement
	color   [4]float32
	visible bool
}

// spacer

func newSpacer(bounds rect) *spacer {
	return &spacer{bounds}
}

type spacer struct {
	rect
}

func (s *spacer) bounds() rect               { return s.rect }
func (s *spacer) setBounds(bounds rect)      { s.rect = bounds }
func (spacer) draw(*graphics)                {}
func (spacer) mouseMovedTo(x, y int)         {}
func (spacer) click(x, y int) (actionID int) { return -1 }
func (spacer) runeTyped(rune)                {}
func (spacer) keyPressed(glfw.Key)           {}
