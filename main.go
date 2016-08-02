package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/vsukhin/product/controllers"
	"github.com/vsukhin/product/logger"
	"github.com/vsukhin/product/repositories"
)

const (
	// PortHTTP contains server port
	PortHTTP = 3000
	// ParameterNamePortHTTP contains parameter name
	ParameterNamePortHTTP = "port"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var (
	httpPort = flag.Int(ParameterNamePortHTTP, PortHTTP, "HTTP server port")
)

func main() {
	logger.Log.Println("Starting server work at ", time.Now())
	flag.Parse()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		close(c)
		logger.Log.Println("Finishing server work at ", time.Now())
		os.Exit(1)
	}()

	productController := controllers.NewProductControllerImplementation(repositories.NewProductRepositoryImplementation())

	router := mux.NewRouter()

	router.HandleFunc("/products", productController.List).Methods("GET")
	router.HandleFunc("/products", productController.Create).Methods("POST")
	router.HandleFunc("/products/{id:[0-9]+}", productController.Get).Methods("GET")
	router.HandleFunc("/products/{id:[0-9]+}", productController.Update).Methods("PUT")
	router.HandleFunc("/products/{id:[0-9]+}", productController.Delete).Methods("DELETE")
	router.HandleFunc("/products/{id:[0-9]+}/prices", productController.SetPrices).Methods("PUT")

	err := http.ListenAndServe(":"+strconv.Itoa(*httpPort), router)
	if err != nil {
		logger.Log.Fatalf("Can't listen http port %v having error %v", httpPort, err)
	}
}
