package tests

import (
	"myapp/tools"
	"testing"
)

func TestSuccessPasswordCompareAndPasswordHash(t *testing.T) {
	type args struct {
		hashed   string
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{}
	tests = append(tests, struct {
		name string
		args args
		want bool
	}{
		name: "Test Case 1",
		args: args{hashed: tools.PasswordHash("12345"), password: "12345"},
		want: true,
	})
	tests = append(tests, struct {
		name string
		args args
		want bool
	}{
		name: "Test Case 2",
		args: args{hashed: tools.PasswordHash("kmdflgfiox2@)(*%asf"), password: "kmdflgfiox2@)(*%asf"},
		want: true,
	})
	tests = append(tests, struct {
		name string
		args args
		want bool
	}{
		name: "Test Case 3",
		args: args{hashed: tools.PasswordHash(""), password: ""},
		want: true,
	})
	tests = append(tests, struct {
		name string
		args args
		want bool
	}{
		name: "Test Case 4",
		args: args{hashed: tools.PasswordHash("thisisatestpassword"), password: "thisisatestpassword"},
		want: true,
	})
	tests = append(tests, struct {
		name string
		args args
		want bool
	}{
		name: "Test Case 5",
		args: args{hashed: tools.PasswordHash("     "), password: "     "},
		want: true,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tools.PasswordCompare(tt.args.hashed, tt.args.password); got != tt.want {
				t.Errorf("PasswordCompare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorPasswordCompareAndPasswordHash(t *testing.T) {
	type args struct {
		hashed   string
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{}
	tests = append(tests, struct {
		name string
		args args
		want bool
	}{
		name: "Test Case 1",
		args: args{hashed: tools.PasswordHash("abcde"), password: "12345"},
		want: false,
	})
	tests = append(tests, struct {
		name string
		args args
		want bool
	}{
		name: "Test Case 2",
		args: args{hashed: tools.PasswordHash("kmdflgfioasgx2@)(*%asf"), password: "kmdflgfiox2@)(*%asf"},
		want: false,
	})
	tests = append(tests, struct {
		name string
		args args
		want bool
	}{
		name: "Test Case 3",
		args: args{hashed: tools.PasswordHash(" dsgdsgds"), password: ""},
		want: false,
	})
	tests = append(tests, struct {
		name string
		args args
		want bool
	}{
		name: "Test Case 4",
		args: args{hashed: tools.PasswordHash("    "), password: "     "},
		want: false,
	})
	tests = append(tests, struct {
		name string
		args args
		want bool
	}{
		name: "Test Case 5",
		args: args{hashed: tools.PasswordHash("asdasg@*$("), password: "gfhjgf@*$("},
		want: false,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tools.PasswordCompare(tt.args.hashed, tt.args.password); got != tt.want {
				t.Errorf("PasswordCompare() = %v, want %v", got, tt.want)
			}
		})
	}
}
