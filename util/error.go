package util

func ErrorMessage(err error) string {
	if err == nil {
		return "non"
	}
	return err.Error()
}
