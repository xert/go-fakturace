package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/kr/pretty"
)

func parseSampleData() (fakturace, error) {
	sampleData := `
[
  {
    "zakazka": "TEST1",
	"typ_dokladu": "TYPE1",
    "popis": "Test1",
    "platnost": "01.01.2016 - ",
    "frekvence": "0,1,0",
    "polozky": [
      {
        "nazev": "Test 1a",
        "cenik": "CEN1A",
        "mnozstvi": 1,
        "cena": 101,
        "platnost": "01.01.2016 -",
        "vyfakturovano_do": "30.11.2016",
        "frekvence": "0,1,0"
      },
      {
        "nazev": "Test 1b",
        "cenik": "CEN1B",
        "mnozstvi": 1,
        "cena": 102,
        "platnost": "01.01.2016 -",
        "vyfakturovano_do": "30.11.2016",
        "frekvence": "0,1,0"
      }
    ]
  },
  {
    "zakazka": "TEST2",
	"typ_dokladu": "TYPE2",
    "popis": "Test2",
    "platnost": "01.01.2016 - 31.12.2016",
    "frekvence": "0,1,0",
	"sleva": 20,
    "polozky": [
      {
        "nazev": "Test 2a",
		"vyfakturovano_do": "30.11.2016",
        "cenik": "CEN2A",
        "mnozstvi": 2,
        "cena": 201
      },
      {
        "nazev": "Test 2b",
		"vyfakturovano_do": "30.11.2016",
        "cenik": "CEN2B",
        "cena": 202,
		"sleva": 22
      }
    ]
  }
]
`

	return NewFakturace([]byte(sampleData))
}

func TestSampleData(t *testing.T) {
	f, err := parseSampleData()

	if err != nil {
		t.Fatal(err)
	}

	if len(f.Smlouvy) == 0 {
		t.Fatal("Expected some sample data")
	}
}

func TestSampleSmlouvyZakazka(t *testing.T) {
	f, _ := parseSampleData()

	for k, s := range f.Smlouvy {
		expected := fmt.Sprintf("TEST%d", k+1)
		if s.Zakazka != expected {
			t.Fatalf("Expected %s got %s", expected, s.Zakazka)
		}
	}

}

func TestIntervalMany(t *testing.T) {
	j := `[{ "platnost": "01.01.2016 - 31.12.2016 - 1.1.2017"}]`

	if _, err := NewFakturace([]byte(j)); err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestIntervalInvalidOd(t *testing.T) {
	j := `[{ "platnost": "A -"}]`

	if _, err := NewFakturace([]byte(j)); err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestIntervalInvalidDo(t *testing.T) {
	j := `[{ "platnost": " - B"}]`

	if _, err := NewFakturace([]byte(j)); err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestSmlouvaZakazka(t *testing.T) {
	s := []smlouva{
		smlouva{Zakazka: "TEST1", Nazev: "TEST1"},
		smlouva{Zakazka: "TEST2", Nazev: "TEST2"},
		smlouva{Zakazka: "TEST3", Nazev: "TEST3"},
	}

	expected := fakturace{Smlouvy: s}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}

	got, err := NewFakturace(data)
	if err != nil {
		t.Fatalf("error: %s, data = %s", err, data)
	}

	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("Expected: %#v\nGot: %#v", expected, got)
	}
}

func TestStruct(t *testing.T) {
	s := []smlouva{
		smlouva{
			Zakazka:    "TEST1",
			Nazev:      "TEST1",
			Popis:      "Test1",
			TypDokladu: "TYPE1",
			Platnost: interval{
				Od: time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC),
				Do: time.Date(2016, time.December, 31, 0, 0, 0, 0, time.UTC),
			},
			Polozky: []polozka{
				polozka{
					Nazev:           "Test 1a",
					Cenik:           "CEN1A",
					Mnozstvi:        1,
					Cena:            101,
					Sleva:           0,
					VyfakturovanoDo: datum{time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)},
					Frekvence:       frekvence{Years: 0, Months: 1, Days: 0},
				},
				polozka{
					Nazev:           "Test 1b",
					Cenik:           "CEN1B",
					Mnozstvi:        1,
					Cena:            102,
					Sleva:           0,
					VyfakturovanoDo: datum{time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)},
					Frekvence:       frekvence{Years: 0, Months: 1, Days: 0},
				},
			},
		},
		smlouva{
			Zakazka:    "TEST2",
			Nazev:      "TEST2",
			Popis:      "Test2",
			TypDokladu: "TYPE2",
			Platnost: interval{
				Od: time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC),
				Do: time.Date(2016, time.December, 31, 0, 0, 0, 0, time.UTC),
			},
			Polozky: []polozka{
				polozka{
					Nazev:           "Test 2a",
					Cenik:           "CEN2A",
					Mnozstvi:        2,
					Cena:            201,
					Sleva:           0,
					VyfakturovanoDo: datum{time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)},
					Frekvence:       frekvence{Years: 1, Months: 0, Days: 0},
				},
				polozka{
					Nazev:           "Test 2b",
					Cenik:           "CEN2B",
					Mnozstvi:        0,
					Cena:            202,
					Sleva:           22,
					VyfakturovanoDo: datum{time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)},
					Frekvence:       frekvence{Years: 1, Months: 0, Days: 0},
				},
			},
		},
	}

	expected := fakturace{Smlouvy: s}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}

	got, err := NewFakturace(data)
	if err != nil {
		t.Fatalf("error: %s, data = %s", err, data)
	}

	if !reflect.DeepEqual(expected, got) {
		pretty.Fdiff(os.Stdout, expected, got)
		t.Fatalf("Structs are not equal")
	}
}

func TestNewEmpty(t *testing.T) {
	_, err := NewFakturace([]byte("[]"))

	if err != nil {
		t.Fatal(err)
	}
}

func TestMultiple(t *testing.T) {
	f, err := NewFakturace([]byte("[{},{},{}]"))

	if err != nil {
		t.Fatal(err)
	}

	if len(f.Smlouvy) != 3 {
		t.Fatalf("Expected %d elements, got %d", 3, len(f.Smlouvy))
	}
}
