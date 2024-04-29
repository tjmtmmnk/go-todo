package timex

import "time"

var JST *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	JST = loc
}
