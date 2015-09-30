package main

type layout interface {
	addElement(bounded) // TODO remove this and add parameter to relayout
	relayout(in rect)
}

type bounded interface {
	bounds() rect
	setBounds(rect)
}

// dummy layout does not change its elements

func newDummyLayout() dummyLayout { return dummyLayout{} }

type dummyLayout struct{}

func (dummyLayout) addElement(bounded) {}
func (dummyLayout) relayout(rect)      {}

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

func newCenterLayout() *centerLayout {
	return &centerLayout{layoutBase{}}
}

type centerLayout struct {
	layoutBase
}

func (l *centerLayout) relayout(in rect) {
	cx, cy := in.x+in.w/2, in.y+in.h/2
	bounds := boundingBox(boundingBoxableItems(l.items))
	boundsCx, boundsCy := bounds.x+bounds.w/2, bounds.y+bounds.h/2
	dx, dy := cx-boundsCx, cy-boundsCy
	for _, item := range l.items {
		b := item.bounds()
		item.setBounds(rect{b.x + dx, b.y + dy, b.w, b.h})
	}
}

type boundingBoxableItems []bounded

func (slice boundingBoxableItems) Len() int         { return len(slice) }
func (slice boundingBoxableItems) At(i int) bounder { return slice[i] }

// The vertical flow layout puts all elements directly underneath each other and
// centers them horizontally on each row. The whole block will be vertically
// centered as well.

func newVerticalFlowLayout(verticalSpaceBetweenElements int) *verticalFlowLayout {
	return &verticalFlowLayout{layoutBase{}, verticalSpaceBetweenElements}
}

type verticalFlowLayout struct {
	layoutBase
	yMargin int
}

func (l *verticalFlowLayout) relayout(in rect) {
	height := 0
	if len(l.items) > 0 {
		height = l.items[0].bounds().h
		for i := 1; i < len(l.items); i++ {
			height += l.yMargin + l.items[i].bounds().h
		}
	}

	y := in.y + (in.h-height)/2
	for _, item := range l.items {
		b := item.bounds()
		item.setBounds(rect{b.x + (in.w-b.w)/2, y, b.w, b.h})
		y += b.h + l.yMargin
	}
}

// The composite layout simply applies all layouts one by one to its items.

func newCompositeLayout(first layout, others ...layout) *compositeLayout {
	return &compositeLayout{append([]layout{first}, others...)}
}

type compositeLayout struct {
	layouts []layout
}

func (l *compositeLayout) addElement(b bounded) {
	for _, layout := range l.layouts {
		layout.addElement(b)
	}
}

func (l *compositeLayout) relayout(in rect) {
	for _, layout := range l.layouts {
		layout.relayout(in)
	}
}

// The top-left layout sets the first element at the top-left of the destination
// rectangle and then simply puts all following elements under it with no margin
// in y and at x=in.x .

func newTopLeftLayout() *topLeftLayout {
	return &topLeftLayout{layoutBase{}}
}

type topLeftLayout struct {
	layoutBase
}

func (l *topLeftLayout) relayout(in rect) {
	y := in.y
	for _, item := range l.items {
		b := item.bounds()
		item.setBounds(rect{in.x, y, b.w, b.h})
		y += b.h
	}
}
