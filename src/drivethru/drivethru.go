/*
 * drivethru.go - An extremely simple order taking application.
 *
 * License: MIT
 */

package main

import (
	"fmt"
	"strconv"
	"flag"
	"time"
	"html/template"
	"net/http"
	"log"
)

type OrderFormContext struct {
	Flash string
	Invoice *Invoice
	Products Products
    Errors []string
	Success string
}

func main() {

	port := flag.String("port", "14002", "Port number to listen on")

    flag.Parse()

	// Authorization code endpoint
	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/html")

		invoice := &Invoice{}
		invoice.Date = time.Now()
		invoice.Number = strconv.FormatInt(invoice.Date.UnixNano(), 10)
		invoice.LineItems = make([]*Item, 0, 10)

		products := load_products("products.json") // Load on demand so we don't have to bounce the server to have new products

		context := &OrderFormContext{Flash: "", Invoice: invoice, Products: products}

		var valid bool
		if r.Method == "POST"  {

			r.ParseForm()

			invoice.ClientName =  r.Form.Get("company")
			invoice.ClientContact = r.Form.Get("name")
			invoice.ClientEmail = r.Form.Get("email")
			invoice.ClientStreet = r.Form.Get("street")
			invoice.ClientCity  = r.Form.Get("city")
			invoice.ClientProvince = r.Form.Get("state")
			invoice.ClientPostal = r.Form.Get("zip_code")			
			invoice.Note = r.Form.Get("purchase_order_number")

			for _, code := range r.Form["product"] {
				p := products.Find(code)
				if p != nil {
					invoice.LineItems = append(invoice.LineItems, 
						&Item{Description: p.Description, Rate: p.Amount, Quantity: 1, ProductCode: code, TaxRate: 0})
				}
			}

			valid, context.Errors = invoice.Valid()
			if valid {
				fm := []byte("Thank you for your order. An invoice will be prepared and sent to your email address.")
				SetFlash(w, "success", fm)

				SaveInvoice(fmt.Sprintf("invoices/%s.json", invoice.Number), invoice)

				http.Redirect(w, r, "/order", 302)
			} else {
				context.Flash = "There was a problem processing your order.  Please correct the following and try again:"
				render_order_form(w, context)
			}
		} else {
			fm, err := GetFlash(w, r, "success")
			if err != nil {
				log.Printf("Error getting flash: %v", err.Error())
			}
			context.Success = string(fm)
			render_order_form(w, context)
		}
	})

	fs := http.FileServer(http.Dir("site"))
	http.Handle("/", http.StripPrefix("/", fs))

	http.ListenAndServe(":" + *port, nil)
}

func render_order_form(w http.ResponseWriter, context *OrderFormContext) {
	var t = template.Must(template.New("order_form.html").ParseFiles("site/order_form.html"))
	t.Execute(w, context)
}
