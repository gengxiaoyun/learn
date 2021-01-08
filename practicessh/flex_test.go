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
//func TestFlagCommand(t *testing.T) {
//	var(
//		str []string
//		arr [][]string
//	)
//	str = []string{"192.168.186.137:3306","192.168.186.137:3307"}
//
//	//r := gin.Default()
//	//r.GET("/",func(c *gin.Context){
//	//	c.JSON(200,c.QueryArray("address"))
//	//	fmt.Println(c.QueryArray("address"))
//	//	arr = FlagCommand(c.QueryArray("address"))
//	//
//	//})
//	//r.Run(":8080")
//
//	arr = FlagCommand(str)
//	if arr == nil{
//		t.Fatal("failed")
//	}
//}