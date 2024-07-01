# rdf2doormud
A RDF Ontology Adventure in Golang

![alt text](https://github.com/warmmike/rdf2doormud/blob/main/1-1.png?raw=true)

## Examples
DOORMUD Game
```
package main

import (
	"fmt"

	"github.com/warmmike/rdf2doormud"
)

func main() {
	p, s, _ := rdf2doormud.GetStart(rdf2doormud.StartUri)

	for {
		fmt.Println(rdf2doormud.StringToAsciiGreen(p.City + ", " + p.Property))
		if p.Comment != "" {
			fmt.Println(rdf2doormud.StringToAsciiWhite(p.Comment))
		}
		exits := rdf2doormud.DirectionsToString(p.Directions)
		fmt.Println(rdf2doormud.StringToAsciiBlue("Obvious exits: " + exits))
		fmt.Print(rdf2doormud.StringToAsciiWhite("[ HP:64 SP:9 MF:100 ]: "))

		ui := rdf2doormud.InputToLong(rdf2doormud.GetInput())

		if ui == "quit" {
			break
		} else if ui == "help" {
			break
		} else {
			err := rdf2doormud.ValidateDirection(p.Directions, ui)
			if err != nil {
				fmt.Println(rdf2doormud.StringToAsciiWhite("You cannot go that way."))
			} else {
				p, s, err = rdf2doormud.GetPlace(s, p.Directions[ui])
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
```
