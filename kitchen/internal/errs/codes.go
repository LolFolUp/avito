package errs

type Kind string

const (
	KindInvalidArgument    Kind = "invalid_argument"
	KindNotFound           Kind = "not_found"
	KindFailedPrecondition Kind = "failed_precondition"
	KindAlreadyExists      Kind = "already_exists"
	KindInternal           Kind = "internal"
)

const (
	CodeInvalidArgument         = "INVALID_ARGUMENT"
	CodeNotFound                = "NOT_FOUND"
	CodeAlreadyExists           = "ALREADY_EXISTS"
	CodeMenuChanged             = "MENU_CHANGED"
	CodeRestaurantNotAccepting  = "RESTAURANT_NOT_ACCEPTING_ORDERS"
	CodeInvalidStatusTransition = "INVALID_STATUS_TRANSITION"
	CodeStatusConflict          = "STATUS_CONFLICT"
	CodeInternal                = "INTERNAL"
)
