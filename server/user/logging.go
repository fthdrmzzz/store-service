package user

import (
	"time"

	"github.com/go-kit/log"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next   UserServer
}

func (mw LoggingMiddleware) CreateUser(email string, name string, password string) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Log(
			"method", "create user",
			"name", name,
			"email", email,
			"password", password,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.Next.CreateUser(email, name, password)
	return
}
