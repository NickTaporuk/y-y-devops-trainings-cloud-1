package gpt

import (
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type IndexData struct {
	Error error
	Photo string
}
type Cat int

const (
	CatUnknown Cat = iota
	CatNocturnal
	CatDiurnal
)

func (c Cat) String() string {
	switch c {
	case CatNocturnal:
		return "nocturnal"
	case CatDiurnal:
		return "diurnal"
	}
	return "unknown"
}

var (
	//go:embed tpl/index.html
	f        embed.FS
	indexTpl = template.Must(template.ParseFS(f, "tpl/index.html"))
)

type Handler struct {
	gpt *CatGPT
}

func NewHandler(gpt *CatGPT) *Handler {
	return &Handler{gpt: gpt}
}

func (handler *Handler) Gpt() *CatGPT {
	return handler.gpt
}

func (handler *Handler) index(w http.ResponseWriter, r *http.Request) {
	id := &IndexData{}
	if r.Method != http.MethodPost {
		indexTpl.Execute(w, id)
		return
	}
	file, _, err := r.FormFile("photo")
	if err != nil {
		id.Error = err
		return
	}
	defer file.Close()
	var img image.Image
	img, err = handler.Gpt().EnsureIsImage(file)
	if err != nil {
		id.Error = err
		return
	}
	catType := CatDiurnal
	if rand.Intn(100) > 75 {
		catType = CatNocturnal
	}
	var enhanced io.Reader
	enhanced, err = handler.Gpt().Enhance(img, catType)
	if err != nil {
		id.Error = err
		return
	}
	var enhancedBytes []byte
	enhancedBytes, err = io.ReadAll(enhanced)
	if err != nil {
		id.Error = err
		return
	}

	sb := &strings.Builder{}
	encoder := base64.NewEncoder(base64.StdEncoding, sb)
	_, err = encoder.Write(enhancedBytes)
	if err != nil {
		id.Error = err
		return
	}
	encoder.Close()
	id.Photo = sb.String()

	catCounterVec.WithLabelValues(catType.String()).Inc()
	indexTpl.Execute(w, id)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

var (
	responseCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_count",
			Help: "by handler and code",
		},
		[]string{"code", "handler", "method"})
	catCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "enhanced_photo_count",
			Help: "by cat type",
		},
		[]string{"cat_type"})
)

func Serve(ctx context.Context, public string, private string, handler *Handler) {
	fmt.Println("Starting the server ...")

	indexChain := promhttp.InstrumentHandlerCounter(
		responseCounterVec.MustCurryWith(prometheus.Labels{"handler": "/"}),
		http.HandlerFunc(handler.index))
	pingChain := promhttp.InstrumentHandlerCounter(
		responseCounterVec.MustCurryWith(prometheus.Labels{"handler": "/ping"}),
		http.HandlerFunc(ping))
	prometheus.MustRegister(responseCounterVec, catCounterVec)

	http.HandleFunc("/", indexChain)
	http.HandleFunc("/ping", pingChain)

	go func() {
		log.Fatal(http.ListenAndServe(private, promhttp.Handler()))
	}()
	log.Fatal(http.ListenAndServe(public, nil))
}
