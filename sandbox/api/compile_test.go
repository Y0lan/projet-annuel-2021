package main

import (
	"os"
	"reflect"
	"testing"
)

func Test_compilePython(t *testing.T) {
	type args struct {
		code string
	}

	tests := []struct {
		name        string
		args        args
		wantSuccess string
	}{
		{
			"can compile python code",
			args{code: "def hello_world():\n    print(\"Hello, World!\")\ndef sum(a, b):\n    return a + b\nif __name__ == \"__main__\":\n    hello_world()\n    sum(1,2)\n"},
			"success",
		},
		{
			"can not compile incorrect python code",
			args{code: "def hello_world()\n    print(\"Hello, World!\")\ndef sum(a, b):\n    return a + b\nif __name__ == \"__main__\":\n    hello_world()\n    sum(1,2)\n"},
			"failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotSuccess := compilePython(tt.args.code)
			if gotSuccess != tt.wantSuccess {
				t.Errorf("compilePython() gotSuccess = %v, want %v", gotSuccess, tt.wantSuccess)
			}
		})
	}
}

func Test_compileGo(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name        string
		args        args
		wantSuccess string
	}{
		{
			"can compile correct go code",
			args{code: "package main\nimport \"fmt\"\nfunc Hello() string {\n\treturn \"Hello, world\"\n}\nfunc main() {\n\tfmt.Println(Hello())\n}\n"},
			"success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotSuccess := compileGo(tt.args.code)
			if gotSuccess != tt.wantSuccess {
				t.Errorf("compileGo() gotSuccess = %v, want %v", gotSuccess, tt.wantSuccess)
			}
		})
	}
}

func Test_compileRust(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name        string
		args        args
		wantOutput  string
		wantSuccess string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, gotSuccess := compileRust(tt.args.code)
			if gotOutput != tt.wantOutput {
				t.Errorf("compileRust() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
			if gotSuccess != tt.wantSuccess {
				t.Errorf("compileRust() gotSuccess = %v, want %v", gotSuccess, tt.wantSuccess)
			}
		})
	}
}

func Test_createFile(t *testing.T) {
	type args struct {
		code      string
		extension string
		fileName  string
		prefix    string
		suffix    string
	}
	tests := []struct {
		name string
		args args
		want os.File
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createFile(tt.args.code, tt.args.extension, tt.args.fileName, tt.args.prefix, tt.args.suffix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isLanguageSupported(t *testing.T) {
	type args struct {
		givenLanguage string
	}
	tests := []struct {
		name            string
		args            args
		wantIsSupported bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsSupported := isLanguageSupported(tt.args.givenLanguage); gotIsSupported != tt.wantIsSupported {
				t.Errorf("isLanguageSupported() = %v, want %v", gotIsSupported, tt.wantIsSupported)
			}
		})
	}
}

func Test_testCode(t *testing.T) {
	type args struct {
		code     string
		test     string
		language string
	}
	tests := []struct {
		name        string
		args        args
		wantOutput  string
		wantSuccess string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, gotSuccess := testCode(tt.args.code, tt.args.test, tt.args.language)
			if gotOutput != tt.wantOutput {
				t.Errorf("testCode() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
			if gotSuccess != tt.wantSuccess {
				t.Errorf("testCode() gotSuccess = %v, want %v", gotSuccess, tt.wantSuccess)
			}
		})
	}
}

func Test_testGo(t *testing.T) {
	type args struct {
		code string
		test string
	}
	tests := []struct {
		name        string
		args        args
		wantOutput  string
		wantSuccess string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, gotSuccess := testGo(tt.args.code, tt.args.test)
			if gotOutput != tt.wantOutput {
				t.Errorf("testGo() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
			if gotSuccess != tt.wantSuccess {
				t.Errorf("testGo() gotSuccess = %v, want %v", gotSuccess, tt.wantSuccess)
			}
		})
	}
}

func Test_testPython(t *testing.T) {
	type args struct {
		code string
		test string
	}
	tests := []struct {
		name        string
		args        args
		wantOutput  string
		wantSuccess string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, gotSuccess := testPython(tt.args.code, tt.args.test)
			if gotOutput != tt.wantOutput {
				t.Errorf("testPython() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
			if gotSuccess != tt.wantSuccess {
				t.Errorf("testPython() gotSuccess = %v, want %v", gotSuccess, tt.wantSuccess)
			}
		})
	}
}

func Test_testRust(t *testing.T) {
	type args struct {
		code string
		test string
	}
	tests := []struct {
		name        string
		args        args
		wantOutput  string
		wantSuccess string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, gotSuccess := testRust(tt.args.code, tt.args.test)
			if gotOutput != tt.wantOutput {
				t.Errorf("testRust() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
			if gotSuccess != tt.wantSuccess {
				t.Errorf("testRust() gotSuccess = %v, want %v", gotSuccess, tt.wantSuccess)
			}
		})
	}
}
