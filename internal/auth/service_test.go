package auth

import "testing"

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		email   string
		userId  uint
		wantErr bool
	}{
		{"test@gmail.com", 1, false},
		{"testing@test.com", 2, false},
		{"", 0, true},
		{"test@gmail.com", 0, true},
		{"", 4, true},
	}

	for _, tt := range tests {
		_, err := GenerateToken(tt.email, tt.userId)

		if (err != nil) != tt.wantErr {
			t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
		}
	}

}

func TestVerifyToken(t *testing.T) {
	tests := []struct {
		token   string
		wantErr bool
	}{
		{"", true},
		{"test", true},
	}

	for _, tt := range tests {
		_, err := VerifyToken(tt.token)

		if (err != nil) != tt.wantErr {
			t.Errorf("VerifyToken() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}
