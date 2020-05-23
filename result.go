package vali

type Result struct {
	BubbleErr     error
	ValidationErr error
	Skip          bool
}
