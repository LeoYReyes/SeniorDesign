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

// TODO: handle apostrophe

package databaseSOT

import (
	"CustomProtocol"
	"device"
	"encoding/json"
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
	case CustomProtocol.ActivateGPS:
		flagStolen("gps", payload[0])
		res := make([]byte, 2)
		res[0] = 1
		req.Response <- res
	case CustomProtocol.FlagStolen:
		flagStolen("laptop", payload[0])
		res := make([]byte, 2)
		res[0] = 1
		req.Response <- res
	case CustomProtocol.FlagNotStolen:
		flagNotStolen("laptop", payload[0])
		res := make([]byte, 2)
		res[0] = 1
		req.Response <- res
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
		accSet := updateAccountInfo(payload[0], payload[1], payload[2])
		res := make([]byte, 1)
		if accSet == true {
			res[0] = 1
		} else {
			res[0] = 0
		}
		req.Response <- res
	case CustomProtocol.GetDevice:
		res := make([]byte, 5)

		if payload[0] == "gps" {
			res = getGpsDevices(payload[1])
		} else if payload[0] == "laptop" {
			res = getLaptopDevices(payload[1])
		} else {
			fmt.Println("CustomProtocol.GetDevice payload[0] must be either gps or laptop")
		}
		req.Response <- res
	case CustomProtocol.SetDevice:
	case CustomProtocol.GetDeviceList:
		res := []byte{}
		res = append(res, getLaptopDevices(payload[0])...)
		res = append(res, getGpsDevices(payload[0])...)
		req.Response <- res
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

func getGpsDevices(email string) []byte {

	var list []device.GPSDevice

	db := connect()

	//finding customerId to be used for selecting devices
	rows, _, err := db.Query("select customerId from account where userName = '" + email + "'")
	if err != nil {
		panic(err)
	}

	customerId := string(rows[0][0].([]byte))

	//adding gpsDevices to the the devices list
	gpsRows, _, gpsErr := db.Query("select * from gpsDevice where customerId = '" + customerId + "'")
	if gpsErr != nil {
		panic(gpsErr)
	}

	for _, gps := range gpsRows {

		//gpsId := string(gps[0].([]byte))
		gpsName := string(gps[1].([]byte))
		deviceId := string(gps[3].([]byte))
		isStolen := gps[4].([]byte)

		// Create GpsDevice struct and append to list of devices
		list = append(list, device.GPSDevice{device.Device{deviceId, gpsName, isStolen[0]}})
	}

	disconnect(db)
	deviceListJson, _ := json.Marshal(list)

	return deviceListJson
}

/*
*
*
*
* Steven Whaley Mar, 18 - created
 */

func registerNewDevice(deviceType string, deviceId string, deviceName string, userId string) {
	db := connect()
	//print(deviceType)

	if deviceType != "gps" && deviceType != "laptop" {
		print("invalid device type")
	} else {
		// query does not work on VARCHAR with length of 50
		if deviceType == "gps" {
			fmt.Println("Writing to gpsDevice...")
			db.Query("INSERT INTO gpsDevice (deviceName, deviceId, customerId) SELECT '" + deviceName + "', '" + deviceId + "', id FROM customer WHERE email='" + userId + "'")
		} else if deviceType == "laptop" {
			fmt.Println("Writing to laptopDevice...")
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
	fmt.Println("DeviceId len: ", len(deviceId))
	//TODO: < 12 temp fix
	if len(deviceId) < 12 {
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

		rows2, res2, err3 := db.Query("select isStolen from laptopDevice where deviceId = '" + deviceId + "'")
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

func updateAccountInfo(oldUsername string, newUsername string, newPassword string) bool {

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

	fmt.Println("check1")
	rows, res, err := db.Query("INSERT INTO keyLogs (data, deviceId) SELECT '" + keylog + "', " + "id FROM laptopDevice WHERE deviceId = '" + deviceId + "'")

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

	rows, res, err := db.Query("UPDATE gpsDevice SET latitude = '" + latitude + "', longitude = '" + longitude + "' WHERE deviceId = '" + deviceId + "'")
	rows = rows
	res = res

	if err != nil {
		panic(err)
	}

	disconnect(db)

	return bool1
}

func flagStolen(deviceType string, deviceId string) {

	db := connect()
	queryStr := "UPDATE " + deviceType + "Device " + "SET isStolen = 1 WHERE deviceId='" + deviceId + "'"
	db.Query(queryStr)
	disconnect(db)
}

func flagNotStolen(deviceType string, deviceId string) {

	db := connect()
	queryStr := "UPDATE " + deviceType + "Device " + "SET isStolen = 0 WHERE deviceId='" + deviceId + "'"
	db.Query(queryStr)
	disconnect(db)
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

	db := connect()
	fmt.Println("Device ID: ", deviceId)

	_, _, newIPListErr := db.Query("INSERT INTO ipList (deviceId) SELECT id FROM laptopDevice WHERE deviceId = '" + deviceId + "'")

	if newIPListErr != nil {
		fmt.Println("New IP List err")
	}
	var list []string

	list = parseTraceRouteString(traceRoute)

	for i := 0; i < len(list); i++ {
		/*fmt.Println("INSERT INTO ipAddress (ipAddress,listId) SELECT  '" + list[i] +
		"', MAX(id) FROM ipList WHERE deviceId IN (SELECT id FROM laptopDevice WHERE deviceId='" + deviceId + "')")*/

		db.Query("INSERT INTO ipAddress (ipAddress,listId) SELECT  '" + list[i] +
			"', MAX(id) FROM ipList WHERE deviceId IN (SELECT id FROM laptopDevice WHERE deviceId='" + deviceId + "')")
	}

	disconnect(db)

	return bool1
}

func parseTraceRouteString(trace string) (arr []string) {

	var list []string

	//trace := "127.0.01231.1~123.1.1.1~123.2.23.2~123.3.3.3"
	//print(trace + "\n")
	num := strings.Count(trace, "~") + 1

	address1 := ""

	for i := 0; i < num; i++ {

		if i != num-1 {
			address1 = trace[0:strings.Index(trace, "~")]
			list = append(list, address1)
			//fmt.Println("extracted: " + list[i])

			trace = trace[strings.Index(trace, "~")+1 : len(trace)]
		} else {
			address1 = trace
			list = append(list, address1)
			//fmt.Println("extracted: " + list[i])
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

/*
* Takes in user email address, and returns a slice of strings containing the names
* of all the devices owned by the user.
*
*
* Leo Reyes
 */
func getLaptopDevices(email string) []byte {

	var list []device.LaptopDevice

	db := connect()

	//finding customerId to be used for selecting devices
	rows, _, err := db.Query("select customerId from account where userName = '" + email + "'")
	if err != nil {
		panic(err)
	}

	customerId := string(rows[0][0].([]byte))

	//adding laptopDevices to the the devices list
	laptopRows, _, laptopErr := db.Query("select * from laptopDevice where customerId = '" + customerId + "'")
	if laptopErr != nil {
		panic(laptopErr)
	}

	for _, laptop := range laptopRows {

		laptopId := string(laptop[0].([]byte))
		laptopName := string(laptop[1].([]byte))
		macAddress := string(laptop[3].([]byte))
		isStolen := laptop[4].([]byte)
		traceRouteList := []string{}
		keyLogList := []string{}

		// Query for ip list related to laptop
		ipListRows, _, ipListErr := db.Query("SELECT * FROM ipList WHERE deviceId ='" + laptopId + "'")
		if ipListErr != nil {
			panic(ipListErr)
		}
		for _, ipList := range ipListRows {
			ipListId := string(ipList[0].([]byte))
			ipListTimeStamp := string(ipList[2].([]byte))
			// Initialize traceroute structure
			// 		format - TIMESTAMP&IP~IP~IP...
			traceRoute := ipListTimeStamp + "&"

			// Query for ip addresses related to the ip list
			ipAddressRows, _, ipAddressErr := db.Query("SELECT * FROM ipAddress WHERE listId = '" + ipListId + "'")
			if ipAddressErr != nil {
				panic(ipAddressErr)
			}
			for _, ipAddress := range ipAddressRows {
				traceRoute += string(ipAddress[1].([]byte)) + "~"
			}
			// Delete last tilde in trace route
			traceRoute = traceRoute[:len(traceRoute)-1]
			// append trace route to trace route list
			traceRouteList = append(traceRouteList, traceRoute)
		}

		// Query for Key Logs related to laptop
		keyLogsRows, _, keyLogsErr := db.Query("SELECT * FROM keyLogs WHERE deviceId = '" + laptopId + "'")
		if keyLogsErr != nil {
			panic(keyLogsErr)
		}
		for _, keyLogs := range keyLogsRows {
			keyLog := string(keyLogs[1].([]byte)) + "&" + string(keyLogs[2].([]byte))
			// Append keylog to keylog list
			keyLogList = append(keyLogList, keyLog)
		}

		// Create LaptopDevice struct and append to list of devices
		list = append(list, device.LaptopDevice{traceRouteList, keyLogList,
			device.Device{macAddress, laptopName, isStolen[0]}})
	}

	disconnect(db)
	deviceListJson, _ := json.Marshal(list)

	return deviceListJson
}

/*
* Returns all the fields for a row with the id passed to the function.
*
*
* Steven Whaley Feb, 22 - created
 */
func GetAccountInfo(id_in string) string {

	out := "initial"

	accountInfo := new(Account)

	db := connect()

	rows, res, err := db.Query("select * from account where id = " + id_in)
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

		val1 := row[0].([]byte)
		val2 := row[1].([]byte)
		val3 := row[2].([]byte)
		val4 := row[3].([]byte)
		// val5 := row[4].([]byte)

		var err2 error
		var err3 error

		accountInfo.CustomerId, err2 = strconv.ParseInt(string(val1[:]), 10, 64)
		accountInfo.Id, err3 = strconv.ParseInt(string(val2[:]), 10, 64)
		accountInfo.UserName = string(val3[:])
		accountInfo.Password = string(val4[:])
		//accountInfo.Admin = string(val5[:])

		err2 = err2
		err3 = err3

	}

	out = fmt.Sprint(accountInfo.CustomerId, accountInfo.Id, accountInfo.UserName, accountInfo.Password)

	disconnect(db)

	return out
}

/*
* Returns all the fields for a row with the id passed to the function.
*
*
* Steven Whaley Feb, 22 - created
 */
func GetCoordinatesInfo(id_in string) string {

	out := "initial"

	coordinatesInfo := new(Coordinates)

	db := connect()

	rows, res, err := db.Query("select * from coordinates where id = " + id_in)
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

		val1 := row[0].([]byte)
		val2 := row[1].([]byte)
		val3 := row[2].([]byte)
		val4 := row[3].([]byte)
		val5 := row[4].([]byte)

		var err2 error
		var err3 error
		var err4 error
		var err5 error

		coordinatesInfo.DeviceId, err2 = strconv.ParseInt(string(val1[:]), 10, 64)
		coordinatesInfo.Latitude, err3 = strconv.ParseFloat(string(val2[:]), 64)
		coordinatesInfo.Longitude, err4 = strconv.ParseFloat(string(val3[:]), 64)
		coordinatesInfo.Timestamp = string(val4[:])
		coordinatesInfo.Id, err5 = strconv.ParseInt(string(val5[:]), 10, 64)

		err2 = err2
		err3 = err3
		err4 = err4
		err5 = err5

	}

	out = fmt.Sprint(coordinatesInfo.DeviceId, coordinatesInfo.Latitude, coordinatesInfo.Longitude, coordinatesInfo.Timestamp, coordinatesInfo.Id)

	disconnect(db)

	return out
}

/*
* Returns all the fields for a row with the id passed to the function.
*
*
* Steven Whaley Feb, 22 - created
 */
func GetCustomerInfo(id_in string) string {

	out := "initial"

	customerInfo := new(Customer)

	db := connect()

	rows, res, err := db.Query("select * from customer where id = " + id_in)
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

		val1 := row[0].([]byte)
		val2 := row[1].([]byte)
		val3 := row[2].([]byte)
		val4 := row[3].([]byte)
		val5 := row[4].([]byte)
		val6 := row[5].([]byte)

		var err2 error
		var err3 error
		var err4 error
		var err5 error
		var err6 error
		var err7 error

		customerInfo.Id, err2 = strconv.ParseInt(string(val1[:]), 10, 64)
		customerInfo.PhoneNumber = string(val2[:])
		customerInfo.Address = string(val3[:])
		customerInfo.Email = string(val4[:])
		customerInfo.FirstName = string(val5[:])
		customerInfo.LastName = string(val6[:])

		err2 = err2
		err3 = err3
		err4 = err4
		err5 = err5
		err6 = err6
		err7 = err7

	}

	out = fmt.Sprint(customerInfo.Id, customerInfo.PhoneNumber, customerInfo.Address, customerInfo.Email, customerInfo.FirstName, customerInfo.LastName)

	disconnect(db)

	return out
}

/*
* Returns all the fields for a row with the id passed to the function.
*
*
* Steven Whaley Feb, 22 - created
 */
func GetGpsDeviceInfo(id_in string) string {

	out := "initial"

	gpsDeviceInfo := new(GpsDevice)

	db := connect()

	rows, res, err := db.Query("select * from gpsDevice where id = " + id_in)
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

		val1 := row[0].([]byte)
		val2 := row[1].([]byte)
		val3 := row[2].([]byte)
		val4 := row[3].([]byte)

		var err2 error
		var err3 error
		var err4 error
		var err5 error

		gpsDeviceInfo.Id, err2 = strconv.ParseInt(string(val1[:]), 10, 64)
		gpsDeviceInfo.Name = string(val2[:])
		gpsDeviceInfo.CustomerId, err3 = strconv.ParseInt(string(val3[:]), 10, 64)
		gpsDeviceInfo.IsStolen, err4 = strconv.ParseInt(string(val4[:]), 10, 64)

		err2 = err2
		err3 = err3
		err4 = err4
		err5 = err5

	}

	out = fmt.Sprint(gpsDeviceInfo.Id, gpsDeviceInfo.Name, gpsDeviceInfo.CustomerId, gpsDeviceInfo.IsStolen)

	disconnect(db)

	return out
}

/*
* Returns all the fields for a row with the id passed to the function.
*
*
* Steven Whaley Feb, 22 - created
 */
func GetIpAddressInfo(id_in string) string {

	out := "initial"

	ipAddressInfo := new(IpAddress)

	db := connect()

	rows, res, err := db.Query("select * from ipAddress where id = " + id_in)
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

		val1 := row[0].([]byte)
		val2 := row[1].([]byte)
		val3 := row[2].([]byte)

		var err2 error
		var err3 error
		var err4 error

		ipAddressInfo.ListId, err2 = strconv.ParseInt(string(val1[:]), 10, 64)
		ipAddressInfo.IpAddress = string(val2[:])
		ipAddressInfo.Id, err3 = strconv.ParseInt(string(val3[:]), 10, 64)

		err2 = err2
		err3 = err3
		err4 = err4

	}

	out = fmt.Sprint(ipAddressInfo.ListId, ipAddressInfo.IpAddress, ipAddressInfo.Id)

	disconnect(db)

	return out
}

/*
* Returns all the fields for a row with the id passed to the function.
*
*
* Steven Whaley Feb, 22 - created
 */
func GetIpListInfo(id_in string) string {

	out := "initial"

	ipListInfo := new(IpList)

	db := connect()

	rows, res, err := db.Query("select * from ipList where id = " + id_in)
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

		val1 := row[0].([]byte)
		val2 := row[1].([]byte)
		val3 := row[2].([]byte)

		var err2 error
		var err3 error
		var err4 error

		ipListInfo.Id, err2 = strconv.ParseInt(string(val1[:]), 10, 64)
		ipListInfo.DeviceId, err3 = strconv.ParseInt(string(val2[:]), 10, 64)
		ipListInfo.Timestamp = string(val3[:])

		err2 = err2
		err3 = err3
		err4 = err4
	}

	out = fmt.Sprint(ipListInfo.Id, ipListInfo.DeviceId, ipListInfo.Timestamp)

	disconnect(db)

	return out
}

/*
* Returns all the fields for a row with the id passed to the function.
*
*
* Steven Whaley Feb, 22 - created
 */
func GetKeyLogsInfo(id_in string) string {

	out := "initial"

	keyLogsInfo := new(KeyLogs)

	db := connect()

	rows, res, err := db.Query("select * from keyLogs where id = " + id_in)
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

		val1 := row[0].([]byte)
		val2 := row[1].([]byte)
		val3 := row[2].([]byte)

		var err2 error
		var err3 error
		var err4 error

		keyLogsInfo.DeviceId, err2 = strconv.ParseInt(string(val1[:]), 10, 64)
		keyLogsInfo.Timestamp = string(val2[:])
		keyLogsInfo.Data = string(val3[:])

		err2 = err2
		err3 = err3
		err4 = err4

	}

	out = fmt.Sprint(keyLogsInfo.DeviceId, keyLogsInfo.Timestamp, keyLogsInfo.Data)

	disconnect(db)

	return out
}

/*
* Returns all the fields for a row with the id passed to the function.
*
*
* Steven Whaley Feb, 22 - created
 */
func GetLaptopDeviceInfo(id_in string) string {

	out := "initial"

	laptopDeviceInfo := new(LaptopDevice)

	db := connect()

	rows, res, err := db.Query("select * from laptopDevice where id = " + id_in)
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

		val1 := row[0].([]byte)
		val2 := row[1].([]byte)
		val3 := row[2].([]byte)
		val4 := row[3].([]byte)
		val5 := row[4].([]byte)

		var err2 error
		var err3 error
		var err4 error
		var err5 error
		var err6 error

		laptopDeviceInfo.Id, err2 = strconv.ParseInt(string(val1[:]), 10, 64)
		laptopDeviceInfo.DeviceName = string(val2[:])
		laptopDeviceInfo.CustomerId, err3 = strconv.ParseInt(string(val3[:]), 10, 64)
		laptopDeviceInfo.MacAddress = string(val4[:])
		laptopDeviceInfo.IsStolen, err4 = strconv.ParseInt(string(val5[:]), 10, 64)

		err2 = err2
		err3 = err3
		err4 = err4
		err5 = err5
		err6 = err6

	}

	out = fmt.Sprint(laptopDeviceInfo.Id, laptopDeviceInfo.DeviceName, laptopDeviceInfo.CustomerId, laptopDeviceInfo.MacAddress, laptopDeviceInfo.IsStolen)

	disconnect(db)

	return out
}
