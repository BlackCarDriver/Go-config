#Go-config

**概要：**   

这是一个轻小的配置模块，主要功能是从配置中读取各种数据，保存并在程序中调用。


**特点：**

方便创建多个配置对象，每个对象绑定一个文件夹路径属性，初始化
时读取本路径下的全部配置文件（配置文件命名：*.con )。在代码中使用
Register() 注册配置的变量名和变量值。并使用 Get() 获取配置的变量值。



**支持配置类型**   
     
整数，浮点数，布尔值，字符串，段落，字符串数组，整数数组    

**包含文件：**   

主要代码位于： /config/config.go，   
配置格式位于：/congig/conf/example.conf

**example**   

    	//create an config object by giving config path
    	tc,err := New("./config/conf/")
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
    	paragraph, err := tc.Get("t_muti_string")
    	if err!=nil {
    		fmt.Println(err)
    	}else{
    		fmt.Println(paragraph)
    	}
    

## 

last change:   
2019/6/16 9:40:53 
