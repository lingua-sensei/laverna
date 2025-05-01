package synthesize

import "testing"

func TestNewSpeed(t *testing.T) {
	tests := []struct {
		str  string
		want Speed
	}{
		{
			str:  "normal",
			want: NormalSpeed,
		},
		{
			str:  "slower",
			want: SlowerSpeed,
		},
		{
			str:  "slowest",
			want: SlowestSpeed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			if got := NewSpeed(tt.str); got != tt.want {
				t.Errorf("NewSpeed(%v): got = %v, want = %v", tt.str, got, tt.want)
			}
		})
	}
}
