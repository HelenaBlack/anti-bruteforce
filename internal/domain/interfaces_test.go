package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInSubnet(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		subnet   string
		expected bool
		wantErr  bool
	}{
		{
			name:     "IP in subnet",
			ip:       "192.168.1.5",
			subnet:   "192.168.1.0/24",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "IP not in subnet",
			ip:       "192.168.2.1",
			subnet:   "192.168.1.0/24",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "IPv6 IP in subnet",
			ip:       "2001:db8::1",
			subnet:   "2001:db8::/32",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "Invalid subnet",
			ip:       "192.168.1.1",
			subnet:   "invalid",
			expected: false,
			wantErr:  true,
		},
		{
			name:     "Invalid IP",
			ip:       "invalid",
			subnet:   "192.168.1.0/24",
			expected: false,
			wantErr:  false, // net.ParseIP returns nil, subnet.Contains(nil) returns false
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsInSubnet(tt.ip, tt.subnet)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
