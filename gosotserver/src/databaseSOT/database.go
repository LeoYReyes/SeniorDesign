/*
* Steven Whaley - created: February 18, 2014  - last updated: February 27, 2014
*
* OVERVIEW:
*
*     This is the database portion of the server. It provides functionality for interacting with the database. Actions
*     such as user sign up and login must use functions provided in this file to update and query the database.
*     See the comments for each method for details.
*
*   useful links:
*       mysql driver for google go:
*         https://github.com/ziutek/mymysql
*     golang documentation:
*         http://golang.org/doc/
*
*   TO DO: Interface with other server components using requests, continue testing existing functions.
*
*        -also what other functionality is needed?
*/

package databaseSOT

import (
	"CustomProtocol"
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	"os"
	"strconv"
	"strings"
	//_ "github.com/ziutek/mymysql/native" // Native engine
	_ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
)

/*
* These structs represent tables in the database.
 */
type Account struct {
	CustomerId int64
	Id         int64
	UserName   string
	Password   string
	Admin      bool
}
type Coordinates struct {
	DeviceId  int64
	Latitude  float64
	Longitude float64
	Timestamp string
	Id        int64
}
type Customer struct {
	Id          int64
	PhoneNumber string
	Address     string
	Email       string
	FirstName   string
	LastName    string
}
type GpsDevice struct {
	Id         int64
	Name       string
	CustomerId int64
	IsStolen   int64
}
type IpAddress struct {
	ListId    int64
	IpAddress string
	Id        int64
}
type IpList struct {
	Id        int64
	DeviceId  int64
	Timestamp string
}
type KeyLogs struct {
	DeviceId  int64
	Timestamp string
	Data      string
	id        int64
}
type LaptopDevice struct {
	Id         int64
	DeviceName string
	CustomerId int64
	MacAddress string
	IsStolen   int64
}

var toServer chan *CustomProtocol.Request
var fromServer chan *CustomProtocol.Request

func StartDatabaseServer(toServerIn chan *CustomProtocol.Request, fromServerIn chan *CustomProtocol.Request) {
	toServer = toServerIn
	fromServer = fromServerIn
	go chanHandler()
}

func chanHandler() {
	for {
		select {
		case req := <-fromServer:
			go processRequest(req)
		}
	}
}

func processRequest(req *CustomProtocol.Request) {
	payload := CustomProtocol.ParsePayload(req.Payload)
	/*for index, element := range payload {
		fmt.Println("Payload element", index, ": ", element)
	}*/
	switch req.OpCode {
	case CustomProtocol.NewAccount:
		SignUp(payload[0], payload[1], payload[2], payload[3], payload[4])
		res := make([]byte, 2)
		res[0] = 1
		req.Response <- res
	case CustomProtocol.NewDevice:
		registerNewDevice(payload[0], payload[1], payload[2], payload[3])
		res := make([]byte, 2)
		res[0] = 1
		req.Response <- res
	case CustomProtocol.UpdateDeviceGPS:
		updated := updateDeviceGps(payload[0], payload[1], payload[2])
		res := make([]byte, 2)
		if updated == true {
			res[0] = 1
		} else {
			res[0] = 0
		}
		req.Response <- res
	case CustomProtocol.VerifyLoginCredentials:
		/*str := []string{}
		pos := 1
		for index, element := range req.Payload {
			if element == 0x1B {
				str = append(str, string(req.Payload[pos:index-1]))
				pos = index + 2
			}
		}
		fmt.Println(str)*/
		accountValid, passwordValid := VerifyAccountInfo(payload[0], payload[1])
		res := make([]byte, 2)
		if accountValid {
			res[0] = 1
			if passwordValid {
				res[1] = 1
			} else {
				res[0] = 0
			}
		} else {
			res[0] = 0
			res[1] = 0
		}
		//fmt.Println("Method Return: ", accountValid, passwordValid)
		//fmt.Println("Byte form: ", res)
		req.Response <- res
	case CustomProtocol.SetAccount:
		accSet := updateAccountInfo(payload[0],payload[1], payload[2])
		res := make([]byte, 1)
		if accSet == true {
			res[0] = 1
		} else {
			res[0] = 0
		}
		req.Response <- res
	case CustomProtocol.GetDevice:
	case CustomProtocol.SetDevice:
	case CustomProtocol.GetDeviceList:
	case CustomProtocol.CheckDeviceStolen: //
		isStolen := IsDeviceStolen(payload[0])
		res := make([]byte, 1)
		if isStolen == true {
			res[0] = 1
		} else {
			res[0] = 0
		}
		req.Response <- res
	case CustomProtocol.UpdateUserKeylogData: //
		boolResult := UpdateKeylog(payload[0], payload[1])
		res := make([]byte, 1)
		if boolResult == true {
			res[0] = 1
		} else {
			res[0] = 0
		}
		req.Response <- res
	case CustomProtocol.UpdateUserIPTraceData: //
		boolResult := UpdateTraceRoute(payload[0], payload[1])
		res := make([]byte, 1)
		if boolResult == true {
			res[0] = 1
		} else {
			res[0] = 0
		}
		req.Response <- res
	default:
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func checkedResult(rows []mysql.Row, res mysql.Result, err error) ([]mysql.Row, mysql.Result) {
	checkError(err)
	return rows, res
}

/*
* Used to form connection with the database.
*
*
* Steven Whaley Feb, 27 - created
 */

func connect() (connection mysql.Conn) {
	user := "root"
	pass := "toor"
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"

	db := mysql.New(proto, "", addr, user, pass, dbname)

	err := db.Connect()
	if err != nil {
		panic(err)
	}

	return db
}

/*
* Used to close connection with the database.
*
*
* Steven Whaley Feb, 27 - created
 */

func disconnect(connection mysql.Conn) {

	checkError(connection.Close())

}

/*
* Takes in user information entered on the sign up page, and creates account and customer entries
* in the database.
*
*
* Steven Whaley Feb, 26 - created
 */

func SignUp(firstname string, lastname string, email string, phoneNumber string, password string) {

	db := connect()

	db.Query("INSERT INTO customer (firstName, lastName, email, phoneNumber) VALUES ('" + firstname + "', '" + lastname + "', '" + email + "', '" + phoneNumber + "')")

	db.Query("INSERT INTO account (userName, password, customerId) SELECT '" + email + "', '" + password + "', id FROM customer WHERE email='" + email + "'")

	disconnect(db)

	return
}

/*
*
*
*
* Steven Whaley Mar, 18 - created
 */

func registerNewDevice(deviceType string, deviceId string, deviceName string, userId string) {
	db := connect()

	if deviceType != "gps" && deviceType != "laptop" {
		print("invalid device type")
	} else {
		if deviceType == "gps" {
			db.Query("INSERT INTO gpsDevice (deviceName, deviceId, customerId) SELECT '" + deviceName + "', '" + deviceId + "', id FROM customer WHERE email='" + userId + "'")
		} else if deviceType == "laptop" {
			fmt.Println("INSERT INTO laptopDevice (deviceName, deviceId, customerId) SELECT '" + deviceName + "', '" + deviceId + "', id FROM customer WHERE email='" + userId + "'")
			db.Query("INSERT INTO laptopDevice (deviceName, deviceId, customerId) SELECT '" + deviceName + "', '" + deviceId + "', id FROM customer WHERE email='" + userId + "'")
		}
	}

	disconnect(db)
}

/*
*  IsDeviceStolen(deviceId string) (bool) takes in device id and return a
*  boolean indicating whether that device is marked stolen.
*
*
* Steven Whaley Mar, 1 - created
 */

func IsDeviceStolen(deviceId string) bool {

	bool1 := false

	db := connect()
	if len(deviceId) != 12 {
		rows, res, err := db.Query("select isStolen from gpsDevice where phoneNumber = '" + deviceId + "'")
		if err != nil {
			panic(err)
		}

		res = res

		for _, row := range rows {
			for _, col := range row {
				if col == nil {
					// col has NULL value
				} else {
					// Do something with text in col (type []byte)
				}
			}

			var err2 error

			val1 := row[0].([]byte)

			temp, err2 := strconv.ParseInt(string(val1[:]), 10, 64)
			err2 = err2

			if temp == 1 {
				bool1 = true
			} else {
				bool1 = false
			}
		}
	} else {

		rows2, res2, err3 := db.Query("select isStolen from laptopDevice where macAddress = '" + deviceId + "'")
		if err3 != nil {
			panic(err3)
		}

		res2 = res2

		for _, row := range rows2 {
			for _, col := range row {
				if col == nil {
					// col has NULL value
				} else {
					// Do something with text in col (type []byte)
				}
			}

			val2 := row[0].([]byte)

			temp2, err4 := strconv.ParseInt(string(val2[:]), 10, 64)
			err4 = err4

			if temp2 == 1 {
				bool1 = true
			} else {
				bool1 = false
			}
		}
	}

	disconnect(db)

	return bool1
}

/*
*
*
*
* Steven Whaley Mar, 23
*/

func updateAccountInfo(oldUsername string, newUsername string, newPassword string) bool{

  bool1 := true

  db := connect()
 
    rows, res, err := db.Query("UPDATE account SET userName = '" + newUsername + "', password = '" + newPassword + "' WHERE userName = '" + oldUsername + "'")
    rows = rows
    res = res

    if err != nil {
    panic(err)
    }

  disconnect(db)

  return bool1
}

/*
*
*
*
* Steven Whaley Mar, 1 - created
* Mar 17 - seems to be working
 */

func UpdateKeylog(deviceId string, keylog string) bool {

	bool1 := true

	db := connect()

	rows, res, err := db.Query("UPDATE keylogs SET data = concat(data, '" + keylog + "') WHERE deviceId = '" + deviceId + "'")
	rows = rows
	res = res

	if err != nil {
		panic(err)
	}

	disconnect(db)

	return bool1
}

/*
*
*
*
* Steven Whaley Mar, 23 - created
*
*/

func updateDeviceGps(deviceId string, latitude string, longitude string) bool {
    bool1 := true

    db := connect()
 
    rows, res, err := db.Query("UPDATE gpsDevice SET latitude = '" +  latitude + "', longitude = '" + longitude + "' WHERE id = '" + deviceId + "'")
    rows = rows
    res = res

    if err != nil {
    panic(err)
    }

    disconnect(db)

    return bool1
}

/*
*
*
*
* Steven Whaley Mar, 17 - created
*
 */

func UpdateTraceRoute(deviceId string, traceRoute string) bool {

	bool1 := true
	max := int64(-1)

	db := connect()

	db.Query("INSERT INTO ipList (deviceId) VALUES ('" + deviceId + "')")

	rows, res, err := db.Query("SELECT MAX(id) FROM ipList")

	//rows, res, err := db.Query("INSERT INTO ipAddress (listId,ipAddress) VALUES ('" + deviceId + "', '" + ip + "')")

	for _, row := range rows {
		for _, col := range row {
			if col == nil {
				// col has NULL value
			} else {
				// Do something with text in col (type []byte)
			}
		}
		err = err

		val1 := row[0].([]byte)

		temp, err2 := strconv.ParseInt(string(val1[:]), 10, 64)

		max = temp

		err2 = err2
		rows = rows
		res = res
	}

	maxs := strconv.FormatInt(max, 10)

	print(max, "\n")
	print(maxs + "\n")

	var list []string

	list = parseTraceRouteString(traceRoute)

	for i := 0; i < len(list); i++ {

		db.Query("INSERT INTO ipAddress (listId,ipAddress) VALUES ('" + maxs + "', '" + list[i] + "')")
	}

	disconnect(db)

	return bool1
}

func parseTraceRouteString(traceRoute string) (arr []string) {

	var list []string

	trace := "127.0.01231.1~123.1.1.1~123.2.23.2~123.3.3.3"
	print(trace + "\n")
	num := strings.Count("127.0.0.1~123.1.1.1~123.2.2.2~123.3.3.3", "~") + 1

	address1 := ""

	for i := 0; i < num; i++ {

		if i != num-1 {
			address1 = trace[0:strings.Index(trace, "~")]
			list = append(list, address1)
			fmt.Println("extracted: " + list[i])

			trace = trace[strings.Index(trace, "~")+1 : len(trace)]
		} else {
			address1 = trace
			list = append(list, address1)
			fmt.Println("extracted: " + list[i])
		}
	}
	return list
}

/*
*
*
* Steven Whaley Mar, 1 - created
 */

func IsGpsDevice(deviceId string) bool {

	bool1 := false

	db := connect()

	rows, res, err := db.Query("select * from gpsDevice where id = '" + deviceId + "'")
	if err != nil {
		panic(err)
	}

	res = res

	for _, row := range rows {
		for _, col := range row {
			if col == nil {
				// col has NULL value
			} else {
				// Do something with text in col (type []byte)
			}
		}

		var err2 error
		temp := int64(-1)

		val1 := row[0].([]byte)

		temp, err2 = strconv.ParseInt(string(val1[:]), 10, 64)
		err2 = err2

		if temp != -1 {
			bool1 = true
		} else {
			bool1 = false
		}
	}

	disconnect(db)

	return bool1
}

/*
* Takes in user information entered on the sign up page, and creates account and customer entries
* in the database.
*
*
* Steven Whaley Feb, 26 - created
 */

func VerifyAccountInfo(username string, password string) (bool, bool) {

	bool1 := false
	bool2 := false
	fmt.Println(username)
	fmt.Println(password)
	accountInfo := new(Account)

	db := connect()

	rows, res, err := db.Query("select * from account")
	if err != nil {
		panic(err)
	}

	res = res

	for _, row := range rows {
		for _, col := range row {
			if col == nil {
				// col has NULL value
			} else {
				// Do something with text in col (type []byte)
			}
		}

		val4 := row[2].([]byte)
		val5 := row[3].([]byte)

		accountInfo.UserName = string(val4[:])
		accountInfo.Password = string(val5[:])

		if accountInfo.UserName == username {
			bool1 = true
		}
		if accountInfo.Password == password {
			bool2 = true
		}
	}

	disconnect(db)

	return bool1, bool2
}