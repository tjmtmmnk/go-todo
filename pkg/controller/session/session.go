package session

import (
	"errors"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const (
	defaultSessionIDKey = "SESSIONID"
	defaultUserIDKey    = "USERID"
)

type Session struct {
	UserID uint64
}

func ExtractFromContext(c echo.Context) (*Session, error) {
	sess, err := session.Get(defaultSessionIDKey, c)
	if err != nil {
		return nil, err
	}
	userID, ok := sess.Values[defaultUserIDKey].(uint64)
	if !ok {
		return nil, errors.New("failed to cast session userID")
	}
	return &Session{UserID: userID}, nil
}

func (s *Session) Save(c echo.Context, options *sessions.Options) error {
	sess, err := session.Get(defaultSessionIDKey, c)
	if err != nil {
		return err
	}
	sess.Options = options

	sess.Values[defaultUserIDKey] = s.UserID

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}

	return nil
}
