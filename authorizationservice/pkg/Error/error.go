package Error

type Error struct {
	InternalError    bool
	Err             error
	ErrorDescription map[string]interface{}
}

func (e *Error) Error() string {
	return e.Err.Error()
}
