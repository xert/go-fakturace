package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/tabwriter"
	"time"
)

func main() {
	var y []byte
	var f, faktury fakturace
	var err error
	var xml, filename string
	var now time.Time

	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Printf("Usage: %s input.json [1.12.2016] > output.xml\n\nWhere 1.12.2016 is invoice date. Otherwise, current date is used.\n\n", os.Args[0])
		os.Exit(1)
	}

	filename = os.Args[1]

	if len(os.Args) == 3 {
		if now, err = time.Parse(dateFormat, os.Args[2]); err != nil {
			log.Fatal(err)
		}
	} else {
		now = time.Now()
	}

	if y, err = ioutil.ReadFile(filename); err != nil {
		log.Fatal(err)
	}

	if f, err = NewFakturace(y); err != nil {
		log.Fatal(err)
	}

	if faktury, err = f.Fakturuj(now); err != nil {
		log.Fatal(err)
	}

	if xml, err = faktury.XML(); err != nil {
		log.Fatal(err)
	}

	fmt.Print(xml)
	if err = f.Save(filename); err != nil {
		log.Fatal(err)
	}

	w := tabwriter.NewWriter(os.Stderr, 40, 8, 1, ' ', tabwriter.AlignRight)
	for _, s := range faktury.Smlouvy {
		var suma float32 = 0.0
		for _, p := range s.Polozky {
			suma += p.Mnozstvi * p.Cena
		}

		fmt.Fprintf(w, "%s\t%s\t%.2f\t\n", s.Nazev, s.Popis, suma)
		w.Flush()
	}
}
