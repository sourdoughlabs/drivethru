package main

import (
    "fmt"
    "os"
    "encoding/json"
    "io/ioutil"
	"time"
)

/* 
 * TODO Move these into a package to be imported.  
 * They are from my mkinvoice program and this a bad way to do it 
 */

type Invoice struct {
    Number           string
    Date             time.Time

	ClientName       string
	ClientContact    string
	ClientEmail      string
	ClientStreet     string
	ClientCity       string
	ClientProvince   string
	ClientPostal     string

    Note             string

    LineItems        []*Item
    Payments         []*Payment
}

type Item struct {
    Description  string
    Rate    int // In pennies - Fixed point is so much better for this.
    Quantity int
	ProductCode string
	TaxRate int // ie 5 for 5%
}

type Payment struct {
	Amount int
	Date time.Time
}

func (invoice *Invoice) Checked(code string) string {

	var checked = " "

	for _, item := range invoice.LineItems {
		if code == item.ProductCode {
			checked = "checked"
			break;
		}
	}
	return checked
}

func (invoice *Invoice) Valid() (invalid bool, errors []string) {
	invalid = false
	errors = make([]string, 0, 10)

	if len(invoice.ClientName) == 0 {
		invalid = true
		errors = append(errors, "Company name is required")
	}

	if len(invoice.ClientContact) == 0 {
		invalid = true
		errors = append(errors, "Name is required")
	}

	if len(invoice.ClientEmail) == 0 {
		invalid = true
		errors = append(errors, "Email is required")
	}

	if len(invoice.ClientStreet) == 0 {
		invalid = true
		errors = append(errors, "Street is required")
	}

	if len(invoice.ClientCity) == 0 {
		invalid = true
		errors = append(errors, "City is required")
	}

	if len(invoice.ClientProvince) == 0 {
		invalid = true
		errors = append(errors, "State/Province is required")
	}

	if len(invoice.ClientPostal) == 0 {
		invalid = true
		errors = append(errors, "Zip/Postal Code is required")
	}	

	if len(invoice.LineItems) == 0 {
		invalid = true
		errors = append(errors, "You must select at least one product")
	}	
	
	return !invalid, errors
}

func SaveInvoice(filename string, invoice *Invoice) {

	invoiceB, _ := json.Marshal(invoice)

    err := ioutil.WriteFile(filename, invoiceB, 0644)

    if err != nil {
        fmt.Printf("File error: %v\n", err)
        os.Exit(1)
    }
}
