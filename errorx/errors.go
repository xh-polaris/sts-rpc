package errorx

import "google.golang.org/grpc/status"

var (
	ErrInvalidUserId = status.Error(1001, "Invalid user ID")
)
