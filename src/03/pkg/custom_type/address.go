package custom_type

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type Address common.Address

func (a *Address) Scan(value any) error {
	switch v := value.(type) {
	case string:
		*a = HexToAddress(v)
		return nil
	case []byte:
		*a = HexToAddress(string(v))
		return nil
	default:
		return fmt.Errorf("unsupported type for address scan: %T", value)
	}
}

func (a Address) Value() (driver.Value, error) {
	return common.Address(a).Hex(), nil
}

// MarshalJSON serializes the custom_type.Address into a JSON string.
func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(common.Address(a).Hex())
}

// UnmarshalJSON deserializes a JSON string into the custom_type.Address.
func (a *Address) UnmarshalJSON(data []byte) error {
	var hex string
	if err := json.Unmarshal(data, &hex); err != nil {
		return fmt.Errorf("failed to unmarshal custom_type.Address: %v", err)
	}
	if !common.IsHexAddress(hex) {
		return fmt.Errorf("invalid hex address: %s", hex)
	}
	*a = HexToAddress(hex)
	return nil
}

func HexToAddress(hex string) Address {
	return Address(common.HexToAddress(hex))
}
