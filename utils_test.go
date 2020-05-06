package tcb

import "testing"

func TestDecodeApiData(t *testing.T) {
	type args struct {
		apiName string
		data    []byte
		obj     interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "TestDecodeApiData: no error",
			args:    args{
				apiName: "no error",
				data:    []byte(`{"errcode":0,"errmsg":"msg","other":"test"}`),
				obj: &struct {
					ResError
					Other string `json:"other"`
				}{},
			},
			wantErr: false,
		},

		{
			name:    "TestDecodeApiData: with error",
			args:    args{
				apiName: "with error",
				data:    []byte(`{"errcode":1,"errmsg":"msg"}`),
				obj: &struct {
					ResError
				}{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DecodeApiData(tt.args.apiName, tt.args.data, tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("DecodeApiData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}