package device

import (
	//"crypto/sha1"
	"fmt"
	//"strings"
	"testing"
)

func TestGoogleMapLinkParser(t *testing.T) {
	str1 := "[1111111111]http://maps.google.com/maps?q=32+36.3143,-085+29.1954+()&z=19|"
	str2 := "[1234567890]Hello World|"
	str3 := "[1111111111]http://maps.google.com/maps?q=2+36.3143,-5+29.1954+()&z=19|"
	expected1 := "32.6052,-085.4866"
	expected2 := ""
	expected3 := "2.6052,-5.4866"
	result1 := googleMapLinkParser(str1)
	result2 := googleMapLinkParser(str2)
	result3 := googleMapLinkParser(str3)
	fmt.Println("1: " + result1)
	fmt.Println("2: " + result2)
	fmt.Println("3: " + result3)
	if expected1 != result1 {
		t.Error("Maps string 1 parsed incorrectly")
	}
	if expected2 != result2 {
		t.Error("Other message parsed incorrectly")
	}
	if expected3 != result3 {
		t.Error("Maps string 2 parsed incorrectly")
	}
}
