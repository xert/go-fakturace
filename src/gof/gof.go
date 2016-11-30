package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"
	"time"
)

const (
	filename string = "./smlouvy.json"
)

func main() {
	var y []byte
	var f, faktury fakturace
	var err error
	var xml string
	var now time.Time

	if len(os.Args) == 2 {
		if now, err = time.Parse(dateFormat, os.Args[1]); err != nil {
			panic(err)
		}
	} else {
		now = time.Now()
	}

	if y, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}

	if f, err = NewFakturace(y); err != nil {
		panic(err)
	}

	if faktury, err = f.Fakturuj(now); err != nil {
		panic(err)
	}

	if xml, err = faktury.XML(); err != nil {
		panic(err)
	}

	fmt.Print(xml)
	if err = f.Save(filename + ".new"); err != nil {
		panic(err)
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
