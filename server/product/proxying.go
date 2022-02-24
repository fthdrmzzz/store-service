package product

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

type proxymw struct {
	ctx        context.Context
	next       ProductServer     // Serve most requests via this service...
	addProduct endpoint.Endpoint // ...except Uppercase, which gets served by this endpoint
}

func (mw proxymw) DeleteProduct(productId int) error {
	return mw.next.DeleteProduct(productId)
}

func (mw proxymw) GetAllProducts() error {
	return mw.next.GetAllProducts()
}

func (mw proxymw) AddProduct(userId int, name string, price int16) (err error) {
	request := addProductRequest{
		UserId: userId,
		Name:   name,
		Price:  price}
	response, err := mw.addProduct(mw.ctx, request)

	if err != nil {
		return err
	}
	resp := response.(addProductResponse)
	if resp.Err != "" {
		return errors.New(resp.Err)
	}
	return nil
}

func ProxyingMiddleware(ctx context.Context, instances string, logger log.Logger) ServiceMiddleware {

	var (
		qps         = 100
		maxAttempts = 3
		maxTime     = 250 * time.Millisecond
	)

	var (
		instanceList = split(instances)
		endpointer   sd.FixedEndpointer
	)

	logger.Log("proxy_to", fmt.Sprint(instanceList))

	for _, instance := range instanceList {
		var e endpoint.Endpoint
		e = makeAddProductProxy(ctx, instance)
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e)
		endpointer = append(endpointer, e)
	}

	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)

	return func(next ProductServer) ProductServer {
		return proxymw{ctx, next, retry}
	}
}

func makeAddProductProxy(ctx context.Context, instance string) endpoint.Endpoint {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	fmt.Println("\n\nPRINTING", instance, "\n\n")
	u, err := url.Parse(instance)

	if err != nil {
		panic(err)
	}
	if u.Path == "" {
		u.Path = "/addProduct"
	}
	return httptransport.NewClient(
		"GET",
		u,
		EncodeRequest,
		DecodeAddProductResponse,
	).Endpoint()
}
func split(s string) []string {
	a := strings.Split(s, ",")
	for i := range a {
		a[i] = strings.TrimSpace(a[i])
	}
	return a
}

type Subscriber interface {
	Endpoints() ([]endpoint.Endpoint, error)
}
