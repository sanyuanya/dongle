package entity

import "encoding/json"

type Base64String string

func (b Base64String) MarshalJSON() ([]byte, error) {
	return json.Marshal("")
}

func (b *Base64String) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*b = Base64String(s)
	return nil
}

func (b Base64String) String() string {
	return string(b)
}
