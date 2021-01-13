package practicessh

import (
	"testing"
)

func TestFlex(t *testing.T) {
	var(
		str []string
		arr [][]string
	)
	str = []string{"192.168.186.137:3306","192.168.186.137:3307"}

	if err = Init(); err != nil {
		t.Fatal("failed")
	}
	arr,err = Flex(str)
	if err != nil{
		t.Fatal("failed")
	}
	if arr == nil{
		t.Fatal("failed")
	}
}