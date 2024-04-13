package handlers

import (
	"log"
	"net/http"
	"product-api/product-api/data"
	"regexp"
	"strconv"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

func (p Product) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}
	if r.Method == http.MethodPut {
		p.l.Println("PUT", r.URL.Path)
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
		p.l.Println("value of g ", g)
		if len(g) != 1 {
			p.l.Panicln("Invalid URI more than one id")
			http.Error(rw, "Invalid URI more than one id", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			p.l.Panicln("Invalid URI more than one capture group")
			http.Error(rw, "Invalid URI more than one capture group", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("unable to convert to integer")
			http.Error(rw, "unable to convert to integer", http.StatusBadRequest)
		}
		p.l.Println("got id ", id)
		p.updateProduct(id, rw, r)
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p Product) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marhshal json", http.StatusInternalServerError)
		return
	}
}
func (p Product) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("handle Post request..")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarhshal json", http.StatusBadRequest)
		return
	}
	p.l.Printf("Prod: %v", prod)
	data.AddProducts(prod)
}

func (p Product) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("handling update request")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarhshal json", http.StatusBadRequest)
		return
	}
	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusBadRequest)
	}
	if err != nil {
		http.Error(rw, "Product not found", http.StatusBadRequest)
	}

}
