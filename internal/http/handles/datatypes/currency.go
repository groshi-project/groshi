package datatypes

import "encoding/json"

type Currency struct {
	String string // string representation of currency

	IsKnown bool // note: IsKnown is true even when currency is an empty string ("")
}

func (c *Currency) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &c.String)
	if err != nil {
		c.IsKnown = false
	} else {
		c.IsKnown = true
	}
	return nil
}
