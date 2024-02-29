package types

import (
	"cosmossdk.io/collections"
)

const (
	// ModuleName is the name of module for external use.
	ModuleName = "slpp"
	// StoreKey is the top-level store key for the oracle module.
	StoreKey = ModuleName
)

var (
	// AVSKeyPrefix is the key-prefix under which currency-pair state is stored.
	AVSKeyPrefix = collections.NewPrefix(0)

	// AVSIDKeyPrefix is the key-prefix under which the next currency-pairID is stored.
	AVSIDKeyPrefix = collections.NewPrefix(1)

	// UniqueIndexAVSKeyPrefix is the key-prefix under which the unique index on
	// currency-pairs is stored.
	UniqueIndexAVSKeyPrefix = collections.NewPrefix(2)

	// IDIndexAVSKeyPrefix is the key-prefix under which a currency-pair index.
	// is stored.
	IDIndexAVSKeyPrefix = collections.NewPrefix(3)
)
