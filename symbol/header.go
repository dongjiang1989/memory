package Symbol

import (
	"bytes"
	"strings"

	log "github.com/cihub/seelog"
)

type Header struct {
	File        string
	Format      string
	Arch        string
	Symbols     string
	ToolVersion string
	FileVersion string
	BuiltTime   string
	Uuid        string
}

//File:   ./SystemSymbol/8.1.1 (12B436)/Symbols/System/Library/PrivateFrameworks/Symbolication.framework/Symbolication
//Format: Mach-O/32-Bit
//Arch:   armv7s
//Symbols:    1570
//Tool Version:   1.0
//File Version:   undefine
//Built Time: 2016-08-05 11:51:32
//UUID:   fffeb18603dd3abda51d0066bedb3c2c

func (header *Header) Add(info string) bool {
	heads := strings.SplitN(info, ":", 2)
	if len(heads) != 2 {
		log.Error("one line to Split is not Symbol header Temp! str:", info, "array:", heads)
	} else {
		key, value := strings.Trim(heads[0], " \t"), strings.Trim(heads[1], " \t")

		switch key {
		case "File":
			header.File = value
		case "Format":
			header.Format = value
		case "Arch":
			header.Arch = value
		case "Symbols":
			header.Symbols = value
		case "Tool Version":
			header.ToolVersion = value
		case "File Version":
			header.FileVersion = value
		case "Built Time":
			header.BuiltTime = value
		case "UUID":
			header.Uuid = value
		default:
			return false
		}
	}
	return true
}

func LoadSymbolHeader(block string) *Header {
	if block == "" {
		return nil
	} else {
		infos := strings.Split(block, "\n")
		var header Header
		for _, info := range infos {
			header.Add(info)
		}

		return &header
	}

	return nil
}

func LoadSymbolHeaderByte(block []byte) *Header {
	if bytes.Equal(block, []byte("")) {
		return nil
	} else {
		infos := bytes.Split(block, []byte{'\n'})
		var header Header
		for _, info := range infos {
			heads := bytes.SplitN(info, []byte(":"), 2)
			if len(heads) != 2 {
				log.Error("one line to Split is not Symbol header Temp! str:", info, "array:", heads)
			} else {
				key, value := strings.Trim(string(heads[0]), " \t"), strings.Trim(string(heads[1]), " \t")

				switch key {
				case "File":
					header.File = value
				case "Format":
					header.Format = value
				case "Arch":
					header.Arch = value
				case "Symbols":
					header.Symbols = value
				case "Tool Version":
					header.ToolVersion = value
				case "File Version":
					header.FileVersion = value
				case "Built Time":
					header.BuiltTime = value
				case "UUID":
					header.Uuid = value
				default:
				}
			}

		}

		return &header
	}

	return nil
}
