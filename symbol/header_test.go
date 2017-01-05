package Symbol

import (
	//"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_symbolheader(t *testing.T) {
	assert := assert.New(t)

	str := `File:   ./SystemSymbol/8.1.1 (12B436)/Symbols/System/Library/PrivateFrameworks/Symbolication.framework/Symbolication
Format: Mach-O/32-Bit
Arch:   armv7s
Symbols:    1570
Tool Version:   1.0
File Version:   undefine
Built Time: 2016-08-05 11:51:32
UUID:   fffeb18603dd3abda51d0066bedb3c2c`

	obj := LoadSymbolHeader(str)

	assert.NotNil(obj)

	assert.Equal(obj.File, "./SystemSymbol/8.1.1 (12B436)/Symbols/System/Library/PrivateFrameworks/Symbolication.framework/Symbolication")
	assert.Equal(obj.Format, "Mach-O/32-Bit")
	assert.Equal(obj.Arch, "armv7s")
	assert.Equal(obj.Symbols, "1570")
	assert.Equal(obj.ToolVersion, "1.0")
	assert.Equal(obj.FileVersion, "undefine")
	assert.Equal(obj.BuiltTime, "2016-08-05 11:51:32")

	assert.Equal(obj.Uuid, "fffeb18603dd3abda51d0066bedb3c2c")
}

func Test_symbolheaderN(t *testing.T) {
	assert := assert.New(t)

	str := `File:   ./SystemSymbol/8.1.1 (12B436)/Symbols/System/Library/PrivateFrameworks/Symbolication.framework/Symbolication
Format: Mach-O/32-Bit
Arch:   armv7s
File Version:   undefine
Built Time: 2016-08-05 11:51:32
UUID:   fffeb18603dd3abda51d0066bedb3c2c`

	obj := LoadSymbolHeader(str)

	assert.NotNil(obj)

	assert.Equal(obj.File, "./SystemSymbol/8.1.1 (12B436)/Symbols/System/Library/PrivateFrameworks/Symbolication.framework/Symbolication")
	assert.Equal(obj.Format, "Mach-O/32-Bit")
	assert.Equal(obj.Arch, "armv7s")
	assert.Equal(obj.FileVersion, "undefine")
	assert.Equal(obj.BuiltTime, "2016-08-05 11:51:32")

	assert.Equal(obj.Uuid, "fffeb18603dd3abda51d0066bedb3c2c")
}

func Test_symbolheaderM(t *testing.T) {
	assert := assert.New(t)

	str := `Filel/8.1.1 (12B436)/Symbols/System/Library/PrivateFrameworks/Symbolication.framework/Symbolication
Format: Mach-O/32-Bit
Arch:   armv7s
File Version:   undefine
Built Time: 2016-08-05 11:51:32
UUID:   fffeb18603dd3abda51d0066bedb3c2c`

	obj := LoadSymbolHeader(str)

	assert.NotNil(obj)

	assert.Equal(obj.Format, "Mach-O/32-Bit")
	assert.Equal(obj.Arch, "armv7s")
	assert.Equal(obj.FileVersion, "undefine")
	assert.Equal(obj.BuiltTime, "2016-08-05 11:51:32")

	assert.Equal(obj.Uuid, "fffeb18603dd3abda51d0066bedb3c2c")
}

func Test_symbolheader1(t *testing.T) {
	assert := assert.New(t)

	str := `File:   ./SystemSymbol/8.1.1 (12B436)/Symbols/System/Library/PrivateFrameworks/Symbolication.framework/Symbolication
Format: Mach-O/32-Bit
Arch:   armv7s
Symbols:    1570
Tool Version:   1.0
File Version:   undefine
Built Time: 2016-08-05 11:51:32
UUID:   fffeb18603dd3abda51d0066bedb3c2c`

	obj := LoadSymbolHeaderByte([]byte(str))

	assert.NotNil(obj)

	assert.Equal(obj.File, "./SystemSymbol/8.1.1 (12B436)/Symbols/System/Library/PrivateFrameworks/Symbolication.framework/Symbolication")
	assert.Equal(obj.Format, "Mach-O/32-Bit")
	assert.Equal(obj.Arch, "armv7s")
	assert.Equal(obj.Symbols, "1570")
	assert.Equal(obj.ToolVersion, "1.0")
	assert.Equal(obj.FileVersion, "undefine")
	assert.Equal(obj.BuiltTime, "2016-08-05 11:51:32")

	assert.Equal(obj.Uuid, "fffeb18603dd3abda51d0066bedb3c2c")
}

func Test_symbolheaderN1(t *testing.T) {
	assert := assert.New(t)

	str := `File:   ./SystemSymbol/8.1.1 (12B436)/Symbols/System/Library/PrivateFrameworks/Symbolication.framework/Symbolication
Format: Mach-O/32-Bit
Arch:   armv7s
File Version:   undefine
Built Time: 2016-08-05 11:51:32
UUID:   fffeb18603dd3abda51d0066bedb3c2c`

	obj := LoadSymbolHeaderByte([]byte(str))

	assert.NotNil(obj)

	assert.Equal(obj.File, "./SystemSymbol/8.1.1 (12B436)/Symbols/System/Library/PrivateFrameworks/Symbolication.framework/Symbolication")
	assert.Equal(obj.Format, "Mach-O/32-Bit")
	assert.Equal(obj.Arch, "armv7s")
	assert.Equal(obj.FileVersion, "undefine")
	assert.Equal(obj.BuiltTime, "2016-08-05 11:51:32")

	assert.Equal(obj.Uuid, "fffeb18603dd3abda51d0066bedb3c2c")
}

func Test_symbolheaderM1(t *testing.T) {
	assert := assert.New(t)

	str := `Filel/8.1.1 (12B436)/Symbols/System/Library/PrivateFrameworks/Symbolication.framework/Symbolication
Format: Mach-O/32-Bit
Arch:   armv7s
File Version:   undefine
Built Time: 2016-08-05 11:51:32
UUID:   fffeb18603dd3abda51d0066bedb3c2c`

	obj := LoadSymbolHeaderByte([]byte(str))

	assert.NotNil(obj)

	assert.Equal(obj.Format, "Mach-O/32-Bit")
	assert.Equal(obj.Arch, "armv7s")
	assert.Equal(obj.FileVersion, "undefine")
	assert.Equal(obj.BuiltTime, "2016-08-05 11:51:32")

	assert.Equal(obj.Uuid, "fffeb18603dd3abda51d0066bedb3c2c")
}
