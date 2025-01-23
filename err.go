package yafasm

import (
	"fmt"
)

var ErrConditionFailed = fmt.Errorf("condition failed")

var ErrStoreSaveFailed = fmt.Errorf("store save failed")

var ErrNoTransitionFound = fmt.Errorf("no transition found")
