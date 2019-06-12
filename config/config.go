package config

import(
	"fmt"
	"reflect"
	"strconv"
	"os"
	"io"
	
	"errors"
	"strings"
	"bufio"
	"bytes"
)
/*
explain of Config struct:
	configPath is the root path of config files, all files with .conf suffix will be read into rawConf
when the struct is init. 
	rawConf is the map tmpely saving the string that read from config files, those string will not be
used until you register then.
	ripeConf is the map saving config value, those config is read from rawConf through Register()
*/
type Config struct{
	configPath string
	rawConf map[string]string;
	ripeConf map[string]interface{}
}

type ConfigMachine interface {
	readAllConfig()	error
	InitWithFilesPath(filesPath string) error
	Register(keyName string , dfValue interface{}, isImportant bool)(err error, warn string)
	Get(keyName string) (value interface{}, err error)
	DisplyConfList()
}

func (c *Config) InitWithFilesPath(Configpath string) error{
	if c.configPath != "" {
		return errors.New("You can't init the Confi twice!")
	}
	c.configPath = Configpath;
	c.readAllConfig()
	return nil
}

func (c *Config) Get(keyName string)(value interface{}, err error){
	if keyName == "" {
		err = errors.New("keyName is null")
		return
	}
	value, ok := c.ripeConf[keyName]
	if ok {
		err = nil
		return
	}else{
		err = fmt.Errorf("KeyName %v not found in config list!", keyName)
		return 
	}
}

func (c *Config) DisplyConfList(){
	fmt.Println("============= rawConf ========")
	for k,v := range c.rawConf {
		fmt.Println(k, " ---> ",v)
	}
	fmt.Println("============ ripefMap ========")
	for k,v := range c.ripeConf {
		fmt.Println(k, " ---> ",v)
	}
}

func (c *Config) Register(confName string, dfValue interface{}, isStrict bool)(err error, warn string){
	rawStr, ok := c.rawConf[confName]
	if !ok && isStrict {	
		err = fmt.Errorf("Can't not load config %v from file!", confName)
		return
	}
	if !ok && !isStrict {
		c.ripeConf[confName] = dfValue
		err = fmt.Errorf("Config %v not font in config files and we create it.", confName)
		return
	}
	tyName := reflect.TypeOf(dfValue).Name()
	switch tyName {
	case "int":
		tmpInt,err := strconv.Atoi(rawStr)
		if err != nil {
			panic(err)
		}
		c.ripeConf[confName] = tmpInt
	case "string":
		c.ripeConf[confName] = rawStr
	case "bool":
		tmpBool, err := strconv.ParseBool(rawStr)
		if err != nil {
			panic(err)
		}
		c.ripeConf[confName] = tmpBool
	}
	return
}

func (c *Config)readAllConfig() error {
	file,err := os.Open(c.configPath)
	if err!=nil {
		return err
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

