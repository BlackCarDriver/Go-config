package main

import(
	"./config"
	"fmt"
)

func main(){
	testInit()
	testUseage()
}

func testInit(){
	//create an Config object by config filename
	conf, _ := config.NewConfig("demo1.conf")
	//register some config valua
	conf.RegisterConf("myname","haha",true)
	conf.RegisterConf("isman",false, true)
	//create an new config value taht dont exist in config file
	conf.RegisterConf("newconf","testnewconf",false)
}

func testUseage(){
	//demo of create more than one Config object
	config.NewConfig("demo2.conf")
	//use GetConfig() to get an Config object
	demo1 := config.GetConfig("demo2.conf")
	demo1.RegisterConf("host","123.123.123.123",true)
	//display the value that in the map
	demo1.DisplayConf()
	demo2 := config.GetConfig("demo1.conf")
	name := demo2.Get("myname")
	isman := demo2.Get("isman")
	fmt.Println(name,"   ",isman)
}