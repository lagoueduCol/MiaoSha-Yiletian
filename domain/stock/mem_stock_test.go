package stock

import (
	"testing"
)

func TestMemStock(t *testing.T) {
	var (
		st  Stock
		err error
		val int64
	)
	if st, err = NewMemStock("101", "1001"); err != nil {
		t.Fatal(err)
	}
	if err = st.Set(10, 1); err != nil {
		t.Fatal(err)
	}
	if val, err = st.Get(); err != nil {
		t.Fatal(err)
	} else if val != 10 {
		t.Fatal("not equal 10")
	}
	if val, err = st.Sub(); err != nil {
		t.Fatal(err)
	} else if val != 9 {
		t.Fatal("not equal 9")
	}
	if err = st.Del(); err != nil {
		t.Fatal(err)
	}
	if val, err = st.Get(); err != nil {
		t.Fatal(err)
	} else if val != 0 {
		t.Fatal("not equal 0")
	}
}
