package kevast

import (
	"reflect"
	"testing"
)

func TestNewSession(t *testing.T) {
	tests := []struct {
		name string
		want Session
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSession(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_Run(t *testing.T) {
	type fields struct {
		db *Kevast
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Session{
				db: tt.fields.db,
			}
			if err := s.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Session.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
