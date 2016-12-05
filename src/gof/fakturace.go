package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

type fakturace struct {
	XMLName xml.Name `json:"-" xml:"winstrom"`
	Version string   `json:"-" xml:"version,attr"`
	Source  string   `json:"-" xml:"source,attr"`

	Smlouvy []smlouva `json:"-" xml:"faktura-vydana"`
}

func (f *fakturace) Fakturuj(now time.Time) (fakturace, error) {
	// zkopirovat f do noveho objektu, vyfakturovat v nem
	var faktury fakturace

	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	enc.Encode(f)
	dec.Decode(&faktury)

	for key, _ := range faktury.Smlouvy {
		if err := faktury.Smlouvy[key].Fakturuj(now); err != nil {
			return faktury, err
		}

		// aktualizace fakturovano do v puvodnim objektu fakturace (pro ulozeni do json)
		for kp, _ := range faktury.Smlouvy[key].Polozky {
			f.Smlouvy[key].Polozky[kp].VyfakturovanoDo = faktury.Smlouvy[key].Polozky[kp].VyfakturovanoDo
		}
	}

	deleted := 0
	for k, _ := range faktury.Smlouvy {
		key := k - deleted
		if len(faktury.Smlouvy[key].Polozky) == 0 { // bez polozek nefakturujeme
			log.Printf("Skipping SMLOUVA %s\n", faktury.Smlouvy[key].Zakazka)
			faktury.Smlouvy = append(faktury.Smlouvy[:key], faktury.Smlouvy[key+1:]...)
			deleted++
		}

	}

	return faktury, nil
}

func (f *fakturace) XML() (string, error) {
	f.Version = "1.0"
	f.Source = "xert fakturace"

	for ks, _ := range f.Smlouvy {
		f.Smlouvy[ks].Zakazka = "code:" + f.Smlouvy[ks].Zakazka
		f.Smlouvy[ks].TypDokladu = "code:" + f.Smlouvy[ks].TypDokladu

		for kp, _ := range f.Smlouvy[ks].Polozky {
			f.Smlouvy[ks].Polozky[kp].CisloRadku = kp + 1
			f.Smlouvy[ks].Polozky[kp].Cenik = "code:" + f.Smlouvy[ks].Polozky[kp].Cenik
		}
	}

	output, err := xml.MarshalIndent(f, "  ", "    ")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s\n", xml.Header, output), nil
}

func (f *fakturace) Save(filename string) error {
	var err error
	var data []byte

	if data, err = json.MarshalIndent(f.Smlouvy, "", "  "); err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func NewFakturace(jsondata []byte) (f fakturace, err error) {

	if err = json.Unmarshal(jsondata, &f.Smlouvy); err != nil {
		return
	}

	for k, s := range f.Smlouvy {
		f.Smlouvy[k].Nazev = s.Zakazka

		for _, p := range s.Polozky {
			if moreThanOneNonZero(p.Frekvence.Years, p.Frekvence.Months, p.Frekvence.Days) {
				err = fmt.Errorf("Error: Smlouva %s, polozka '%s', frekvence %d,%d,%d - more than one field is non zero", s.Zakazka, p.Nazev, p.Frekvence.Years, p.Frekvence.Months, p.Frekvence.Days)
				return
			}
		}
	}

	return
}

func moreThanOneNonZero(items ...int) bool {
	nonzeroes := 0

	for _, i := range items {
		if i > 0 {
			nonzeroes++
		}

		if nonzeroes > 1 {
			return true
		}
	}

	return false
}
