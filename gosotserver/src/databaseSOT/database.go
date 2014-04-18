/*
 *	DATABASE.GO
 *
 * 	Created			- Steven Whaley: February 18, 2014
 *	Last updated	- Steven Whaley: April 17, 2014
 *
 * 	OVERVIEW:
 *
 *  		This is the database portion of the server. It provides functionality for interacting with the database. Actions
 *  		such as user sign up and login must use functions provided in this file to update and query the database.
 *  		See the comments for each method for details.
 *
 *		user := "root"
 *		pass := "toor"
 *		dbname := "trackerdb"
 *		proto := "tcp"
 *		addr := "127.0.0.1:3306"
 *
 *   Useful Links:
 *
 *   	Mysql driver for google go	-	https://github.com/ziutek/mymysql
 *
 *   	Golang documentation		-	http://golang.org/doc/
 *
 */

package databaseSOT

import (
	"CustomProtocol"
	"device"
	"encoding/json"
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
	"strconv"
	"strings"
)

/*
 * These structs represent tables in the database.
 *
 *	Steven Whaley - February 2014
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

/*
 *	Processes server requests and calls database functions
 *	to perform requested functionality. Resulting response
 *	is sent back.
 *
 *	Steven Whaley - March 2014
 */
func processRequest(req *CustomProtocol.Request) {

	payload := CustomProtocol.ParsePayload(req.Payload)
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
		//TODO: temp fix < 12
		if len(payload[0]) < 12 {
			flagNotStolen("gps", payload[0])
		} else {
			flagNotStolen("laptop", payload[0])
		}
		res := make([]byte, 2)
		res[0] = 1 //TO DO CHANGE
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
		res = append(res, 0x1B)
		res = append(res, getGpsDevices(payload[0])...)
		req.Response <- res
	case CustomProtocol.CheckDeviceStolen:
		isStolen := IsDeviceStolen(payload[0])
		res := make([]byte, 1)
		if isStolen == true {
			res[0] = 1
		} else {
			res[0] = 0
		}
		req.Response <- res
	case CustomProtocol.UpdateUserKeylogData:
		boolResult := UpdateKeylog(payload[0], payload[1])
		res := make([]byte, 1)
		if boolResult == true {
			res[0] = 1
		} else {
			res[0] = 0
		}
		req.Response <- res
	case CustomProtocol.UpdateUserIPTraceData:
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

/*
 * 	Used to form connection with the database.
 *
 * 	Steven Whaley Feb, 27 - created
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
		fmt.Println("Database Connection Error:", err)
	}

	return db
}

/*
 * 	Used to close connection with the database.
 *
 * 	Steven Whaley Feb, 27 - created
 */

func disconnect(connection mysql.Conn) {

	connection.Close()

}

/*
 * 	Takes in user information entered on the sign up page, and creates account and customer entries
 * 	in the database.
 *
 * 	Steven Whaley Feb, 26 - created
 */

func SignUp(firstname string, lastname string, email string, phoneNumber string, password string) bool {

	db := connect()

	output := true

	if strings.Contains(firstname, "'") || strings.Contains(lastname, "'") || strings.Contains(email, "'") || strings.Contains(phoneNumber, "'") || strings.Contains(password, "'") {
		firstname = strings.Replace(firstname, "'", "\\'", -1)
		lastname = strings.Replace(lastname, "'", "\\'", -1)
		email = strings.Replace(email, "'", "\\'", -1)
		phoneNumber = strings.Replace(phoneNumber, "'", "\\'", -1)
		password = strings.Replace(password, "'", "\\'", -1)
	}

	res, _, err := db.Query("SELECT * FROM customer WHERE email = '" + email + "'")

	if err != nil {
		fmt.Println("Database Query Error:", err)
	}

	if len(res) != 0 {
		output = false
	} else {
		db.Query("INSERT INTO customer (firstName, lastName, email, phoneNumber) VALUES ('" + firstname + "', '" + lastname + "', '" + email + "', '" + phoneNumber + "')")

		db.Query("INSERT INTO account (userName, password, customerId) SELECT '" + email + "', '" + password + "', id FROM customer WHERE email='" + email + "'")
	}

	disconnect(db)

	return output
}

/*
 * 	Return a byte array of gps devices associated
 *	with the email input.
 *
 * 	Steven Whaley March 2014 - created
 */
func getGpsDevices(email string) []byte {
	var list []device.GPSDevice

	db := connect()

	rows, _, err := db.Query("select customerId from account where userName = '" + email + "'")
	if err != nil {
		fmt.Println("Database Query Error:", err)
	}

	customerId := string(rows[0][0].([]byte))

	gpsRows, _, gpsErr := db.Query("select * from gpsDevice where customerId = '" + customerId + "'")
	if gpsErr != nil {
		fmt.Println("Database Query Error:", gpsErr)
	}

	for _, gps := range gpsRows {

		gpsId := string(gps[0].([]byte))
		gpsName := string(gps[1].([]byte))
		deviceId := string(gps[3].([]byte))
		isStolen := gps[4].([]byte)

		gpsCoordRows, _, _ := db.Query("SELECT * FROM coordinates WHERE deviceId='" + gpsId + "'")

		coordList := []string{}
		for _, coordRow := range gpsCoordRows {
			str := ""
			str = string(coordRow[3].([]byte)) + "&" + string(coordRow[1].([]byte)) + string(0x1B) + string(coordRow[2].([]byte))
			coordList = append(coordList, str)
		}
		list = append(list, device.GPSDevice{coordList, device.Device{deviceId, gpsName, isStolen[0]}})
	}

	disconnect(db)
	deviceListJson, _ := json.Marshal(list)
	return deviceListJson
}

/*
 *	Registers a new device based on the type, id, name
 *	and user email provided.
 *
 *
 *	Steven Whaley Mar, 18 - created
 */

func registerNewDevice(deviceType string, deviceId string, deviceName string, email string) bool {
	db := connect()
	output := true
	var res []mysql.Row
	if strings.Contains(deviceType, "'") || strings.Contains(deviceId, "'") || strings.Contains(deviceName, "'") || strings.Contains(email, "'") {
		deviceType = strings.Replace(deviceType, "'", "\\'", -1)
		deviceId = strings.Replace(deviceId, "'", "\\'", -1)
		deviceName = strings.Replace(deviceName, "'", "\\'", -1)
		email = strings.Replace(email, "'", "\\'", -1)
	}
	if deviceType == "gps" {
		res, _, _ = db.Query("SELECT * FROM gpsDevice WHERE deviceId = '" + deviceId + "'")
	} else if deviceType == "laptop" {
		res, _, _ = db.Query("SELECT * FROM laptopDevice WHERE deviceId = '" + deviceId + "'")
	}

	if len(res) != 0 {
		output = false
		fmt.Println("check")
	} else {
		if deviceType != "gps" && deviceType != "laptop" {
			print("invalid device type")
		} else {
			if deviceType == "gps" {
				fmt.Println("Writing to gpsDevice...")
				db.Query("INSERT INTO gpsDevice (deviceName, deviceId, customerId) SELECT '" + deviceName + "', '" + deviceId + "', id FROM customer WHERE email='" + email + "'")
			} else if deviceType == "laptop" {
				fmt.Println("Writing to laptopDevice...")
				db.Query("INSERT INTO laptopDevice (deviceName, deviceId, customerId) SELECT '" + deviceName + "', '" + deviceId + "', id FROM customer WHERE email='" + email + "'")
			}
		}
	}
	disconnect(db)

	return output
}

/*
 *  IsDeviceStolen(deviceId string) (bool) takes in device id and return a
 *  boolean indicating whether that device is marked stolen.
 *
 *
 *	Steven Whaley Mar, 1 - created
 */

func IsDeviceStolen(deviceId string) bool {

	bool1 := false

	db := connect()

	//TODO: < 12 temp fix
	if len(deviceId) < 12 {
		rows, _, err := db.Query("select isStolen from gpsDevice where deviceId = '" + deviceId + "'")
		if err != nil {
			fmt.Println("Database Query Error:", err)
		}

		for _, row := range rows {

			val1 := row[0].([]byte)

			temp, _ := strconv.ParseInt(string(val1[:]), 10, 64)

			if temp == 1 {
				bool1 = true
			} else {
				bool1 = false
			}
		}
	} else {

		rows2, _, err3 := db.Query("select isStolen from laptopDevice where deviceId = '" + deviceId + "'")
		if err3 != nil {
			fmt.Println("Database Query Error:", err3)
		}

		for _, row := range rows2 {

			val2 := row[0].([]byte)

			temp2, _ := strconv.ParseInt(string(val2[:]), 10, 64)

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
 *	Updates account fields for already existing account.
 *
 *
 *	Steven Whaley Mar, 23
 */

func updateAccountInfo(oldUsername string, newUsername string, newPassword string) bool {

	bool1 := true

	db := connect()

	if strings.Contains(newUsername, "'") || strings.Contains(newPassword, "'") {
		newUsername = strings.Replace(newUsername, "'", "\\'", -1)
		newPassword = strings.Replace(newPassword, "'", "\\'", -1)
	}

	_, _, err := db.Query("UPDATE account SET userName = '" + newUsername + "', password = '" + newPassword + "' WHERE userName = '" + oldUsername + "'")

	if err != nil {
		fmt.Println("Database Query Error:", err)
	}

	disconnect(db)

	return bool1
}

/*
 *  Inserts new keylog data into database.
 *
 *
 * 	Steven Whaley - Mar, 1 - created
 * 	Steven Whaley - Mar, 17 - seems to be working
 */

func UpdateKeylog(deviceId string, keylog string) bool {

	bool1 := true

	db := connect()

	if strings.Contains(keylog, "'") {
		keylog = strings.Replace(keylog, "'", "\\'", -1)
	}

	_, _, err := db.Query("INSERT INTO keyLogs (data, deviceId) SELECT '" + keylog + "', " + "id FROM laptopDevice WHERE deviceId = '" + deviceId + "'")

	if err != nil {
		fmt.Println("Database Query Error:", err)
	}

	disconnect(db)

	return bool1
}

/*
 *	Updates the latitude and longitude coordinates associated
 *	with a gps device.
 *
 *	Steven Whaley - March 23, 2014 - created
 */
func updateDeviceGps(deviceId string, latitude string, longitude string) bool {
	bool1 := true

	db := connect()

	_, _, err := db.Query("INSERT INTO coordinates(latitude, longitude, deviceId) SELECT '" + latitude + "', '" + longitude + "', id FROM gpsDevice WHERE deviceId = '" + deviceId + "'")

	if err != nil {
		fmt.Println("Database Query Error:", err)
	}

	disconnect(db)

	return bool1
}

/*
 *	Sets isStolen to 1 in the database to indicate that a device
 *	is stolen.
 *
 *	Steven Whaley - March 2014
 */
func flagStolen(deviceType string, deviceId string) {

	db := connect()
	queryStr := "UPDATE " + deviceType + "Device " + "SET isStolen = 1 WHERE deviceId='" + deviceId + "'"
	db.Query(queryStr)
	disconnect(db)
}

/*
 *	Sets isStolen to 1 in the database to indicate that a device
 *	is stolen.
 *
 *	Steven Whaley - March 2014
 */
func flagNotStolen(deviceType string, deviceId string) {

	db := connect()
	queryStr := "UPDATE " + deviceType + "Device " + "SET isStolen = 0 WHERE deviceId='" + deviceId + "'"
	db.Query(queryStr)
	disconnect(db)
}

/*
 *	UpdateTraceRoute calls parseTraceRouteString
 *
 *	Steven Whaley - March 2014
 */
func UpdateTraceRoute(deviceId string, traceRoute string) bool {

	bool1 := true

	db := connect()

	_, _, newIPListErr := db.Query("INSERT INTO ipList (deviceId) SELECT id FROM laptopDevice WHERE deviceId = '" + deviceId + "'")

	if newIPListErr != nil {
		fmt.Println("New IP List err", newIPListErr)
	}
	var list []string

	list = parseTraceRouteString(traceRoute)

	for i := 0; i < len(list); i++ {
		db.Query("INSERT INTO ipAddress (ipAddress,listId) SELECT  '" + list[i] +
			"', MAX(id) FROM ipList WHERE deviceId IN (SELECT id FROM laptopDevice WHERE deviceId='" + deviceId + "')")
	}

	disconnect(db)

	return bool1
}

/*
 *	The traceroute string is delimited by ~
 *	This function parse out the the ip addresses into
 *	an array.
 *
 *	Steven Whaley - March 2014
 */
func parseTraceRouteString(trace string) (arr []string) {

	var list []string

	num := strings.Count(trace, "~") + 1

	address1 := ""

	for i := 0; i < num; i++ {

		if i != num-1 {
			address1 = trace[0:strings.Index(trace, "~")]
			list = append(list, address1)
			trace = trace[strings.Index(trace, "~")+1 : len(trace)]
		} else {
			address1 = trace
			list = append(list, address1)
		}
	}
	return list
}

/*
 *	Returns True or False
 *
 *	Steven Whaley Mar, 1 - created
 */
func IsGpsDevice(deviceId string) bool {

	bool1 := false

	db := connect()

	rows, res, err := db.Query("select * from gpsDevice where id = '" + deviceId + "'")
	if err != nil {
		fmt.Println("Database Query Error:", err)
	}

	res = res

	for _, row := range rows {

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
 *	Takes in user information entered on the sign up page, and creates account and customer entries
 *	in the database.
 *
 *	Steven Whaley Feb, 26 - created
 */

func VerifyAccountInfo(username string, password string) (bool, bool) {

	bool1 := false
	bool2 := false

	accountInfo := new(Account)

	db := connect()

	rows, res, err := db.Query("select * from account")
	if err != nil {
		fmt.Println("Database Query Error:", err)
	}
	res = res

	for _, row := range rows {

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
 *	Takes in user email address, and returns a slice of strings containing the names
 *	of all the devices owned by the user.
 *
 *	Steven Whaley
 */
func getLaptopDevices(email string) []byte {

	var list []device.LaptopDevice

	db := connect()

	//finding customerId to be used for selecting devices
	rows, _, err := db.Query("select customerId from account where userName = '" + email + "'")
	if err != nil {
		fmt.Println("Database Query Error:", err)
	}

	customerId := string(rows[0][0].([]byte))

	//adding laptopDevices to the the devices list
	laptopRows, _, laptopErr := db.Query("select * from laptopDevice where customerId = '" + customerId + "'")
	if laptopErr != nil {
		fmt.Println(laptopErr)
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
			fmt.Println("Database Query Error:", ipListErr)
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
				fmt.Println(ipAddressErr)
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
			fmt.Println(keyLogsErr)
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
