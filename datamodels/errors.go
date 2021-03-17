package datamodels

var (
	ERR_SUCCESS            = Error{ID: 00, Message: "OK"}
	ERR_NOT_LOGGED         = Error{ID: 10, Message: "User must be logged"}
	ERR_UNKNOWN_USER       = Error{ID: 12, Message: "Unknown User, try to log again"}
	ERR_BAD_TOKEN          = Error{ID: 14, Message: "Bad token"}
	ERR_UNREADABLE_TOKEN   = Error{ID: 16, Message: "Cannot decode token"}
	ERR_UNKNOW_RIDE        = Error{ID: 50, Message: "Unknown Ride"}
	ERR_RIDE_NOT_AVAILABLE = Error{ID: 52, Message: "Ride already dispatched"}
	ERR_EMPTY_REQUEST      = Error{ID: 90, Message: "Empty or malformed request"}
	ERR_NOT_IMPLEMENTED    = Error{ID: 99, Message: "Function Not implemented"}
)
