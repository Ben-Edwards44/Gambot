package uci


import "gambot/src/engine/search"


type spinOption struct {
	name string
	defaultVal int
	min int
	max int
	setVal int
}


func createTTOpt() *spinOption {
	name := "Hash"
	defaultVal := search.DefaultTTSizeMib
	min := 1
	max := 128
	setVal := defaultVal

	opt := spinOption{name: name, defaultVal: defaultVal, min: min, max: max, setVal: setVal}

	return &opt
}


func (opt *spinOption) changeTTSize(newSize int) {
	opt.setVal = newSize
	search.NewTT(newSize)
}