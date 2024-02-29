package types

import (
	"fmt"
)

type AVSContractAlreadyExistsError struct {
	contractBytesHex string
}

func NewAVSContractAlreadyExistsError(contractBytesHex string) AVSContractAlreadyExistsError {
	return AVSContractAlreadyExistsError{contractBytesHex}
}

func (e AVSContractAlreadyExistsError) Error() string {
	return fmt.Sprintf("avs for contract bytes already exists: %d", e.contractBytesHex)
}
