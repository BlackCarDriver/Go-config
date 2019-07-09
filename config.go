package config

import(
	"fmt"
	"strconv"
	"os"
	"io"
	
	"errors"
	"strings"
	"bufio"
	"regexp"
	"encoding/json"
)

//recorde the filename that alread read, each file can only read once
var readHistory = make(map[string]bool)

/*

*/
type Config struct{
	configPath string
	isStrict bool
	rawConf map[string]string 
}

type ConfigMachine interface { 
	InitWithFilesPath(filesPath string) error
	SetIsStrict(bool)
	Display() 
	GetInt(keyName string) (value int, err error)
	GetInts(keyName string) (value []int, err error)
	GetFloat(keyName string) (value float64, err error)
	GetString(keyName string) (value string, err error)
	GetStrings(keyName string) (value []string, err error) 
	GetBool(keyName string) (value bool, err error)
	GetStruct(keyName string, container interface{}) error
}

//the mainly way of obtain a Config
func NewConfig(confPath string)(ConfigMachine, error) {
	newMachine := new(Config)
	err := newMachine.InitWithFilesPath(confPath)
	return newMachine, err
}


//=========== method in interface ===============
func (c *Config) InitWithFilesPath(Configpath string) error{
	if c.configPath != "" {
		return errors.New("You can't init the Confi twice!")
	}
	if !strings.HasSuffix(c.configPath, "/") {
		c.configPath += "/"
	}
	c.rawConf = make(map[string]string)
	c.configPath = Configpath;
	errList := c.readAllConfig()
	return errList
}

func (c *Config)SetIsStrict(strict bool){
	c.isStrict = strict
}

//display the key name and value name in rawMap and ripeMap
func (c *Config) Display(){
	fmt.Println( "======================== config lists ======================" )
	for k,v := range c.rawConf {
		fmt.Printf("----------- %v ----------- \n%-20v \n", k,v)
	}
	fmt.Println( "===========================================================" )
	fmt.Println()
}

func (c *Config)GetInt(keyName string) (value int, err error) {
	rawStr, ok := "", false
	if !isLegalName(keyName) {
		err = errors.New("keyName is not legal!")
		goto tail
	}
	rawStr, ok = c.rawConf[keyName]
	if !ok {
		err = errors.New("Can't not find config " + keyName)
		goto tail
	}
	value,err = strconv.Atoi(rawStr)
	if err != nil {
		goto tail
	}
	tail:
	if c.isStrict && err != nil {
		panic(err)
	}
	return value, err
}

func (c *Config)GetInts(keyName string) (value []int, err error) {
	tmpInt,rawStr,ok := -1, "", false
	tmpStrArry:= make([]string, 0)
	if !isLegalName(keyName) {
		err = errors.New("keyName is not legal!")
		goto tail
	}
	rawStr, ok = c.rawConf[keyName]
	if !ok {
		err = errors.New("Can't not find config " + keyName)
		goto tail
	}
	tmpStrArry = strings.Split(rawStr,",")
	for _,strInt := range tmpStrArry {
		tmpInt, err = strconv.Atoi(strInt)
		if err!=nil {
				goto tail
		}
		value = append(value, tmpInt)
	}
	tail:
	if c.isStrict && err != nil {
		panic(err)
	}
	return value, err
}

func (c *Config)GetString(keyName string) (value string, err error) {
	rawStr, ok := "", false
	if !isLegalName(keyName) {
		err = errors.New("keyName is not legal!")
		goto tail
	}
	rawStr, ok = c.rawConf[keyName]
	if !ok {
		err = errors.New("Can't not find config " + keyName)
		goto tail
	}
	value = rawStr
	tail:
	if c.isStrict && err != nil {
		panic(err)
	}
	return value, err
}

func (c *Config)GetStrings(keyName string) (value []string, err error) {
	rawStr, ok := "", false
	if !isLegalName(keyName) {
		err = errors.New("keyName is not legal!")
		goto tail
	}
	rawStr, ok = c.rawConf[keyName]
	if !ok {
		err = errors.New("Can't not find config " + keyName)
		goto tail
	}
	rawStr = strings.Trim(rawStr,`"`)
    value = strings.Split(rawStr, `","`)
	tail:
	if c.isStrict && err != nil {
		panic(err)
	}
	return value, err
}

func (c *Config)GetBool(keyName string) (value bool, err error) {
	rawStr, ok, tmpBool := "", false, false
	if !isLegalName(keyName) {
		err = errors.New("keyName is not legal!")
		goto tail
	}
	rawStr, ok = c.rawConf[keyName]
	if !ok {
		err = errors.New("Can't not find config " + keyName)
		goto tail
	}
	tmpBool, err = strconv.ParseBool(rawStr)
	if err != nil {
		goto tail 
	}
	value = tmpBool
	tail:
	if c.isStrict && err != nil {
		panic(err)
	}
	return value, err
}

func (c *Config)GetFloat(keyName string) (value float64, err error) {
	rawStr, ok := "", false
	if !isLegalName(keyName) {
		err = errors.New("keyName is not legal!")
		goto tail
	}
	rawStr, ok = c.rawConf[keyName]
	if !ok {
		err = errors.New("Can't not find config " + keyName)
		goto tail
	}
	value, err = strconv.ParseFloat(rawStr, 64)
	if err != nil {
		goto tail
	}
	tail:
	if c.isStrict && err != nil {
		panic(err)
	}
	return value, err
}

func (c *Config)GetStruct(keyName string, container interface{}) error{ 
	jsonText, err := c.GetString(keyName)
	if err != nil {
		goto tail
	}
	jsonText = "{" + jsonText + "}"
	err = json.Unmarshal([]byte(jsonText), &container)
	tail:
	if c.isStrict && err!=nil {
		panic(err)
	}
	return err
}

//=============== tools function ==========

//read all files with .conf suffix in configPath
func (c *Config)readAllConfig() error {
	filesPath := c.configPath
	file ,err := os.Open( filesPath )
	if err != nil {
		return err
	}
	defer file.Close()
	fi, err := file.Readdir(0)
	if err != nil {
		return err
	}
	errReport := ""
	for _, info := range fi {
		//only read files that name like *.conf
		if strings.HasSuffix(info.Name(), ".conf") == false {
			continue
		}
		//guarante each file only read one times
		if readHistory[info.Name()] {
			errReport += fmt.Sprintf("can not read %v, already read before...", info.Name())
			continue;
		}
		readHistory[info.Name()] = true
		tmpPath := filesPath + info.Name()
		err := c.readConfig(tmpPath)
		if err != nil {
			errReport += fmt.Sprintf("\n %v", err)
		}
	}
	if errReport == ""{
		return nil
	}
	return errors.New(errReport)
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
			return errors.New("Reading config was interupt because unexpect fomat of config (index <= 0): " +  string(lineByte) )
		}
		confName := strings.TrimSpace(line[:index])
		confValue := strings.TrimSpace(line[index+1:])
		if len(confName) == 0 || len(confValue) == 0 {	//unknow format
			return errors.New("Reading config was interupt because unexpect fomat of config (len==0): " +  string(lineByte) )
		}

		if isLegalName(confName) == false {				//config name not legal
			return errors.New("Config Name not legal at line : " + string(lineByte) )
		}

		if isStringType(confValue) {	//match string type
			confValue = strings.Trim(confValue, `"`)
			goto saveConf
		}

		if isNumberType(confValue) {	//match int or float type
			goto saveConf
		}

		if confValue=="true" || confValue == "false" {	//match bool type
			goto saveConf
		}
		//read an multi line string to rawMap, dont
		if confValue == `{`	 {		
			tmpStr := ""
			for {
				tmplineByte, _, tmpErr := buf.ReadLine()
				if tmpErr != nil { 	
					return fmt.Errorf("Readding worng by mistack after ‘%v’ , error: %v ", string(lineByte), tmpErr)
				}
				tmpline := string(tmplineByte)
				if strings.HasPrefix(strings.TrimSpace(tmpline), `}`) {
					break
				}
				tmpStr += tmpline 
				tmpStr += "\n"
			}
			confValue = tmpStr
			goto saveConf
		}	

		if confValue == "[" {		//mathch an array
			tmpStr := ""
			for {
				tmplineByte, _, tmpErr := buf.ReadLine()
				if tmpErr != nil { 	
					return fmt.Errorf("Unexpect error when reading array type config in or near : ‘%v’, error: %v ", string(lineByte), tmpErr)
				}
				tmpline := string(tmplineByte)
				tmpline = strings.TrimSpace(tmpline)
				if tmpline == "]" {
					break
				}
				if strings.HasSuffix(tmpline, ",") {
					tmpline = strings.TrimRight(tmpline, ",")
				}
				tmpStr += tmpline
				tmpStr += ","
			}
			confValue = strings.TrimRight(tmpStr, ",")
			goto saveConf
		}

	saveConf:
		c.rawConf[confName] = confValue
	}
	return nil
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

//judge if a name of config is legal
func isLegalName(confName string) bool {
	legalNameReg, _ := regexp.Compile(`^[a-zA-Z0-9_]+$`) 
	isLegal := legalNameReg.MatchString(confName)
	return isLegal
}

//judege if a config value match a string type, scuh as `"Is is config value"`
func isStringType(confValue string) bool {
	tmpStr := confValue
	counter := strings.Count(tmpStr, `"`)
	if counter != 2 {
		return false
	}
	tmpStr = strings.Trim(tmpStr, `"`)
	return (strings.Count(tmpStr, `"`) == 0)
}

//judege if a config value match a integer or float type
func isNumberType(confValue string) bool {
	_, isInt := strconv.Atoi(confValue)
	_, isFlo := strconv.ParseFloat(confValue, 64)
	return (isInt==nil || isFlo==nil)
}

