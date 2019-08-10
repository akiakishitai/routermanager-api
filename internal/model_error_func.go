package openapi

func (e *Error) httpStatus() int {
	switch e.Code {
	case 400:
		return 400
	case 404:
		return 404
	default:
		return 500
	}
}
