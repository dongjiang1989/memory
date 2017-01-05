package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestJson_StartsWithArray(t *testing.T) {
	assert := assert.New(t)

	const jsoncontextwitharray = `[
	{
		"url": "user",
		"serviceAPI": "http://www.test.com/user"
	},
	{
		"url": "employee",
		"serviceAPI": "http://www.test.com/employee"
	}
]`
	f, err := os.Create("testjsonWithArray.conf")

	assert.Nil(err, "create conf file err")

	_, err = f.WriteString(jsoncontextwitharray)

	assert.Nil(err, "write conf file err")

	defer f.Close()
	defer os.Remove("testjsonWithArray.conf")

	jsonconf, err := NewConfig("json", "testjsonWithArray.conf")

	assert.Nil(err, "read conf file err")

	rootArray, err := jsonconf.DIY("rootArray")

	assert.Nil(err, "get root body err: array does not exist as element")

	rootArrayCasted := rootArray.([]interface{})

	assert.NotNil(rootArrayCasted, "array from root is nil")

	elem := rootArrayCasted[0].(map[string]interface{})

	assert.Equal(elem["url"], "user", "array[0] values are not valid")
	assert.Equal(elem["serviceAPI"], "http://www.test.com/user", "array[0] values are not valid")

	elem2 := rootArrayCasted[1].(map[string]interface{})
	assert.Equal(elem2["url"], "employee", "array[0] values are not valid")
	assert.Equal(elem2["serviceAPI"], "http://www.test.com/employee", "array[1] values are not valid")
}

func TestJson_WithDict(t *testing.T) {
	assert := assert.New(t)

	var (
		jsoncontext = `{
"appname": "beeapi",
"testnames": "foo;bar",
"httpport": 8080,
"mysqlport": 3600,
"PI": 3.1415976, 
"runmode": "dev",
"autorender": false,
"copyrequestbody": true,
"session": "on",
"cookieon": "off",
"newreg": "OFF",
"needlogin": "ON",
"enableSession": "Y",
"enableCookie": "N",
"flag": 1,
"database": {
        "host": "host",
        "port": "port",
        "database": "database",
        "username": "username",
        "password": "password",
		"conns":{
			"maxconnection":12,
			"autoconnect":true,
			"connectioninfo":"info"
		}
    }
}`
		keyValue = map[string]interface{}{
			"appname":                         "beeapi",
			"testnames":                       []string{"foo", "bar"},
			"httpport":                        8080,
			"mysqlport":                       int64(3600),
			"PI":                              3.1415976,
			"runmode":                         "dev",
			"autorender":                      false,
			"copyrequestbody":                 true,
			"session":                         true,
			"cookieon":                        false,
			"newreg":                          false,
			"needlogin":                       true,
			"enableSession":                   true,
			"enableCookie":                    false,
			"flag":                            true,
			"database::host":                  "host",
			"database::port":                  "port",
			"database::database":              "database",
			"database::password":              "password",
			"database::conns::maxconnection":  12,
			"database::conns::autoconnect":    true,
			"database::conns::connectioninfo": "info",
			"unknown":                         "",
		}
	)

	f, err := os.Create("testjson.conf")
	assert.Nil(err, "create conf file err")

	_, err = f.WriteString(jsoncontext)

	assert.Nil(err, "write conf file err")

	defer f.Close()
	defer os.Remove("testjson.conf")
	jsonconf, err := NewConfig("json", "testjson.conf")

	assert.Nil(err, "read conf file err")

	for k, v := range keyValue {
		var err error
		var value interface{}
		switch v.(type) {
		case int:
			value, err = jsonconf.Int(k)
		case int64:
			value, err = jsonconf.Int64(k)
		case float64:
			value, err = jsonconf.Float(k)
		case bool:
			value, err = jsonconf.Bool(k)
		case []string:
			value = jsonconf.Strings(k)
		case string:
			value = jsonconf.String(k)
		default:
			value, err = jsonconf.DIY(k)
		}

		assert.Nil(err, fmt.Sprintf("get key %q value fatal,%v err %s", k, v, err))
		assert.Equal(fmt.Sprintf("%v", v), fmt.Sprintf("%v", value), fmt.Sprintf("get key %q value, want %v got %v .", k, v, value))

	}

	err = jsonconf.Set("name", "dongjiang")
	assert.Nil(err, "set data not succ")

	assert.Equal(jsonconf.String("name"), "dongjiang", "get data not equal")

	db, err := jsonconf.DIY("database")

	assert.Nil(err, "get sub interface not succ")

	m, ok := db.(map[string]interface{})
	t.Log(db)

	assert.True(ok, "db not map[string]interface{}")

	assert.Equal(m["host"].(string), "host", "get host err")

	_, err = jsonconf.Int("unknown")

	assert.NotNil(err, "unknown keys should return an error when expecting an Int")

	i, err := jsonconf.Int("httpport")

	assert.Nil(err, "get sub interface not succ")
	assert.Equal(i, 8080, fmt.Sprintf("httpport is not %d", i))

	_, err = jsonconf.Int64("unknown")
	assert.NotNil(err, "unknown keys should return an error when expecting an Int64")

	i6, err := jsonconf.Int64("mysqlport")
	assert.Nil(err, "get sub interface not succ")
	assert.Equal(i6, int64(3600), fmt.Sprintf("mysqlport is not %d", i6))

	_, err = jsonconf.Float("unknown")
	assert.NotNil(err, "unknown keys should return an error when expecting a Float")

	_, err = jsonconf.DIY("unknown")
	assert.NotNil(err, "unknown keys should return an error when expecting an interface{}")

	val := jsonconf.String("unknown")

	assert.Empty(val, "unknown keys should return an empty string when expecting a String")

	_, err = jsonconf.Bool("unknown")
	assert.NotNil(err, "unknown keys should return an error when expecting a Bool")

	assert.True(jsonconf.DefaultBool("unknow", true), "unknown keys with default value wrong")
}
