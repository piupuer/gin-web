package response

import (
	"testing"
)

func TestPageInfo_getLimit(t *testing.T) {
	type fields struct {
		Size uint
		Num  uint
	}
	tests := []struct {
		name       string
		fields     fields
		wantLimit  int
		wantOffset int
	}{
		{
			name: "case1",
			fields: fields{
				Size: 0,
				Num:  0,
			},
			wantLimit:  10,
			wantOffset: 0,
		},
		{
			name: "case2",
			fields: fields{
				Size: 10,
				Num:  5,
			},
			wantLimit:  10,
			wantOffset: 40,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PageInfo{
				PageSize: tt.fields.Size,
				PageNum:  tt.fields.Num,
			}
			gotLimit, gotOffset := s.GetLimit()
			if gotLimit != tt.wantLimit {
				t.Errorf("GetLimit() gotLimit = %v, want %v", gotLimit, tt.wantLimit)
			}
			if gotOffset != tt.wantOffset {
				t.Errorf("GetLimit() gotOffset = %v, want %v", gotOffset, tt.wantOffset)
			}
		})
	}
}
