package golang

func (e HttpError) Error() string {
	return e.Message
}
