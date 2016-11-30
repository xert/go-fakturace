package main

import (
	"fmt"
	"time"
)

type polozka struct {
	CisloRadku      int       `json:"-" xml:"cisRad"`
	Nazev           string    `json:"nazev" xml:"nazev"`
	VyfakturovanoDo datum     `json:"vyfakturovano_do" xml:"-"`
	Cenik           string    `json:"cenik" xml:"cenik"`
	Mnozstvi        float32   `json:"mnozstvi" xml:"mnozMj"`
	Cena            float32   `json:"cena" xml:"cenaMj"`
	Sleva           float32   `json:"sleva,omitempty" xml:"slevaPol"`
	Platnost        *interval `json:"platnost,omitempty"  xml:"-"`
	Frekvence       frekvence `json:"frekvence" xml:"-"`
}

func (p *polozka) Fakturuj(now time.Time, s *smlouva) error {
	if p.Platnost != nil {
		if p.Platnost.Od.IsZero() {
			p.Platnost.Od = s.Platnost.Od
		}
		if p.Platnost.Do.IsZero() {
			p.Platnost.Do = s.Platnost.Do
		}
	}
	if p.VyfakturovanoDo.IsZero() {
		// prazdne vyfakturovano do znaci novou polozku
		p.VyfakturovanoDo = datum{now.AddDate(-p.Frekvence.Years, -p.Frekvence.Months, -p.Frekvence.Days)}
	}
	if p.Frekvence.IsZero() {
		p.Frekvence = s.Frekvence
	}

	if p.Mnozstvi == 0 {
		p.Mnozstvi = 1
	}

	startDate := p.VyfakturovanoDo.AddDate(0, 0, 1)
	endDate := startDate

	// plati pro fakturaci dopredu, s jinou zatim nepocitame
	count := 0
	for !endDate.After(now) {
		endDate = endDate.AddDate(p.Frekvence.Years, p.Frekvence.Months, p.Frekvence.Days)
		count++
	}

	endDate = endDate.AddDate(0, 0, -1)

	p.Nazev = fmt.Sprintf("%s za %d období (%s – %s)", p.Nazev, count, startDate.Format("2.1.2006"), endDate.Format("2.1.2006"))
	p.Mnozstvi *= float32(count)
	p.VyfakturovanoDo = datum{endDate}

	return nil
}
