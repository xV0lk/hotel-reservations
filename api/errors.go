package api

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e *Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func ErrNotFound() Error {
	return NewError(404, "The id you provided is invalid")
}

func ErrBadRequest() Error {
	return NewError(400, "Bad Request")
}

func ErrInternal() Error {
	return NewError(500, "Internal Server Error")
}

func ErrUnauthorized() Error {
	return NewError(401, "Unauthorized")
}
