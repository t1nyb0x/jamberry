package usecase

// ValidationError はバリデーションエラーを表します
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NotFoundError は結果が見つからないエラーを表します
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

// IsValidationError はエラーがValidationErrorかどうかを判定します
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// IsNotFoundError はエラーがNotFoundErrorかどうかを判定します
func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}
