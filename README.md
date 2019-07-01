#Go-config

###概要

这是一个实现配置功能的工具包,仍然在不断修改完善中，
导入方法：import "github.com/BlackCarDriver/config"


###支持配置类型
整数，浮点数，布尔值，字符串，段落，字符串数组，整数数组 , json   

###配置格式参考
    # match a string type, note that space in begin and end whill be trim
    t_string = it_is_a_string
    
    # match a string type, recommend format
    t_string2 = " it is a string "
    
    # match a integer type 
    t_int = 123456
    
    # match a float type
    t_float = 123.456
    
    # match a bool type
    t_bool = true
    
    # match a string array type,
    # note that comma should appended to every line
    t_str_arry = [
    "it is array one ...",
    "it is array two ...",
    "ti is array three ...",
    ]
    
    # match a integer array
    t_int_array = [
    12345,
    34567,
    67887,
    ]
    
    # match a multi line string
    t_muti_string = {
    it is the frist line of paragraph
    it is the second line of paragraph
    it is the last line of paragraph...
    }
    
    # match a struct in json format
    jsonText = {
    	"age" : 23,
    	"hobby" : ["game","foot","spoot"]
    }
    

### 使用案例

   	//create an config object by giving config path
	tc,err := NewConfig("./config/conf/")
	if err!=nil {
		fmt.Println("the following is the errors during reading config file :")
		fmt.Println(err)
	}

	//registe config by giving default value
	tc.Register("t_string", "test", true)
	tc.Register("t_string2", "test", true)
	tc.Register("t_muti_string", "test", true)
	tc.Register("t_int", 0, true)
	tc.Register("t_float", 0.1, true)
	tc.Register("t_bool", false, true)
	tc.Register("t_str_arry", make([]string,1), true)
	tc.Register("t_int_array", make([]int, 1), true)
	//register a new config key that don't exist in config file by setting isStrict = false
	tc.Register("newCOnfig", "t_muti_string", false)
	
	//display the config in map
	tc.Display()

	//get a config value by configName
	fmt.Println("test muti link string :")
	paragraph, err := tc.GetString("t_muti_string")
	if err!=nil {
		fmt.Println(err)
	}else{
		fmt.Println(paragraph)
	}

	//read struct from json syntax string
	fmt.Println("test read struct :")
	tc.Register("jsonText", "", true)
	child := childen{}
	err = tc.GetStruct("jsonText", &child)	
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println(child)
	}


    

## 

last change:   
2019/7/1 19:58:00 
