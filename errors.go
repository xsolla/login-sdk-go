package login_sdk_go

type WrappedError struct {
	Inner error
}

func WrapError(err error) *WrappedError {
	return &WrappedError{
		err,
	}
}

func (we *WrappedError) Error() string {
	return we.Inner.Error()
}

func (we *WrappedError) Valid() bool {
	return we.Inner == nil
}

// IsExpired Deprecated: After change the main jwt library this code is for the backward compatibility.
func (we *WrappedError) IsExpired() bool {
	return false
}
