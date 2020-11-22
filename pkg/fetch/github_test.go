package fetch

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-version"
)

func testData() []Series {
	return []Series{
		{v("1.20"), v("1.20.1"), v("1.20.2")},
		{v("1.18"), v("1.18.1"), v("1.18.2")},
		{v("1.16"), v("1.16.1"), v("1.16.2")},
		{v("1.14"), v("1.14.1"), v("1.14.2")},
		{v("1.12"), v("1.12.1"), v("1.12.2")},
	}
}

func TestInjectFront(t *testing.T) {
	cur := testData()
	want := []Series{
		{v("1.21"), v("1.21.69"), nil},
		{v("1.20"), v("1.20.1"), v("1.20.2")},
		{v("1.18"), v("1.18.1"), v("1.18.2")},
		{v("1.16"), v("1.16.1"), v("1.16.2")},
		{v("1.14"), v("1.14.1"), v("1.14.2")},
	}
	n := v("1.21.69")

	inject(cur, n, false, 2)

	if !reflect.DeepEqual(cur, want) {
		t.Errorf("\nWant: %v\nGot:  %v", want, cur)
	}
}
func TestInjectMiddle(t *testing.T) {
	cur := testData()
	want := []Series{
		{v("1.20"), v("1.20.1"), v("1.20.2")},
		{v("1.18"), v("1.18.1"), v("1.18.2")},
		{v("1.17"), v("1.17.69"), nil},
		{v("1.16"), v("1.16.1"), v("1.16.2")},
		{v("1.14"), v("1.14.1"), v("1.14.2")},
	}
	n := v("1.17.69")

	inject(cur, n, false, 2)

	if !reflect.DeepEqual(cur, want) {
		t.Errorf("\nWant: %v\nGot:  %v", want, cur)
	}
}
func TestInjectMiddlePre(t *testing.T) {
	cur := testData()
	want := []Series{
		{v("1.20"), v("1.20.1"), v("1.20.2")},
		{v("1.18"), v("1.18.1"), v("1.18.2")},
		{v("1.17"), nil, v("1.17.69")},
		{v("1.16"), v("1.16.1"), v("1.16.2")},
		{v("1.14"), v("1.14.1"), v("1.14.2")},
	}
	n := v("1.17.69")

	inject(cur, n, true, 2)

	if !reflect.DeepEqual(cur, want) {
		t.Errorf("\nWant: %v\nGot:  %v", want, cur)
	}
}
func TestInjectMiddleBoth(t *testing.T) {
	cur := testData()
	want := []Series{
		{v("1.20"), v("1.20.1"), v("1.20.2")},
		{v("1.18"), v("1.18.1"), v("1.18.2")},
		{v("1.17"), v("1.17.69"), v("1.17.70")},
		{v("1.16"), v("1.16.1"), v("1.16.2")},
		{v("1.14"), v("1.14.1"), v("1.14.2")},
	}
	n := v("1.17.69")
	p := v("1.17.70")

	inject(cur, n, false, 2)
	inject(cur, p, true, 2)

	if !reflect.DeepEqual(cur, want) {
		t.Errorf("\nWant: %v\nGot:  %v", want, cur)
	}
}
func TestInjectEnd(t *testing.T) {
	cur := testData()
	want := []Series{
		{v("1.20"), v("1.20.1"), v("1.20.2")},
		{v("1.18"), v("1.18.1"), v("1.18.2")},
		{v("1.16"), v("1.16.1"), v("1.16.2")},
		{v("1.14"), v("1.14.1"), v("1.14.2")},
		{v("1.13"), v("1.13.69"), nil},
	}
	n := v("1.13.69")

	inject(cur, n, false, 2)

	if !reflect.DeepEqual(cur, want) {
		t.Errorf("\nWant: %v\nGot:  %v", want, cur)
	}
}
func TestInjectMiss(t *testing.T) {
	cur := testData()
	want := []Series{
		{v("1.20"), v("1.20.1"), v("1.20.2")},
		{v("1.18"), v("1.18.1"), v("1.18.2")},
		{v("1.16"), v("1.16.1"), v("1.16.2")},
		{v("1.14"), v("1.14.1"), v("1.14.2")},
		{v("1.12"), v("1.12.1"), v("1.12.2")},
	}
	n := v("1.11.69")

	inject(cur, n, false, 2)

	if !reflect.DeepEqual(cur, want) {
		t.Errorf("\nWant: %v\nGot:  %v", want, cur)
	}
}

func TestTrumpHigh(t *testing.T) {
	cur := testData()
	want := []Series{
		{v("1.20"), v("1.20.1"), v("1.20.2")},
		{v("1.18"), v("1.18.69"), v("1.18.2")},
		{v("1.16"), v("1.16.1"), v("1.16.2")},
		{v("1.14"), v("1.14.1"), v("1.14.2")},
		{v("1.12"), v("1.12.1"), v("1.12.2")},
	}
	n := v("1.18.69")

	inject(cur, n, false, 2)

	if !reflect.DeepEqual(cur, want) {
		t.Errorf("\nWant: %v\nGot:  %v", want, cur)
	}
}
func TestTrumpHighPre(t *testing.T) {
	cur := testData()
	want := []Series{
		{v("1.20"), v("1.20.1"), v("1.20.2")},
		{v("1.18"), v("1.18.1"), v("1.18.69")},
		{v("1.16"), v("1.16.1"), v("1.16.2")},
		{v("1.14"), v("1.14.1"), v("1.14.2")},
		{v("1.12"), v("1.12.1"), v("1.12.2")},
	}
	n := v("1.18.69")

	inject(cur, n, true, 2)

	if !reflect.DeepEqual(cur, want) {
		t.Errorf("\nWant: %v\nGot:  %v", want, cur)
	}
}
func TestTrumpEqual(t *testing.T) {
	cur := testData()
	want := []Series{
		{v("1.20"), v("1.20.1"), v("1.20.2")},
		{v("1.18"), v("1.18.1"), v("1.18.2")},
		{v("1.16"), v("1.16.1"), v("1.16.2")},
		{v("1.14"), v("1.14.1"), v("1.14.2")},
		{v("1.12"), v("1.12.1"), v("1.12.2")},
	}
	n := v("1.18.1")

	inject(cur, n, false, 2)

	if !reflect.DeepEqual(cur, want) {
		t.Errorf("\nWant: %v\nGot:  %v", want, cur)
	}
}
func TestTrumpLow(t *testing.T) {
	cur := testData()
	want := []Series{
		{v("1.20"), v("1.20.1"), v("1.20.2")},
		{v("1.18"), v("1.18.1"), v("1.18.2")},
		{v("1.16"), v("1.16.1"), v("1.16.2")},
		{v("1.14"), v("1.14.1"), v("1.14.2")},
		{v("1.12"), v("1.12.1"), v("1.12.2")},
	}
	n := v("1.18.0")

	inject(cur, n, false, 2)

	if !reflect.DeepEqual(cur, want) {
		t.Errorf("\nWant: %v\nGot:  %v", want, cur)
	}
}

func v(s string) *version.Version {
	return version.Must(version.NewVersion(s))
}
