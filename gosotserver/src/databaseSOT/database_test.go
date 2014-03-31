package databaseSOT

import (
	"CustomProtocol"
	"crypto/sha1"
	"fmt"
	"strings"
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
	payload := CustomProtocol.ParsePayload(buf)
	testArray := []string{"Param1", "Param2", "Param3"}
	for index, element := range payload {
		if element != testArray[index] {
			t.Error("parsePayload([]byte) did not work as expected.")
		}
	}
	t.Log("parsePayload([]byte) test passed.")
}

func TestIsDeviceStolen(t *testing.T) {

	//fmt.Println(IsDeviceStolen("1"))
	if IsDeviceStolen("1") == true {
		t.Log("IsDeviceStolen(string) test passed.")
	} else {
		t.Error("IsDeviceStolen(string) did not work as expected.")
	}
}

func TestVerifyAccountInfo(t *testing.T) {
	userName := "saw0019@auburn.edu"
	password := "steven"
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

func TestGetLaptopDevices(t *testing.T) {
	//gps device must be present
	userName := "saw0019@auburn.edu"
	jsonList := getLaptopDevices(userName)
	fmt.Println("GetLaptopDevice TEST")
	fmt.Println(string(jsonList))
}

func TestGetGpsDevices(t *testing.T) {
	//gps device must be present
	userName := "saw0019@auburn.edu"
	jsonList := getGpsDevices(userName)
	fmt.Println("GetGpsDevice TEST")
	fmt.Println(string(jsonList))
}

func TestParseTraceRoute(t *testing.T) {
	traceroute := parseTraceRouteString("127.0.0.1:4096~123.1.1.1~123.2.23.2~123.3.3.3")
	fmt.Println("PARSE TRACEROUTE")
	fmt.Println(len(traceroute))
}
