package database

type ErrorStatus uint

const (
	NoError ErrorStatus = iota
	NotFoundError
	InternalError
)

type Err struct {
	Msg    string
	Status ErrorStatus
	Err    error
}

func (e Err) Error() string {
	return e.Msg
}

func (e Err) NotFound() bool {
	return e.Status == NotFoundError
}

func (e Err) InternalError() bool {
	return e.Status == InternalError
}
