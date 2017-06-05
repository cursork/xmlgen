package main

import (
	"encoding/xml"
	x "github.com/cursork/xmlgen"
	"os"
)

type Xxx struct {
	XMLName xml.Name `xml:"http://something/ abc"`
	X       string   `xml:"xyz"`
}

func main() {
	e :=
		x.E("Foo", x.NoAttrs(),
			x.E("Bar", map[string]interface{}{"someattr": "&& a <value>"},
				123.0234843920,
				" and a string",
				x.E("AndAnElement", x.NoAttrs(),
					x.E("Foo", x.NoAttrs(), "BO&OM")),
			),
			// Use a custom ToElement() func
			&Xxx{X: "test"},
		)
	if err := e.Marshal(os.Stdout); err != nil {
		panic(err)
	}
}

func (*Xxx) ToElement() *x.Element {
	return x.E("thisisatest", x.NoAttrs())
}

/*
$ go run main.go | xmllint --format -
<?xml version="1.0"?>
<Foo>
  <Bar someattr="&amp;&amp; a &lt;value&gt;">123.023484 and a string<AndAnElement><Foo>BO&amp;OM</Foo></AndAnElement></Bar>
  <thisisatest/>
</Foo>
*/
