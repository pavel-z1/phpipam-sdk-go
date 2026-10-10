package phpipam

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
)

const testBoolIntStringJSONTrue = `{"foo":"1"}`
const testBoolIntStringJSONFalse = `{"foo":"0"}`
const testBoolIntStringJSONError = `{"foo":"2"}`

type testBoolIntStringType struct {
	Foo BoolIntString `json:"foo"`
}

const testJSONIntStringJSONZeroEmpty = `{"foo":""}`
const testJSONIntStringJSONZeroNumber = `{"foo":"0"}`
const testJSONIntStringJSONNonZeroNumber = `{"foo":"2"}`
const testJSONIntStringJSONError = `{"foo":"a"}`

type testJSONIntStringType struct {
	Foo JSONIntString `json:"foo"`
}

type testJSONIntStringTypeOmitEmpty struct {
	Foo JSONIntString `json:"foo,omitempty"`
}

func setPHPIPAMenv() {
	os.Setenv("PHPIPAM_APP_ID", "foobar")
	os.Setenv("PHPIPAM_ENDPOINT_ADDR", "https://example.com/phpipam/api")
	os.Setenv("PHPIPAM_PASSWORD", "abcdefgh0123456789")
	os.Setenv("PHPIPAM_USER_NAME", "nobody")
}

func unsetPHPIPAMenv() {
	os.Unsetenv("PHPIPAM_APP_ID")
	os.Unsetenv("PHPIPAM_ENDPOINT_ADDR")
	os.Unsetenv("PHPIPAM_PASSWORD")
	os.Unsetenv("PHPIPAM_USER_NAME")
}

func TestPHPIPAMDefaultConfigProviderWithEnv(t *testing.T) {
	setPHPIPAMenv()
	c := DefaultConfigProvider()
	if c.Endpoint != "https://example.com/phpipam/api" {
		t.Fatalf("Expected Endpoint to be https://example.com/phpipam/api, got %s", c.Endpoint)
	}
	if c.Username != "nobody" {
		t.Fatalf("Expected Username to be nobody, got %s", c.Username)
	}
	if c.Password != "abcdefgh0123456789" {
		t.Fatalf("Expected Password to be abcdefgh0123456789, got %s", c.Password)
	}
	if c.AppID != "foobar" {
		t.Fatalf("Expected AppID to be foobar, got %s", c.AppID)
	}
}

func TestPHPIPAMDefaultConfigProviderNoEnv(t *testing.T) {
	unsetPHPIPAMenv()
	c := DefaultConfigProvider()
	if c.Endpoint != "http://localhost/api" {
		t.Fatalf("Expected Endpoint to be http://localhost/api, got %s", c.Endpoint)
	}
	if c.Username != "" {
		t.Fatalf("Expected Username to be empty, got %s", c.Username)
	}
	if c.Password != "" {
		t.Fatalf("Expected Password to be empty, got %s", c.Password)
	}
	if c.AppID != "" {
		t.Fatalf("Expected AppID to be empty, got %s", c.AppID)
	}
}

func TestBoolIntStringUnmarshalJSONTrue(t *testing.T) {
	var actual testBoolIntStringType
	if err := json.Unmarshal([]byte(testBoolIntStringJSONTrue), &actual); err != nil {
		t.Fatalf("Bad: %s", err)
	}
	if actual.Foo != true {
		t.Fatalf("Expected value to be true, got %t", actual)
	}
}

func TestBoolIntStringUnmarshalJSONFalse(t *testing.T) {
	var actual testBoolIntStringType
	if err := json.Unmarshal([]byte(testBoolIntStringJSONFalse), &actual); err != nil {
		t.Fatalf("Bad: %s", err)
	}
	if actual.Foo != false {
		t.Fatalf("Expected value to be false, got %t", actual)
	}
}

func TestBoolIntStringUnmarshalJSONError(t *testing.T) {
	var v testBoolIntStringType
	err := json.Unmarshal([]byte(testBoolIntStringJSONError), &v)
	if err == nil {
		t.Fatalf("Expected error, got none")
	}

	// The error text produced by encoding/json differs between Go versions, so
	// check the error type and value instead.
	var typeErr *json.UnmarshalTypeError
	if !errors.As(err, &typeErr) || typeErr.Value != "bool" {
		t.Fatalf("Expected an UnmarshalTypeError for value bool, got %#v", err)
	}
}

func TestBoolIntStringMarshalJSONTrue(t *testing.T) {
	v := testBoolIntStringType{
		Foo: true,
	}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	expected := testBoolIntStringJSONTrue
	actual := string(b)
	if expected != actual {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}

func TestBoolIntStringMarshalJSONFalse(t *testing.T) {
	v := testBoolIntStringType{
		Foo: false,
	}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	expected := testBoolIntStringJSONFalse
	actual := string(b)
	if expected != actual {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}

func TestJSONIntStringUnmarshalJSONZeroEmpty(t *testing.T) {
	var actual testJSONIntStringType
	if err := json.Unmarshal([]byte(testJSONIntStringJSONZeroEmpty), &actual); err != nil {
		t.Fatalf("Bad: %s", err)
	}
	if actual.Foo != 0 {
		t.Fatalf("Expected value to be 0, got %d", actual)
	}
}

func TestJSONIntStringUnmarshalJSONZeroNumber(t *testing.T) {
	var actual testJSONIntStringType
	if err := json.Unmarshal([]byte(testJSONIntStringJSONZeroNumber), &actual); err != nil {
		t.Fatalf("Bad: %s", err)
	}
	if actual.Foo != 0 {
		t.Fatalf("Expected value to be 0, got %d", actual)
	}
}

func TestJSONIntStringUnmarshalJSONNonZeroNumber(t *testing.T) {
	var actual testJSONIntStringType
	if err := json.Unmarshal([]byte(testJSONIntStringJSONNonZeroNumber), &actual); err != nil {
		t.Fatalf("Bad: %s", err)
	}
	if actual.Foo != 2 {
		t.Fatalf("Expected value to be 2, got %d", actual)
	}
}

func TestJSONIntStringUnmarshalJSONError(t *testing.T) {
	var v testJSONIntStringType
	err := json.Unmarshal([]byte(testJSONIntStringJSONError), &v)
	if err == nil {
		t.Fatalf("Expected error, got none")
	}

	// The error text produced by encoding/json differs between Go versions, so
	// check the error type and value instead.
	var typeErr *json.UnmarshalTypeError
	if !errors.As(err, &typeErr) || typeErr.Value != "int" {
		t.Fatalf("Expected an UnmarshalTypeError for value int, got %#v", err)
	}
}

func TestJSONIntStringMarshalJSONZero(t *testing.T) {
	v := testJSONIntStringType{
		Foo: 0,
	}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	expected := testJSONIntStringJSONZeroNumber
	actual := string(b)
	if expected != actual {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}

func TestJSONIntStringMarshalJSONOmitEmpty(t *testing.T) {
	v := testJSONIntStringTypeOmitEmpty{
		Foo: 0,
	}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	expected := "{}"
	actual := string(b)
	if expected != actual {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}

func TestBoolIntStringUnmarshalJSONNative(t *testing.T) {
	cases := map[string]BoolIntString{
		`{"foo":true}`:  true,
		`{"foo":false}`: false,
		`{"foo":1}`:     true,
		`{"foo":0}`:     false,
		`{"foo":null}`:  false,
	}
	for in, expected := range cases {
		var actual testBoolIntStringType
		if err := json.Unmarshal([]byte(in), &actual); err != nil {
			t.Fatalf("Bad (%s): %s", in, err)
		}
		if actual.Foo != expected {
			t.Fatalf("%s: expected %t, got %t", in, expected, actual.Foo)
		}
	}
}

func TestJSONIntStringUnmarshalJSONNative(t *testing.T) {
	cases := map[string]JSONIntString{
		`{"foo":24}`:   24,
		`{"foo":0}`:    0,
		`{"foo":null}`: 0,
	}
	for in, expected := range cases {
		var actual testJSONIntStringType
		if err := json.Unmarshal([]byte(in), &actual); err != nil {
			t.Fatalf("Bad (%s): %s", in, err)
		}
		if actual.Foo != expected {
			t.Fatalf("%s: expected %d, got %d", in, expected, actual.Foo)
		}
	}
}
