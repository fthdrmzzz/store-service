package product

import (
	"time"

	"github.com/go-kit/log"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next   ProductServer
}

func (mw LoggingMiddleware) AddProduct(userId int, name string, price int16) (err error) {
	defer func(begin time.Time) {
		mw.Logger.Log(
			"method", "AddProduct",
			"userId", userId,
			"name", name,
			"price", price,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Next.AddProduct(userId, name, price)
	return
}

func (mw LoggingMiddleware) DeleteProduct(productId int) (err error) {
	defer func(begin time.Time) {
		mw.Logger.Log(
			"method", "DeleteProduct",
			"productId", productId,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.Next.DeleteProduct(productId)
	return
}

func (mw LoggingMiddleware) GetAllProducts() (err error) {
	defer func(begin time.Time) {
		mw.Logger.Log(
			"method", "GetAllProducts",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.Next.GetAllProducts()
	return
}
