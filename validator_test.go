package govalidate

import (
	"testing"
)

func TestBodyRequest(t *testing.T) {
	type Embedded struct {
		A string
	}
	type RequestNumberMin struct {
		A int `validate:"min=3"`
	}
	type RequestNumberMax struct {
		A int `validate:"max=6"`
	}
	type RequestString struct {
		A string `validate:"required"`
		B string
	}
	type RequestPassword struct {
		A string `validate:"password"`
		B string
	}
	type RequestObject struct {
		A Embedded `validate:"required"`
	}

	type args struct {
		d interface{}
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"TestBodyRequestStringRequiredSuccess", args{d: RequestString{A: "test", B: ""}}, false},
		{"TestBodyRequestStringRequiredFailingWithEmptyString", args{d: RequestString{A: "", B: ""}}, true},
		{"TestBodyRequestStringRequiredFailingWithNoField", args{d: RequestString{B: ""}}, true},

		{"TestBodyRequestPasswordSuccess", args{d: RequestPassword{A: "P@ssw0rdStrong", B: ""}}, false},
		{"TestBodyRequestPasswordLowerThan8", args{d: RequestPassword{A: "week", B: ""}}, true},
		{"TestBodyRequestPasswordWithoutAtLeastOneUppercase", args{d: RequestPassword{A: "nouppercase!1", B: ""}}, true},
		{"TestBodyRequestPasswordWithoutAtLeastOneLowercase", args{d: RequestPassword{A: "NOLOWERCASE!1", B: ""}}, true},
		{"TestBodyRequestPasswordWithoutAtLeastOneSpecialChar", args{d: RequestPassword{A: "Nospecial1", B: ""}}, true},
		{"TestBodyRequestPasswordWithoutAtLeastOneDigit", args{d: RequestPassword{A: "Nonumber@", B: ""}}, true},
		{"TestBodyRequestPasswordWithMoreThan20", args{d: RequestPassword{A: "P@ssw0rdtolongtofitonthevalidation", B: ""}}, true},

		{"TestBodyRequestObjectRequiredSuccess", args{d: RequestObject{Embedded{}}}, false},

		{"TestBodyRequestIntegerMinSuccess", args{d: RequestNumberMin{5}}, false},
		{"TestBodyRequestIntegerMaxSuccess", args{d: RequestNumberMax{5}}, false},
		{"TestBodyRequestIntegerMinFailingRange", args{d: RequestNumberMin{2}}, true},
		{"TestBodyRequestIntegerMaxFailingRange", args{d: RequestNumberMax{7}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BodyRequest(tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("BodyRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
