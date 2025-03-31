package utils

import "time"

type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(rawDate []byte) (err error) {
	date, err := time.Parse(`"2006-01-02T15:04:05.000"`, string(rawDate))
	if err != nil {
		date, err = time.Parse(`"2006-01-02T15:04:05.000Z"`, string(rawDate))
		if err != nil {
			return err
		}
	}
	t.Time = date
	return
}
