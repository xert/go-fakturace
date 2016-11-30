package main

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateFormat string = "2.1.2006"

type interval struct {
	Od time.Time
	Do time.Time
}

func (i *interval) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	date := string(b)

	i.Od = time.Time{}
	i.Do = time.Time{}

	if date != "" {
		parts := strings.Split(date, "-")
		for key, _ := range parts {
			parts[key] = strings.TrimSpace(parts[key])
		}

		if len(parts) != 2 {
			return fmt.Errorf("Spatny pocet datumu v intervalu %s. Ocekavam format dd.mm.yyyy - dd.mm.yyyy", date)
		}

		if parts[0] != "" {
			if i.Od, err = time.Parse(dateFormat, parts[0]); err != nil {
				return err
			}
		}
		if parts[1] != "" {
			if i.Do, err = time.Parse(dateFormat, parts[1]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *interval) MarshalJSON() ([]byte, error) {
	if i.Od.IsZero() && i.Do.IsZero() {
		return []byte("\"\""), nil
	}

	var od, do string

	if !i.Od.IsZero() {
		od = i.Od.Format(dateFormat)
	}
	if !i.Do.IsZero() {
		do = i.Do.Format(dateFormat)
	}
	return []byte(fmt.Sprintf("\"%s - %s\"", od, do)), nil
}

type datum struct {
	time.Time
}

func (d *datum) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	date := string(b)
	if date == "" {
		return
	}
	d.Time, err = time.Parse(dateFormat, date)
	return
}

func (d *datum) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("\"\""), nil
	}
	return []byte("\"" + d.Time.Format(dateFormat) + "\""), nil
}

func (d *datum) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	return e.EncodeElement(d.Format("2006-01-02"), start)
}

type frekvence struct {
	Years  int
	Months int
	Days   int
}

func (f *frekvence) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	tmp := strings.Split(string(b), ",")
	if f.Years, err = strconv.Atoi(tmp[0]); err != nil {
		return
	}
	if f.Months, err = strconv.Atoi(tmp[1]); err != nil {
		return
	}
	if f.Days, err = strconv.Atoi(tmp[2]); err != nil {
		return
	}

	return
}

func (f *frekvence) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%d,%d,%d\"", f.Years, f.Months, f.Days)), nil
}

func (f *frekvence) IsZero() bool {
	return f.Years == 0 && f.Months == 0 && f.Days == 0
}

func (f frekvence) String() string {
	return fmt.Sprintf("%d,%d,%d", f.Years, f.Months, f.Days)
}
