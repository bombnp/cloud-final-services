package ethutils

import (
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// ParseABI parses abi string
func ParseABI(ctrABI string) abi.ABI {
	a, err := abi.JSON(strings.NewReader(ctrABI))
	if err != nil {
		log.Fatal("Can't parse contract abi", err)
	}
	return a
}
