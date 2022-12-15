package catcher

type Catcher interface {
	StartCatching() error
	StopCatching() error
}
