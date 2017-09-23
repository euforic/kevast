package kevast

import (
	"strconv"
	"testing"
)

func BenchmarkKevast_NewDB(b *testing.B) {
	ts := []struct {
		name string
	}{
		{name: "Basic"},
	}
	for _, tt := range ts {
		b.ResetTimer()
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				NewDB()
			}
		})
	}
}

func BenchmarkKevast_Write(b *testing.B) {
	ts := []struct {
		name string
		db   *Kevast
	}{
		{name: "1", db: dbGen(0, 10)},
		{name: "10", db: dbGen(10, 10)},
		{name: "100", db: dbGen(100, 10)},
	}
	for _, tt := range ts {
		b.ResetTimer()
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if err := tt.db.Write("foo", "bar"); err != nil {
					b.Error(err)
					b.FailNow()
				}
			}
		})
	}
}

func BenchmarkKevast_Read(b *testing.B) {
	ts := []struct {
		name string
		db   *Kevast
	}{
		{name: "1", db: dbGen(0, 10)},
		{name: "10", db: dbGen(10, 10)},
		{name: "100", db: dbGen(100, 10)},
	}
	for _, tt := range ts {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if _, err := tt.db.Read("key0"); err != nil {
					b.Error(err)
					b.FailNow()
				}
			}
		})
	}
}

func BenchmarkKevast_Start(b *testing.B) {
	ts := []struct {
		name string
	}{
		{name: "Basic"},
	}
	for _, tt := range ts {
		db := NewDB()
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if err := db.Start(); err != nil {
					b.Error(err)
					b.FailNow()
				}
			}
		})
	}
}

func BenchmarkKevast_Abort(b *testing.B) {
	ts := []struct {
		name string
		keys int
	}{
		{name: "1x1", keys: 1},
	}
	for _, tt := range ts {
		b.Run(tt.name, func(b *testing.B) {
			db := dbGen(0, tt.keys)
			db.idx = int64(b.N)
			db.stores = append(db.stores, make([]store, b.N)...)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := db.Abort(); err != nil {
					b.Error(err)
					b.FailNow()
				}
			}
		})
	}
}

func BenchmarkKevast_Commit(b *testing.B) {
	ts := []struct {
		name string
		keys int
	}{
		{name: "1x1", keys: 1},
		{name: "1x10", keys: 10},
		{name: "1x100", keys: 100},
	}
	for _, tt := range ts {
		b.Run(tt.name, func(b *testing.B) {
			db := dbGen(b.N, tt.keys)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := db.Commit(); err != nil {
					b.Error(err)
					b.FailNow()
				}
			}
		})
	}
}

func dbGen(txs int, keys int) *Kevast {
	stores := []store{}

	for i := 0; i <= txs; i++ {
		s := store{}
		for j := 0; j < keys; j++ {
			s["key"+strconv.Itoa(j)] = "bar" + strconv.Itoa(i)
		}
		stores = append(stores, s)
	}

	return &Kevast{
		idx:    int64(txs),
		stores: stores,
	}
}
