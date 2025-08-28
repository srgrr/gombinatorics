package gombinator

import "errors"

var ErrInvalidPartitionKValue = errors.New("k must be at least 1")
var ErrInvalidRepeatKValue = errors.New("k must be -1 for infinite repetition, 0 for empty channel, or a positive integer")
