package practicessh

import (
	"testing"
)

func TestFlex(t *testing.T) {
	var(
		address string
		arr [][]string
	)
	address = "192.168.186.137:3306,192.168.186.137:3307"

	if err = Init(); err != nil {
		t.Fatal("failed")
	}
	arr,err = Flex(address)
	if err != nil{
		t.Fatal("failed")
	}
	if arr == nil{
		t.Fatal("failed")
	}
}
