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
	expected1 := "32.6052,-085.4866"
	expected2 := ""
	result1 := googleMapLinkParser(str1)
	result2 := googleMapLinkParser(str2)
	fmt.Println(result1)
	fmt.Println(result2)
	if expected1 != result1 {
		t.Error("Maps string parsed incorrectly")
	}
	if expected2 != result2 {
		t.Error("Other message parsed incorrectly")
	}
}
