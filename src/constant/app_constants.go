package constant

type ResponseStatus int
type Headers int
type General int

// Auth methods
const (
	Telegram string = "telegram"
)

// Constant Api
const (
	Success ResponseStatus = iota + 1
	DataNotFound
	UnknownError
	InvalidRequest
	Unauthorized
	WrongBody
	WrongMethod
	WrongDataBody
)

func (r ResponseStatus) GetResponseStatus() string {
	return [...]string{"SUCCESS", "DATA_NOT_FOUND", "UNKNOWN_ERROR", "INVALID_REQUEST", "UNAUTHORIZED", "WRONG_BODY", "WRONG_METHOD", "WRONG_DATA_BODY"}[r-1]
}

type UpgradeLvl int

func GetLvlMaxFields(r int) int {
	return [...]int{4, 8, 16}[r-1]
}
