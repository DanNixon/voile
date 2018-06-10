package db

import (
	"encoding/json"
	"net/url"
)

type Url struct {
	Url url.URL
}

func (u *Url) Parse(urlStr string) error {
	uu, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	u.Url = *uu
	return nil
}

func (u *Url) String() string {
	a := u.Url.String()
	return a
}

func (u *Url) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *Url) UnmarshalJSON(b []byte) error {
	var urlStr string
	json.Unmarshal(b, &urlStr)
	return u.Parse(urlStr)
}
