// core implements shared functionality that both
// client and server can use.
//
// By exporting all the messages (Request and Response)
// it becomes very easy for the client to communicate
// back and forth with the server.

package core

import (
	"encoding/xml"
)

type Response struct {
	Ok         bool
	Message    string
	Vegitables Vegitables
}

type Request struct {
	Command []string
}

// A struct which contains the complete
// array of all vegitables in the file
type Vegitables struct {
	XMLName    xml.Name    `xml:"vegitables"`
	Vegitables []Vegitable `xml:"vegitable"`
}

// the vegitable struct, this contains
// vegitable name name, price per kg and
// remaining kgs
type Vegitable struct {
	XMLName      xml.Name `xml:"vegitable"`
	Name         string   `xml:"name"`
	PricePerKg   string   `xml:"pricePerKg"`
	RemainingKgs string   `xml:"remainingKgs"`
}
