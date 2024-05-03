package timex

import (
	"github.com/labstack/echo/v4"
	"github.com/moznion/go-optional"
	"time"
)

var JST *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	JST = loc
}

type OptionalDate optional.Option[time.Time]
type OptionalDateTime optional.Option[time.Time]

var _ echo.BindUnmarshaler = (*OptionalDate)(nil)
var _ echo.BindUnmarshaler = (*OptionalDateTime)(nil)

func UnwrapAsUTCPtr(t optional.Option[time.Time]) *time.Time {
	return optional.Map(t, func(_t time.Time) time.Time { return _t.UTC() }).UnwrapAsPtr()
}

func (od *OptionalDate) UnmarshalParam(src string) error {
	if src == "" {
		return nil
	}
	parsedTime, err := time.Parse(time.DateOnly, src)
	if err != nil {
		return err
	}

	*od = OptionalDate(optional.Some(parsedTime))

	return nil
}

func (odt *OptionalDateTime) UnmarshalParam(src string) error {
	if src == "" {
		return nil
	}
	parsedTime, err := time.Parse(time.RFC3339, src)
	if err != nil {
		return err
	}

	*odt = OptionalDateTime(optional.Some(parsedTime))

	return nil
}
