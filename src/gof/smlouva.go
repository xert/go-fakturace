package main

import (
	"log"
	"time"
)

type smlouva struct {
	Nazev          string    `json:"-" xml:"cisSml"`
	Zakazka        string    `json:"zakazka" xml:"zakazka"`
	DatumVystaveni datum     `json:"-"  xml:"datVyst"`
	DatumDUZP      datum     `json:"-"  xml:"duzpPuv"`
	Popis          string    `json:"popis" xml:"popis"`
	Platnost       interval  `json:"platnost" xml:"-"`
	Frekvence      frekvence `json:"frekvence" xml:"-"`
	TypDokladu     string    `json:"typ_dokladu" xml:"typDokl"`

	Polozky []polozka `json:"polozky" xml:"polozkyFaktury>faktura-vydana-polozka"`
}

func (s *smlouva) Fakturuj(now time.Time) error {
	s.DatumVystaveni = datum{now}
	s.DatumDUZP = s.DatumVystaveni

	deleted := 0
	for k, _ := range s.Polozky {
		key := k - deleted

		if err := s.Polozky[key].Fakturuj(now, s); err != nil {
			return err
		}

		if s.Polozky[key].Mnozstvi == 0 { // nefakturovat
			log.Printf("Skipping polozka %s\n", s.Polozky[key].Nazev)
			s.Polozky = append(s.Polozky[:key], s.Polozky[key+1:]...)
			deleted++
		}

	}

	return nil
}
