package config

import(
	"fmt"
	"reflect"
	"strconv"
	"os"
	"io"
	
	"io/ioutil"
	"strings"
	"bufio"
	"bytes"
	"log"
)

var rootPath = "./config/conf/"

//confFileName save the file name that already used by Config
var confFileName = make(map[string]bool)

var ConfigSaver = make(map[string]Config)

//a Config behalf of an config file, all config files should save in an same floder
//fileName is the name of config files such as database.conf
type Config struct{
	fileName string
	rawConf  map[string]string 
	confMap  map[string]interface{}
} 

//set up the path of floder that saving the config files
func SetRoot(root string){
	if rootPath != "./config/conf/" {
		log.Fatal("Can't call SetRoot() for twice! ")
	}
	_,err := ioutil.ReadDir(root)
	if os.IsNotExist(err) {
		wd,_ := os.Getwd()
		fmt.Println("config rootPath not exist ! please check : wd is :",wd , " + ", root)
		os.Exit(2)
	}
	rootPath = root
}

//Create an Config, fileName is the name of config file
func NewConfig(fileName string)(c Config, err error){
	if _,have := confFileName[fileName]; have {
		log.Fatal("the config file already have been used ! filename: ",fileName)
	}
	filepath := rootPath + fileName
	_, err = os.Open(filepath)
	if err !=nil {
		log.Fatal("config file not exist : ",err)
	}
	c.fileName = fileName
	c.rawConf = make(map[string]string)
	c.confMap = make(map[string]interface{})
	c.readConf()
	confFileName[fileName] = true
	ConfigSaver[fileName] = c
	return
}

//Get a Config that already create before
func GetConfig(ConfigName string) Config {
	c,ok := ConfigSaver[ConfigName]
	if ok {
		return c
	}
	log.Fatal("Try to get a Config that do not exist! name :", ConfigName)
	return c
}

//register a config on confmap, dfValue is default value, mainly used for defind the 
//type of config value. if isstrict is true, the config value will be only read in config file
//if isStrict is false, the config will use dfValue when config file don't exist confName
func (c *Config) RegisterConf(confName string, dfValue interface{}, isStrict bool){
	rawStr, ok := c.rawConf[confName]
	if isStrict && !ok {
		err := fmt.Errorf("Can't not load config %v from file!", confName)
		panic(err)
	}
	if !ok && !isStrict {
		c.confMap[confName] = dfValue
		return
	}
	tyName := reflect.TypeOf(dfValue).Name()

	switch tyName {
	case "int":
		tmpInt,err := strconv.Atoi(rawStr)
		if err != nil {
			panic(err)
		}
		c.confMap[confName] = tmpInt
	case "string":
		c.confMap[confName] = rawStr
	case "bool":
		tmpBool, err := strconv.ParseBool(rawStr)
		if err != nil {
			panic(err)
		}
		c.confMap[confName] = tmpBool
	}

	return
}

//get a config value by name, if this name is not exist in confMap, it will call panic
func (c *Config) Get(confName string) interface{} {
	res, ok := c.confMap[confName]
	if !ok {
		err := fmt.Errorf("Config %v don't exit !", confName)
		log.Fatal(err)
	}
	return res
}

//printf config map for test
func (c *Config) DisplayConf(){
	fmt.Println("============= rawConf ========")
	for k,v := range c.rawConf {
		fmt.Println(k, " ---> ",v)
	}
	fmt.Println("============ confMap ========")
	for k,v := range c.confMap {
		fmt.Println(k, " ---> ",v)
	}
}


//================== private method =============================================


//read config file and record thos values in rawConf, the type of those config as yet
//are all strings, the config can be used only after RegisterConf() 
func (c *Config) readConf() error {
	file,err := os.Open(rootPath + c.fileName)
	if err!=nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	for{
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil
		}
		if bytes.Equal(line,  []byte{}) {
			continue
		}
		if bytes.HasPrefix(line, []byte{'#'} ) {
			continue
		}
		tmpStr := strings.TrimSpace(string(line))
		index := strings.Index(tmpStr, "=")
		if index <= 0 {
			continue
		}
		key := strings.TrimSpace(tmpStr[:index])
		value := strings.TrimSpace(tmpStr[index+1:])
		if len(key) == 0 || len(value) == 0 {
			continue
		}
		c.rawConf[key] = value
	}
	return nil
}

