package main

type layout interface {
	addElement(bounded)
	relayout()
}

type bounded interface {
	bounds() rect
	setBounds(rect)
}

// layoutBase implements addElement with a slice. This way not every layout has
// to re-implement it.
type layoutBase struct {
	items []bounded
}

func (l *layoutBase) addElement(item bounded) {
	l.items = append(l.items, item)
}

// center layout aligns the bounding box of its elements so its center lies on a
// given point

func newCenterLayout(centerX, centerY int) *centerLayout {
	return &centerLayout{new(layoutBase), centerX, centerY}
}

type centerLayout struct {
	*layoutBase
	cx, cy int
}

func (l *centerLayout) relayout() {
	bounds := boundingBox(boundingBoxableItems(l.items))
	boundsCx, boundsCy := bounds.x+bounds.w/2, bounds.y+bounds.h/2
	dx, dy := l.cx-boundsCx, l.cy-boundsCy
	for _, item := range l.items {
		b := item.bounds()
		item.setBounds(rect{b.x + dx, b.y + dy, b.w, b.h})
	}
}

type boundingBoxableItems []bounded

func (slice boundingBoxableItems) Len() int         { return len(slice) }
func (slice boundingBoxableItems) At(i int) bounder { return slice[i] }

// The vertical flow layout puts all elements directly underneath each other and
// centers them horizontally on each row. The bounding rect around the result
// starts at 0,0 so the first item will start at y=0 and x=0+offset where the
// offset depends on the maximum width of all elements.

func newVerticalFlowLayout() *verticalFlowLayout {
	return &verticalFlowLayout{new(layoutBase)}
}

type verticalFlowLayout struct {
	*layoutBase
}

func (l *verticalFlowLayout) relayout() {
	maxWidth, height := 0, 0
	for _, item := range l.items {
		b := item.bounds()
		if b.w > maxWidth {
			maxWidth = b.w
		}
		height += b.h
	}

	y := 0
	for _, item := range l.items {
		b := item.bounds()
		item.setBounds(rect{b.x + (maxWidth-b.w)/2, y, b.w, b.h})
	}
}
