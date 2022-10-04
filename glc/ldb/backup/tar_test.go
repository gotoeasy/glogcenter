package backup

import (
	"testing"
)

func Test_tar(t *testing.T) {
	//TarDir("F:\\222/default", "F:\\222\\000.tar")
	UnTar("F:\\222\\000.tar", "F:\\111\\xxxooo/aaaa")
}
