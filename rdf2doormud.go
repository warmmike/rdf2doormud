package rdf2doormud

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/deiu/rdf2go"
)

const (
	startUri          = "file://./1-1.ttl"
	startLabel        = "start"
	labelPredicateUri = "http://www.w3.org/2004/02/skos/core#hiddenLabel"
	placeObjectUri    = "http://schema.org#Place"
)

var (
	ErrNotFound = errors.New("place was not found in graph")
)

var (
	ldirections = []string{"north", "south", "east", "west", "northwest", "southwest", "northeast", "southeast"}
)

type datavalue struct {
	Value string `json:"@value"`
}
type dataid struct {
	Value string `json:"@id"`
}

type jsonMap struct {
	Url       []datavalue `json:"http://schema.org#url"`
	Comment   []datavalue `json:"http://www.w3.org/2000/01/rdf-schema#comment"`
	Homepage  []dataid    `json:"http://xmlns.com/foaf/0.1/homepage"`
	Property  []datavalue `json:"http://schema.org#property"`
	City      []datavalue `json:"http://schema.org#city"`
	Level     []datavalue `json:"http://schema.org#floorLevel"`
	Latitude  []datavalue `json:"http://opengis.net/ont/geosparql#lat"`
	Longitude []datavalue `json:"http://opengis.net/ont/geosparql#long"`
	North     []dataid    `json:"http://geosysiot.in/rsso/ApplicationSchema#north"`
	South     []dataid    `json:"http://geosysiot.in/rsso/ApplicationSchema#south"`
	East      []dataid    `json:"http://geosysiot.in/rsso/ApplicationSchema#east"`
	West      []dataid    `json:"http://geosysiot.in/rsso/ApplicationSchema#west"`
	Northeast []dataid    `json:"http://geosysiot.in/rsso/ApplicationSchema#northeast"`
	Northwest []dataid    `json:"http://geosysiot.in/rsso/ApplicationSchema#northwest"`
	Southeast []dataid    `json:"http://geosysiot.in/rsso/ApplicationSchema#southeast"`
	Southwest []dataid    `json:"http://geosysiot.in/rsso/ApplicationSchema#southwest"`
}

type Place struct {
	Url        string `json:"http://schema.org#url"`
	Comment    string `json:"http://www.w3.org/2000/01/rdf-schema#comment"`
	Homepage   string `json:"http://xmlns.com/foaf/0.1/homepage"`
	Property   string `json:"http://schema.org#property"`
	City       string `json:"http://schema.org#city"`
	Level      string `json:"http://schema.org#floorLevel"`
	Latitude   string `json:"http://opengis.net/ont/geosparql#lat"`
	Longitude  string `json:"http://opengis.net/ont/geosparql#long"`
	North      string `json:"http://geosysiot.in/rsso/ApplicationSchema#north"`
	South      string `json:"http://geosysiot.in/rsso/ApplicationSchema#south"`
	East       string `json:"http://geosysiot.in/rsso/ApplicationSchema#east"`
	West       string `json:"http://geosysiot.in/rsso/ApplicationSchema#west"`
	Northeast  string `json:"http://geosysiot.in/rsso/ApplicationSchema#northeast"`
	Northwest  string `json:"http://geosysiot.in/rsso/ApplicationSchema#northwest"`
	Southeast  string `json:"http://geosysiot.in/rsso/ApplicationSchema#southeast"`
	Southwest  string `json:"http://geosysiot.in/rsso/ApplicationSchema#southwest"`
	Directions map[string]string
}

func RaiseErrNotFound(s string) error {
	return fmt.Errorf("place %s does not exist, error: %w", s, ErrNotFound)
}

func (r *Place) UnmarshalJSON(b []byte) error {
	rsp := []jsonMap{}
	err := json.Unmarshal(b, &rsp)
	if err != nil {
		return err
	} else {
		r.Directions = make(map[string]string)
		for _, i := range rsp {
			switch {
			case i.Url != nil:
				r.Url = i.Url[0].Value
			case i.Comment != nil:
				r.Comment = i.Comment[0].Value
			case i.Homepage != nil:
				r.Homepage = i.Homepage[0].Value
			case i.Property != nil:
				r.Property = i.Property[0].Value
			case i.City != nil:
				r.City = i.City[0].Value
			case i.Level != nil:
				r.Level = i.Level[0].Value
			case i.Latitude != nil:
				r.Latitude = i.Latitude[0].Value
			case i.Longitude != nil:
				r.Longitude = i.Longitude[0].Value
			case i.North != nil:
				r.North = i.North[0].Value
				r.Directions["north"] = i.North[0].Value
			case i.South != nil:
				r.South = i.South[0].Value
				r.Directions["south"] = i.South[0].Value
			case i.East != nil:
				r.East = i.East[0].Value
				r.Directions["east"] = i.East[0].Value
			case i.West != nil:
				r.West = i.West[0].Value
				r.Directions["west"] = i.West[0].Value
			case i.Northeast != nil:
				r.Northeast = i.Northeast[0].Value
				r.Directions["northeast"] = i.Northeast[0].Value
			case i.Northwest != nil:
				r.Northwest = i.Northwest[0].Value
				r.Directions["northwest"] = i.Northwest[0].Value
			case i.Southeast != nil:
				r.Southeast = i.Southeast[0].Value
				r.Directions["southeast"] = i.Southeast[0].Value
			case i.Southwest != nil:
				r.Southwest = i.Southwest[0].Value
				r.Directions["southwest"] = i.Southwest[0].Value
			}
		}
	}
	//
	return err
}

func TrimQuotes(s string) string {
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func ConvertGraphString(s string) string {
	return strings.Replace(strings.Replace(TrimQuotes(strings.Split(s, `^^`)[0]), "<", "", 1), ">", "", 1)
}

func NewGraphFromTTL(ttl string) (*rdf2go.Graph, error) {
	baseUri := "https://example.org/foo"
	g := rdf2go.NewGraph(baseUri)

	r := strings.NewReader(ttl)
	err := g.Parse(r, "text/turtle")
	if err != nil {
		return g, err
	}
	return g, err
}

func NewGraphFromSubject(g *rdf2go.Graph, subject rdf2go.Term) *rdf2go.Graph {
	g2 := rdf2go.NewGraph("https://example.org")

	triples3 := g.All(subject, nil, nil)
	for _, triple3 := range triples3 {
		g2.Add(rdf2go.NewTriple(triple3.Subject, triple3.Predicate, triple3.Object))
	}
	return g2
}

func GraphToArea(g *rdf2go.Graph) (Place, error) {
	var err error
	var w bytes.Buffer

	g.Serialize(&w, "application/ld+json")

	var in Place
	err = json.Unmarshal([]byte(w.Bytes()), &in)
	if err != nil {
		return in, err
	}

	return in, err
}

func DirectionsToString(dirs map[string](string)) string {
	ldirs := []string{}
	for _, dir := range ldirections {
		for i := range dirs {
			if i == dir {
				ldirs = append(ldirs, dir)
			}
		}
	}
	out := strings.Join(ldirs, ", ")
	return out
}

func OpenTTLFile(f string) (string, error) {
	var err error
	var t string
	u, _ := url.ParseRequestURI(f)
	if u != nil {
		if u.Scheme == "file" {
			f, err := os.Open(u.Host + u.Path)
			if err != nil {
				return t, err
			}
			defer f.Close()
			bytes, _ := io.ReadAll(f)
			t = string(bytes)
		}
	}
	return t, err
}

func GetPlaceObject(ttl string, placeUri string) (Place, error) {
	var out Place
	g, err := NewGraphFromTTL(ttl)
	if err != nil {
		return out, err
	}

	triples := g.All(nil, nil, rdf2go.NewResource(placeObjectUri))
	for _, triple := range triples {
		if ConvertGraphString(triple.Subject.String()) == placeUri {
			g2 := NewGraphFromSubject(g, triple.Subject)
			out, err := GraphToArea(g2)
			if err != nil {
				return out, err
			}
			return out, err
		}
	}
	err = RaiseErrNotFound(placeUri)
	return out, err
}

func GetPlace(ttl string, placeUri string) (Place, string, error) {
	out, err := GetPlaceObject(ttl, placeUri)
	if err != nil {
		return out, ttl, err
	}

	if len(out.Directions) == 0 {
		if out.Url != "" {
			out, ttl, _ = GetStart(out.Url)
		}
	}

	return out, ttl, nil
}

func GetStartObject(ttl string) (Place, error) {
	var out Place

	g, err := NewGraphFromTTL(ttl)
	if err != nil {
		return out, err
	}

	triples := g.All(nil, rdf2go.NewResource(labelPredicateUri), nil)
	for _, triple := range triples {
		if ConvertGraphString(triple.Object.String()) == startLabel {
			g2 := NewGraphFromSubject(g, triple.Subject)
			out, err := GraphToArea(g2)
			if err != nil {
				return out, err
			}
			return out, err
		}
	}
	return out, err
}

func GetStart(uri string) (Place, string, error) {
	var out Place
	ttl, err := OpenTTLFile(uri)
	if err != nil {
		return out, ttl, err
	}

	out, err = GetStartObject(ttl)
	if err != nil {
		return out, ttl, err
	}
	return out, ttl, err
}

func GetInput() string {
	reader := bufio.NewReader(os.Stdin)
	ui, _ := reader.ReadString('\n')
	return ui
}

func InputToLong(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	var (
		cmd_map = strings.NewReplacer(
			"?", "help",
			"h", "help",
			"q", "quit",
			"n", "north",
			"s", "south",
			"e", "east",
			"w", "west",
		)
		dir_map = strings.NewReplacer(
			"ne", "northeast",
			"nw", "northwest",
			"se", "southeast",
			"sw", "southwest",
		)
	)
	if len(s) < 2 {
		s = cmd_map.Replace(s)
	} else if len(s) < 3 {
		s = dir_map.Replace(s)
	}
	return s
}

func ValidateDirection(dirs map[string](string), in string) error {
	var err error
	for i := range dirs {
		if i == in {
			return err
		}
	}
	return errors.New("not found")
}

func StringToAsciiWhite(s string) string {
	return string("\033[0;37;40m") + s
}

func StringToAsciiBlue(s string) string {
	return string("\033[1;36;40m") + s
}

func StringToAsciiGreen(s string) string {
	return string("\n\033[1;32;40m") + s
}
