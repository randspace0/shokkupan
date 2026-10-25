package stats

import "testing"

func TestParse(t *testing.T) {
	out := "myapp\tweb\t3f2a1b9c0d4e\t0.35%\t41.2MiB / 1.94GiB\n" +
		"myapp\tworker\taaaaaaaaaaaa\t150.00%\t512B / 2GB\n"
	got, err := Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d", len(got))
	}
	web := got[0]
	if web.App != "myapp" || web.Process != "web" || web.CPU != 0.35 {
		t.Errorf("web: %+v", web)
	}
	if web.MemBytes != 43201331 || web.MemLimit != 2083059138 {
		t.Errorf("web mem: %+v", web)
	}
	if got[1].CPU != 150 || got[1].MemBytes != 512 || got[1].MemLimit != 2e9 {
		t.Errorf("worker: %+v", got[1])
	}
}

func TestParseEmptyAndBad(t *testing.T) {
	if got, err := Parse(""); err != nil || len(got) != 0 {
		t.Fatalf("empty: %v %v", got, err)
	}
	if _, err := Parse("only\ttwo"); err == nil {
		t.Fatal("expected error")
	}
}
