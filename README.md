#Go-config

### 概要

这是一个实现配置功能的工具包，用于为需要读取配置功能的小型项目提供读取一些常见类型数据的功能，
主要优点是简单易用，仍然在不断修改完善中。


### 支持配置类型
整数，整数数组, 浮点数，布尔值，字符串，字符串数组，段落，json  

### 配置格式参考
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
    
### 字段说明
    一个配置对象，与一个文件夹绑定，将会读取这个文件夹中的全部*.conf文件并保存其中格式正确的配置  
     type Config struct {...}

    配置接口, 包括了本包提供使用的函数
    type ConfigMachine interface{...}
    
    得到一个配置对象，通过这个对象来获取配置文件中的配置值
    func NewConfig(confPath string)(ConfigMachine, error) 

    设置获取一个错误配置时的操作，若strict为true,则读取错误配置后引发panic,否则返回遇到的错误
    func (c *Config)SetIsStrict(strict bool)
    
    得到相应的配置值，值的类型又函数名可知
    (c *Config)GetInt(keyName string) (value int, err error)
    (c *Config)GetInts(keyName string) (value []int, err error)
    (c *Config)GetString(keyName string) (value string, err error)
    (c *Config)GetStrings(keyName string) (value []string, err error)
    (c *Config)GetBool(keyName string) (value bool, err error)
    (c *Config)GetFloat(keyName string) (value float64, err error)
    
    将配置文件中的json格式的值转换为对应的结构体
    (c *Config)GetStruct(keyName string, container interface{}) error


### 使用例子
    package main
    
    import(
    	"github.com/BlackCarDriver/config"
    	"fmt"
    )
    
    type childen struct {
    	Age int `json:"age"`
    	Hobby []string `json:"hobby"`
    }
    
    func main(){
    	tc,err := config.NewConfig("./")
    	if err != nil {
    		panic(err)
    	}
    	//tc.Display()
    	tc.SetIsStrict(true)
    	
    	t_string, _ := tc.GetString("t_string")
    	fmt.Println("tstring : ", t_string)
    	
    	t_int, _ := tc.GetInt("t_int")
    	fmt.Println("t_int : ", t_int)
    	
    	t_float, _ := tc.GetFloat("t_float")
    	fmt.Println("t_float : ", t_float)
    
    	t_bool, _ := tc.GetBool("t_bool")
    	fmt.Println("t_bool : ", t_bool)
    
    	t_str_arry, _ := tc.GetStrings("t_str_arry")
    	fmt.Println("t_str_arry : ", t_str_arry)
    
    	t_int_array, _ := tc.GetInts("t_int_array")
    	fmt.Println("t_int_array : ", t_int_array)
    
    	tstruct := childen{}
    	tc.GetStruct("jsonText", &tstruct)
    	fmt.Println(tstruct)
    }
    

## 

last change:   
2019/7/9 10:55:33 
