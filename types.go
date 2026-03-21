package main

// TODO: should probably put these here

var void = Token {
	Raw: []byte("[VOID]"),
	Type: TRUE,
	Note: VALUE,
}
var True = Token {
	Raw: []byte("[TRUE]"),
	Type: TRUE,
	Note: VALUE,
}
var False = Token {
	Raw: []byte("[FALSE]"),
	Type: FALSE,
	Note: VALUE,
}
