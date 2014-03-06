package databaseSOT

import (
	//"crypto/sha1"
	//"fmt"
	//"strings"
	"testing"
)

func TestParsePayload(t *testing.T) {
	a := []byte("Param1")
	b := []byte("Param2")
	c := []byte("Param3")
	buf := []byte{}
	buf = append(buf, a...)
	buf = append(buf, 0x1B)
	buf = append(buf, b...)
	buf = append(buf, 0x1B)
	buf = append(buf, c...)
	buf = append(buf, 0x1B)
	payload := parsePayload(buf)
	testArray := []string{"Param1", "Param2", "Param3"}
	for index, element := range payload {
		if element != testArray[index] {
			t.Error("parsePayload([]byte) did not work as expected.")
		}
	}
	t.Log("parsePayload([]byte) test passed.")
}

func TestIsDeviceStolen(t *testing.T) {
	if IsDeviceStolen("3344449464") {
		t.Log("IsDeviceStolen(string) test passed.")
	} else {
		t.Error("IsDeviceStolen(string) did not work as expected.")
	}
}

func TestVerifyAccountInfo(t *testing.T) {
	userName := "Test@Test.com"
	password := "test"
	h := sha1.New()
	h.Write([]byte(strings.Join([]string{userName, password}, "")))
	hashedPass := fmt.Sprintf("%x", h.Sum(nil))
	accountValid, passwordValid := VerifyAccountInfo(userName, string(hashedPass))
	if accountValid && passwordValid {
		t.Log("VerifyAccount(string, string) test passed.")
	} else {
		t.Error("VerifyAccount(string, string) did not work as expected.")
	}

}

func TestGetUserDevices(t *testing.T) {
	devices := GetUserDevices("Test@Test.com")
	testArray := []string{"Stolen Laptop Test", "Stolen GPS Device Test"}
	for index, element := range devices {
		if element != testArray[index] {
			t.Error("GetUserDevices(string) did not work as expected.")
		}
	}
	t.Log("GetUserDevices(string) test passed.")
}
