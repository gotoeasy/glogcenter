package cmn

import (
	"log"
	"testing"
)

func Test_all(t *testing.T) {

	//	u64 := StringToUint64("18446744073709551615", 0)
	s64 := Uint64ToStringBase(10, 36)
	uu64 := StringToUint(s64, 36, 0)
	log.Println(s64, "-------------------", uu64, uu64)
}
