// TODO: write more tests to cover all cases
package kevast

import (
	"reflect"
	"strconv"
	"testing"
)

func TestNewDB(t *testing.T) {
	tests := []struct {
		name string
		want *Kevast
	}{
		{name: "Basic", want: NewDB()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDB(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKevast_Write(t *testing.T) {
	type args struct {
		key string
		val string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Basic", args: args{key: "foo", val: "bar"}, wantErr: false},
		{name: "IntString", args: args{key: "foo", val: "234"}, wantErr: false},
		{name: "MissingKey", args: args{key: "", val: "bar"}, wantErr: true},
		{name: "MissingVal", args: args{key: "foo", val: ""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewDB()
			if err := s.Write(tt.args.key, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("Kevast.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if val, ok := s.stores[0][tt.args.key]; !tt.wantErr && (!ok || val != tt.args.val) {
				t.Errorf("Kevast.Write() error = invalid value")
			}
		})
	}
}

func TestKevast_Read(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "Basic", args: args{key: "foo"}, want: "bar0", wantErr: false},
		{name: "MissingKey", args: args{key: ""}, want: "", wantErr: true},
	}
	s := dbWithData(1)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Read(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kevast.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Kevast.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKevast_Del(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		db      Kevast
		args    args
		wantErr bool
	}{
		{name: "Basic", db: dbWithData(1), args: args{key: "foo"}, wantErr: false},
		{name: "MissingKey", db: dbWithData(1), args: args{key: ""}, wantErr: true},
		{name: "KeyNotExist", db: dbWithData(1), args: args{key: "hello"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.Del(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Kevast.Del() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKevast_Start(t *testing.T) {
	tests := []struct {
		name    string
		wantLen int
		wantErr bool
	}{
		{name: "Basic", wantLen: 2, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := dbWithData(1)
			if err := s.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Kevast.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
			if l := len(s.stores); l != tt.wantLen {
				t.Errorf("Kevast.Start() Failed to create store for Tx wantLen:%d, got:%d ", tt.wantLen, l)
			}
		})
	}
}

func TestKevast_Abort(t *testing.T) {
	tests := []struct {
		name    string
		depth   int
		wantLen int
		wantErr bool
	}{
		{name: "Basic", depth: 2, wantLen: 1, wantErr: false},
		{name: "Nested", depth: 3, wantLen: 2, wantErr: false},
		{name: "NotInTx", depth: 1, wantLen: 1, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := dbWithData(tt.depth)
			if err := s.Abort(); (err != nil) != tt.wantErr {
				t.Errorf("Kevast.Abort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotLen := len(s.stores); gotLen != tt.wantLen {
				t.Errorf("Kevast.Abort() = %v, want %v", gotLen, tt.wantLen)
			}
		})
	}
}

func TestKevast_Commit(t *testing.T) {
	tests := []struct {
		name    string
		wantLen int
		want    Kevast
		wantErr bool
	}{
		{
			name:    "Basic",
			wantLen: 1,
			want:    Kevast{idx: 0, stores: []store{store{"foo": "bar1", "baz": "qux1"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := dbWithData(2)
			if err := s.Commit(); (err != nil) != tt.wantErr {
				t.Errorf("Kevast.Commit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLen := len(s.stores); gotLen != tt.wantLen {
				t.Errorf("Kevast.Commit() Tx store not cleared wantLen: %v gotLen: %v", gotLen, tt.wantLen)
			}
			if !compareKevasts(tt.want, s) {
				t.Errorf("Kevast.Commit() faild to commit Tx want: %v got: %v", tt.want, s)
			}
		})
	}
}

func compareKevasts(want Kevast, got Kevast) bool {
	// Check if indexes are equal
	if want.idx != got.idx {
		return false
	}
	// Check if stores len is equal
	if len(want.stores) != len(got.stores) {
		return false
	}

	for i, gotStore := range got.stores {
		// Check if store's len are equal
		if len(want.stores[i]) != len(gotStore) {
			return false
		}

		// Check if store's values are equal for each key
		for k, v := range gotStore {
			if want.stores[i][k] != v {
				return false
			}
		}
	}
	return true
}

func dbWithData(layers int) Kevast {
	stores := []store{}

	for i := 0; i < layers; i++ {
		stores = append(stores, store{
			"foo": "bar" + strconv.Itoa(i),
			"baz": "qux" + strconv.Itoa(i),
		})
	}

	return Kevast{
		idx:    int64(layers) - 1,
		stores: stores,
	}
}
