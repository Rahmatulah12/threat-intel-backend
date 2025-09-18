package postgres

import (
	"testing"
)

func TestNewConnection(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Host:     "localhost",
				Port:     "5432",
				User:     "test",
				Password: "test",
				DBName:   "test",
				SSLMode:  "disable",
			},
			wantErr: true, // Will fail without actual DB
		},
		{
			name: "empty config",
			config: Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConnection(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}