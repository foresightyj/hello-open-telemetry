package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const name = "weather-internal"

type WeatherData struct {
	Date         time.Time
	TemperatureC int
	TemprartureF int
	Summary      string
}

type ApiResponse[T any] struct {
	Success bool
	Error   string
	Data    T
}

var summaries = []string{"Freezing", "Bracing", "Chilly", "Cool", "Mild", "Warm", "Balmy", "Hot", "Sweltering", "Scorching"}

func getWeathers(limit int) []WeatherData {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	if limit <= 0 {
		log.Fatalln("limit not positive")
	}
	var res []WeatherData
	for i := 0; i < limit; i++ {
		r := r1.Int()
		sum := summaries[r%len(summaries)]
		res = append(res, WeatherData{
			Date:         time.Date(2022, 10, r%31, 0, 0, 0, 0, time.UTC),
			TemperatureC: r % 40,
			TemprartureF: 32 + (int)((float64(r%40))/0.5556),
			Summary:      sum,
		})
	}
	return res
}

var tracer = otel.Tracer(name)

func procastinate(ctx context.Context, ms time.Duration) {
	_, span := tracer.Start(ctx, "Procastinate...")
	defer span.End()
	time.Sleep(ms * time.Millisecond)
}

func initTracer() (*sdktrace.TracerProvider, error) {
	url := "http://localhost:14268/api/traces"
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	//exporter, err := stdout.New(stdout.WithPrettyPrint())

	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func main() {
	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	r := gin.Default()
	r.Use(otelgin.Middleware("weather-go-server"))

	r.GET("/weather", func(c *gin.Context) {
		limit, err := strconv.ParseInt(c.DefaultQuery("limit", "5"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, ApiResponse[[]WeatherData]{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		ctx := c.Request.Context()

		procastinate(ctx, 1500)

		_, span := tracer.Start(ctx, "CalculateWeatherForecasts")
		defer span.End()

		j, _ := json.MarshalIndent(c.Request.Header, "", "  ")
		fmt.Println("========== REQ ===========", string(j))
		weathers := getWeathers(int(limit))
		res := ApiResponse[[]WeatherData]{
			Success: true,
			Error:   "",
			Data:    weathers,
		}
		time.Sleep(1 * time.Second)
		c.JSON(http.StatusOK, res)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
