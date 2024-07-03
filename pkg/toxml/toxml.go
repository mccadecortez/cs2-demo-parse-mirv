package toxml

import (
	"encoding/xml"
	"fmt"
)

// XXX: AI generation used

// Define struct types for the XML structure
type Command struct {
	Tick    string `xml:"tick,attr"`
	Command string `xml:"body"`
}

type Commands struct {
	CommandList []Command `xml:"c"`
}

type CommandSystem struct {
	XMLName  xml.Name `xml:"commandSystem"`
	Commands Commands `xml:"commands"`
}

func ToXML(commands []Command) string {
	commandSystem := CommandSystem{
		Commands: Commands{
			CommandList: commands,
		},
	}

	// Marshal the struct to XML
	xmlData, err := xml.MarshalIndent(commandSystem, "", "\t")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return ""
	}

	// Add the XML header
	xmlHeader := []byte(xml.Header)

	return string(xmlHeader) + string(xmlData)
}
