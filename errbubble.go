package vali

type bubbleErr struct {
	err error
}

// BubbleErr can be used to wrap an error returned from
// a custom tag, to pass it as is from the validation func.
// It will preserve the type and the error itself allowing you to
// call `errors.Is` on the returne error from the `Validate()` method.
func BubbleErr(err error) *bubbleErr {
	return &bubbleErr{err}
}

func (b *bubbleErr) Error() string {
	return b.err.Error()
}
