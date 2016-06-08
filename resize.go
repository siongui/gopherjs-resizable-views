package resizeviews

import (
	"github.com/gopherjs/gopherjs/js"
	"strconv"
	"strings"
)

type resizeManager struct {
	allContainer        *js.Object
	leftView            *js.Object
	viewwrapper         *js.Object
	arrow               *js.Object
	separator           *js.Object
	rightView           *js.Object
	leftViewWidth       int
	viewwrapperWidth    int
	rightViewWidth      int
	startLeftViewWidth  int
	startRightViewWidth int
	initialMouseX       int
}

func (rm *resizeManager) allWidth() int {
	return rm.allContainer.Get("offsetWidth").Int()
}

func (rm *resizeManager) initWidth() {
	rm.leftViewWidth = 250  //px
	rm.viewwrapperWidth = 7 //px
	rm.rightViewWidth = rm.allWidth() - rm.leftViewWidth - rm.viewwrapperWidth

	rm.leftView.Get("style").Set("width", strconv.Itoa(rm.leftViewWidth)+"px")
	rm.viewwrapper.Get("style").Set("width", strconv.Itoa(rm.viewwrapperWidth)+"px")
	rm.rightView.Get("style").Set("width", strconv.Itoa(rm.rightViewWidth)+"px")
}

func (rm *resizeManager) arrowOnClick(event *js.Object) {
	rm.leftViewWidth = 0
	rm.rightViewWidth = rm.allWidth() - rm.viewwrapperWidth
	rm.leftView.Get("style").Set("width", "0")
	rm.rightView.Get("style").Set("width", strconv.Itoa(rm.rightViewWidth)+"px")
}

func (rm *resizeManager) getLeftViewWidthFromCss() int {
	i, _ := strconv.Atoi(strings.Replace(rm.leftView.Get("style").Get("width").String(), "px", "", 1))
	return i
}

func (rm *resizeManager) getRightViewWidthFromCss() int {
	i, _ := strconv.Atoi(strings.Replace(rm.rightView.Get("style").Get("width").String(), "px", "", 1))
	return i
}

func (rm *resizeManager) separatorOnMouseDown(event *js.Object) {
	doc := js.Global.Get("document")
	event.Call("preventDefault")
	rm.startLeftViewWidth = rm.getLeftViewWidthFromCss()
	rm.startRightViewWidth = rm.getRightViewWidthFromCss()
	rm.initialMouseX = event.Get("clientX").Int()
	doc.Call("addEventListener", "mousemove", rm.documentOnMouseMove)
	doc.Call("addEventListener", "mouseup", rm.documentOnMouseUp)
}

func (rm *resizeManager) documentOnMouseMove(event *js.Object) {
	// calculate the delta of mouse cursor movement
	dx := event.Get("clientX").Int() - rm.initialMouseX

	newlw := rm.startLeftViewWidth + dx
	if newlw < 0 {
		rm.leftView.Get("style").Set("width", "0")
		rm.rightView.Get("style").Set("width", strconv.Itoa(rm.startLeftViewWidth+rm.startRightViewWidth)+"px")
		return
	}

	newrw := rm.startRightViewWidth - dx
	if newrw < 0 {
		rm.leftView.Get("style").Set("width", strconv.Itoa(rm.startLeftViewWidth+rm.startRightViewWidth)+"px")
		rm.rightView.Get("style").Set("width", "0")
		return
	}

	rm.leftView.Get("style").Set("width", strconv.Itoa(newlw)+"px")
	rm.rightView.Get("style").Set("width", strconv.Itoa(newrw)+"px")
}

func (rm *resizeManager) documentOnMouseUp(event *js.Object) {
	doc := js.Global.Get("document")
	doc.Call("removeEventListener", "mousemove", rm.documentOnMouseMove)
	doc.Call("removeEventListener", "mouseup", rm.documentOnMouseUp)
}

func NewResizeManager(allContainerId, leftViewId, viewwrapperId, arrowId, separatorId, rightViewId string) *resizeManager {
	doc := js.Global.Get("document")
	rm := &resizeManager{
		allContainer: doc.Call("getElementById", allContainerId),
		leftView:     doc.Call("getElementById", leftViewId),
		viewwrapper:  doc.Call("getElementById", viewwrapperId),
		arrow:        doc.Call("getElementById", arrowId),
		separator:    doc.Call("getElementById", separatorId),
		rightView:    doc.Call("getElementById", rightViewId),
	}

	rm.initWidth()
	rm.arrow.Set("onclick", rm.arrowOnClick)
	rm.separator.Set("onmousedown", rm.separatorOnMouseDown)

	return rm
}
