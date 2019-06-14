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
)
/*
explain of Config struct:
configPath is the root path of config files, all files with .conf suffix will be read into rawConf
when the struct is init. 
rawConf is the map tmpely saving the string that read from config files, those string will not be
used until you register then.
ripeConf is the map saving config value, those config is read from rawConf through Register()
*/

var readHistory = make(map[string]bool)

type Config struct{
	configPath string
	rawConf map[string]string 
	ripeConf map[string]interface{}
}

type ConfigMachine interface {
	InitWithFilesPath(filesPath string) error
	Register(keyName string , dfValue interface{}, isImportant bool)(err error, warn string)
	Get(keyName string) (value interface{}, err error)
}

func handleErr(prefix string ,err error, isSeriou bool) ( errNotNull bool) {
	if err == nil {
		return false
	}
	fmt.Println(prefix , err)
	if isSeriou {
		os.Exit(2)
	}
	return true	 
}

//=========== method in interface ===============
func (c *Config) InitWithFilesPath(Configpath string) error{
	if c.configPath != "" {
		return errors.New("You can't init the Confi twice!")
	}
	c.rawConf = make(map[string]string)
	c.ripeConf = make(map[string]interface{}) 
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

//==============================================

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

//read all files with .conf suffix in configPath
func (c *Config)readAllConfig()( haveError bool ) {
	filesPath := c.configPath
	file ,err := os.Open( filesPath )
	handleErr("os.Open(filesPaht) ",err, true)
	defer file.Close()
	fi, err := file.Readdir(0)
	handleErr("file.Readdir(0) ",err, true)
	errCounter := 0
	for _, info := range fi {
		//guarante each file only read one times
		if readHistory[info.Name()] {
			fmt.Println(info.Name(), " already read before ...")
			continue;
		}
		readHistory[info.Name()] = true
		tmpPath := filesPath + info.Name()
		err := c.readConfig(tmpPath)
		if handleErr("c.readConfig(tmpPath) ",err,false) {
			errCounter ++
		}
	}
	return errCounter==0
}

//read a config file and save message into Conf.rawMap
func (c *Config)readConfig(path string) error {
	file,err := os.Open(path)
	if handleErr("os.Open(path) ", err, false) {
		return err
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	for{
		lineByte, _, err := buf.ReadLine()
		line := strings.TrimSpace( string(lineByte) )
		if err == io.EOF {	//end of file
			break
		}
		if err != nil {		//other error
			fmt.Println(err)
			return err
		}
		if line == "" {		//ignore empty line
			continue
		}
		if strings.HasPrefix(line, "#") {	//ignore cmment
			continue
		}
		index := strings.Index(line, "=")
		if index <= 0 {						//unknow format
			return errors.New("Unknow format in config files when reading : " + line)
		}
		confName := strings.TrimSpace(line[:index])
		confValue := strings.TrimSpace(line[index+1:])
		if len(confName) == 0 || len(confValue) == 0 {	//unknow format
			return errors.New("Unknow format in config files when reading : " + line)
		}
		//============need to do somtthing
	
		fmt.Println(confName, " ----> ", confValue)

		//=================================	
	}
	return nil
}

//judge if a name of config is legal
func isLegalName(name string) bool {
	//need to do something...
	return true
}

func Test(){
	var tc Config
	tc.InitWithFilesPath("./config/conf/")
}