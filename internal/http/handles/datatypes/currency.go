package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Currency struct {
	Name string // string representation of currency

	IsKnown bool // note: IsKnown is false when Name is an empty string ("")
}

func (c *Currency) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &c.Name); err != nil {
		return err
	}

	for _, currency := range currencies {
		if c.Name == currency {
			c.IsKnown = true
			break
		}
	}
	return nil
}

func (c *Currency) String() string { // todo: check where used
	return c.Name
}

func (c *Currency) Scan(src interface{}) (err error) {
	fmt.Println(src)
	return nil
}

func (c Currency) Value() (driver.Value, error) {
	return c.Name, nil
}
