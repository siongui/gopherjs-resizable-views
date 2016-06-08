package main

import (
	rv "github.com/siongui/gopherjs-resizable-views"
)

func main() {
	rv.NewResizeManager("allContainer", "leftview", "viewwrapper", "viewarrow", "viewseparator", "rightview")
}
