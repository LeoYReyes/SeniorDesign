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
*   changes:
*       Steven Whaley:
*       February 18 - The server is now being implemented in Google Go instead of Java. Started researching mysql libraries
*             in go, and wrote some database connection code. Experimented with some simple database connection code.
*
*       Steven Whaley:
*       February 19 - Successfully connected to database, but code is all in main and isn't ready to
*               interface with server code. Worked on getter functions, such as GetAccountInfo().
*             Currently working on outputting in JSON format.
*
*       Steven Whaley:
*       February 20 - VerifyAccountInfo() takes in strings for username and password and
*             returns two booleans for whether they match the database values or not.
*
*     Steven Whaley:
*     February 23 - Issue Parsing some mysql datatypes properly. Testing code using my main function.
*
*     Steven Whaley:
*     February 26 - SignUp() takes in user information entered on the sign up page, and creates account and customer entries
*             in the database.
*
*     Steven Whaley:
*     February 27 - We have now tested using the website for sign up and login. It works, but we now need to properly
*             implement a request system. I
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

func parsePayload(payload []byte) []string {
	str := []string{}
	pos := 1
	for index, element := range payload {
		if element == 0x1B {
			str = append(str, string(payload[pos:index-1]))
			pos = index + 2
		}
	}
	return str
}

func processRequest(req *CustomProtocol.Request) {
	payload := parsePayload(req.Payload)
	switch req.OpCode {
	case CustomProtocol.NewAccount:
	case CustomProtocol.NewDevice:
	case CustomProtocol.UpdateDeviceGPS:
	case CustomProtocol.UpdateDeviceIP:
	case CustomProtocol.UpdateDeviceKeylog:
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
	case CustomProtocol.GetDevice:
	case CustomProtocol.SetDevice:
	case CustomProtocol.GetDeviceList:
	case CustomProtocol.CheckDeviceStolen: // in: id string, out: boolean
	case CustomProtocol.UpdateUserKeylogData: // in: id string, string keylog out: boolean
	case CustomProtocol.UpdateUserIPTraceData: // in: id string, traceroute string, out: boolean
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
*  IsDeviceStolen(deviceId string) (bool) takes in device id and return a
*  boolean indicating whether that device is marked stolen.
*
*
* Steven Whaley Mar, 1 - created
 */

func IsDeviceStolen(deviceId string) (bool) {

  bool1 := false
  
  db := connect()

  rows, res, err := db.Query("select isStolen from gpsDevice where id = '" + deviceId + "'")
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

  rows2, res2, err3 := db.Query("select isStolen from laptopDevice where id = '" + deviceId + "'")
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
* Steven Whaley Feb, 26 - created
 */
func GetUserDevices(email string) []string {

	customerId := "initial"

	var list []string

	db := connect()

	//finding customerId to be used for selecting devices
	rows, res, err := db.Query("select customerId from account where userName = '" + email + "'")
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
		customerId = string(val1[:])

	}

	//adding laptopDevices to the the devices list
	rows2, res2, err2 := db.Query("select deviceName from laptopDevice where customerId = '" + customerId + "'")
	if err2 != nil {
		panic(err2)
	}

	res2 = res2

	for _, row2 := range rows2 {
		for _, col2 := range row2 {
			if col2 == nil {
				// col has NULL value
			} else {
				// Do something with text in col (type []byte)
			}
		}

		val2 := row2[0].([]byte)

		list = append(list, string(val2[:]))

	}

	//adding gpsDevices to the devices list
	rows3, res3, err3 := db.Query("select name from gpsDevice where customerId = '" + customerId + "'")
	if err3 != nil {
		panic(err3)
	}

	res3 = res3

	for _, row3 := range rows3 {
		for _, col3 := range row3 {
			if col3 == nil {
				// col has NULL value
			} else {
				// Do something with text in col (type []byte)
			}
		}

		val3 := row3[0].([]byte)

		list = append(list, string(val3[:]))
	}

	disconnect(db)

	return list
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
