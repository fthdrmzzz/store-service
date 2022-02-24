package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	pkProduct "github.com/fthdrmzzz/store-service/server/product"
	pkUser "github.com/fthdrmzzz/store-service/server/user"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//admin
//uD5zvmnN92QeTwLM

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
		proxy  = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy uppercase requests")
	)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "user_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "user_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "user_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	var userSvc pkUser.UserServer
	userSvc = pkUser.UserService{}

	userSvc = pkUser.LoggingMiddleware{logger, userSvc}
	userSvc = pkUser.InstrumentingMiddleware{requestCount, requestLatency, countResult, userSvc}
	createUserHandler := httptransport.NewServer(
		pkUser.MakeCreateUserEndpoint(userSvc),
		pkUser.DecodeCreateUserRequest,
		pkUser.EncodeResponse,
	)

	requestCount = kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "product_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency = kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "product_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult = kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "product_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	var productSvc pkProduct.ProductServer
	productSvc = pkProduct.ProductService{}
	productSvc = pkProduct.ProxyingMiddleware(context.Background(), *proxy, logger)(productSvc)
	productSvc = pkProduct.LoggingMiddleware{logger, productSvc}
	productSvc = pkProduct.InstrumentingMiddleware{requestCount, requestLatency, countResult, productSvc}
	addProductHandler := httptransport.NewServer(
		pkProduct.MakeAddProductEndpoint(productSvc),
		pkProduct.DecodeAddProductRequest,
		pkUser.EncodeResponse,
	)

	deleteProductHandler := httptransport.NewServer(
		pkProduct.MakeDeleteProductEndpoint(productSvc),
		pkProduct.DecodeDeleteProductRequest,
		pkUser.EncodeResponse,
	)

	http.Handle("/createUser", createUserHandler)
	http.Handle("/addProduct", addProductHandler)
	http.Handle("/deleteProduct", deleteProductHandler)

	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}
