package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestIni_Parse(t *testing.T) {
	assert := assert.New(t)

	var (
		inicontext = `
;comment one
#comment two
appname = beeapi
httpport = 8080
mysqlport = 3600
PI = 3.1415976
runmode = "dev"
autorender = false
copyrequestbody = true
session= on
cookieon= off
newreg = OFF
needlogin = ON
enableSession = Y
enableCookie = N
flag = 1
[demo]
key1="asta"
key2 = "xie"
CaseInsensitive = true
peers = one;two;three
`

		keyValue = map[string]interface{}{
			"appname":               "beeapi",
			"httpport":              8080,
			"mysqlport":             int64(3600),
			"pi":                    3.1415976,
			"runmode":               "dev",
			"autorender":            false,
			"copyrequestbody":       true,
			"session":               true,
			"cookieon":              false,
			"newreg":                false,
			"needlogin":             true,
			"enableSession":         true,
			"enableCookie":          false,
			"flag":                  true,
			"demo::key1":            "asta",
			"demo::key2":            "xie",
			"demo::CaseInsensitive": true,
			"demo::peers":           []string{"one", "two", "three"},
			"null":                  "",
			"demo2::key1":           "",
			"error":                 "",
			"emptystrings":          []string{},
		}
	)

	f, err := os.Create("testini.conf")

	assert.Nil(err, "Create file not succ!")

	_, err = f.WriteString(inicontext)
	defer f.Close()

	assert.Nil(err, "open file and write not succ!")

	defer os.Remove("testini.conf")
	iniconf, err := NewConfig("ini", "testini.conf")

	assert.Nil(err, "open file  not succ!")

	for k, v := range keyValue {
		var err error
		var value interface{}
		switch v.(type) {
		case int:
			value, err = iniconf.Int(k)
		case int64:
			value, err = iniconf.Int64(k)
		case float64:
			value, err = iniconf.Float(k)
		case bool:
			value, err = iniconf.Bool(k)
		case []string:
			value = iniconf.Strings(k)
		case string:
			value = iniconf.String(k)
		default:
			value, err = iniconf.DIY(k)
		}

		assert.Nil(err, fmt.Sprintf("get key %q value fail,err %s", k, err))
		assert.Equal(fmt.Sprintf("%v", v), fmt.Sprintf("%v", value), fmt.Sprintf("%v is not equal %v", v, value))
	}

	err = iniconf.Set("name", "dongjiang")
	assert.Nil(err, "set in iniconf error")
	t.Log(iniconf)
	assert.Equal(iniconf.String("name"), "dongjiang", "get name error")

}

func TestIni_Save(t *testing.T) {
	assert := assert.New(t)

	const (
		inicontext = `
app = app
;comment one
#comment two
# comment three
appname = beeapi
httpport = 8080
# DB Info
# enable db
[dbinfo]
# db type name
# suport mysql,sqlserver
name = mysql
`

		saveResult = `
app=app
#comment one
#comment two
# comment three
appname=beeapi
httpport=8080
# DB Info
# enable db
[dbinfo]
# db type name
# suport mysql,sqlserver
name=mysql
`
	)
	cfg, err := NewConfigData("ini", []byte(inicontext))

	assert.Nil(err, "read conf ini for []byte err")

	name := "newIniConfig.ini"

	err = cfg.SaveConfigFile(name)

	assert.Nil(err, "Save conf err")

	defer os.Remove(name)

	data, err := ioutil.ReadFile(name)

	assert.Nil(err, "Save conf to read err")

	cfgData := string(data)

	datas := strings.Split(saveResult, "\n")

	for _, line := range datas {
		assert.True(strings.Contains(cfgData, line+"\n"), fmt.Sprintf("different after save ini config file. need contains %q", line))
	}
}
