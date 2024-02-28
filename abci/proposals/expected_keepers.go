package proposals

import (
	slpptypes "githbu "
)

// SLPPKeeper represents the expected interface for the slpp keeper
type SLPPKeeper interface {
	GetAVSPerID(id uint64)
}