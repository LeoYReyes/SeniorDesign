// testproj project testproj.go
// personal note: export GOPATH=/Users/stevenwhaley/go/

package databaseSOT

import (
	"fmt"
	"os"
	"strconv"
	//"bytes"
	//"encoding/json"

	"github.com/ziutek/mymysql/mysql"
	//_ "github.com/ziutek/mymysql/native" // Native engine
	_ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
)

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

func printOK() {
	fmt.Println("OK")
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

func SignUp(firstname string, lastname string, email string, phoneNumber string, password string) {

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

	db.Query("INSERT INTO customer (firstName, lastName, email, phoneNumber) VALUES ('" + firstname + "', '" + lastname + "', '" + email + "', '" + phoneNumber + "')")

	db.Query("INSERT INTO account (userName, password, customerId) SELECT '" + email + "', '" + password + "', id FROM customer WHERE email='" + email + "'")

	checkError(db.Close())

	return
}

func VerifyAccountInfo(username string, password string) (bool, bool) {

	bool1 := false
	bool2 := false

	accountInfo := new(Account)

	user := "root"
	pass := ""
	dbname := "trackerdb"
	//proto := "unix"
	//addr := "/var/run/mysqld/mysqld.sock"
	proto := "tcp"
	addr := "127.0.0.1:3306"

	db := mysql.New(proto, "", addr, user, pass, dbname)

	err := db.Connect()
	if err != nil {
		panic(err)
	}

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

		//val2 := row[0].([]byte)
		//val3 := row[1].([]byte)
		val4 := row[2].([]byte)
		val5 := row[3].([]byte)

		accountInfo.UserName = string(val4[:])
		accountInfo.Password = string(val5[:])

		//  fmt.Println(accountInfo.UserName)
		//  fmt.Println(accountInfo.Password)

		//jsonx, _ := json.Marshal(accountInfo)
		//  fmt.Println(string(jsonx))

		//  m := &Account{UserName: "leo", Password: "pass"}
		//json1, _ := json.Marshal(m)
		//fmt.Println(string(json1))

		//x := &Account{UserName: "test", Password: "pass"}
		//acc1, _ := json.Marshal(x)
		//fmt.Println(string(acc1))

		if accountInfo.UserName == username {
			bool1 = true
		}
		if accountInfo.Password == password {
			bool2 = true
		}

		//  out = string(jsonx)

	}

	checkError(db.Close())
	//fmt.Print("Connection Closed... ")
	//printOK()

	return bool1, bool2
}

/*func GetUserDevices(customerId string) (string)  {

    out := "initial"

    user := "root"
    pass := ""
    dbname := "trackerdb"
    //proto := "unix"
    //addr := "/var/run/mysqld/mysqld.sock"
    proto := "tcp"
    addr := "127.0.0.1:3306"

    db := mysql.New(proto, "", addr, user, pass, dbname)

    err := db.Connect()
    if err != nil {
        panic(err)
    }

    rows, res, err := db.Query("select name from gpsDevice, laptopDevice where customerId = 15")
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

        out = string(val1[:])
    }

    checkError(db.Close())


    return out
}*/

func GetAccountInfo(id_in string) string {

	out := "initial"

	accountInfo := new(Account)

	user := "root"
	pass := ""
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"

	db := mysql.New(proto, "", addr, user, pass, dbname)

	err := db.Connect()
	if err != nil {
		panic(err)
	}

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

	checkError(db.Close())

	return out
}

func GetCoordinatesInfo(id_in string) string {

	out := "initial"

	coordinatesInfo := new(Coordinates)

	user := "root"
	pass := ""
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"

	db := mysql.New(proto, "", addr, user, pass, dbname)

	err := db.Connect()
	if err != nil {
		panic(err)
	}

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

	checkError(db.Close())

	return out
}

func GetCustomerInfo(id_in string) string {

	out := "initial"

	customerInfo := new(Customer)

	user := "root"
	pass := ""
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"

	db := mysql.New(proto, "", addr, user, pass, dbname)

	err := db.Connect()
	if err != nil {
		panic(err)
	}

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

	checkError(db.Close())

	return out
}

func GetGpsDeviceInfo(id_in string) string {

	out := "initial"

	gpsDeviceInfo := new(GpsDevice)

	user := "root"
	pass := ""
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"

	db := mysql.New(proto, "", addr, user, pass, dbname)

	err := db.Connect()
	if err != nil {
		panic(err)
	}

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

	checkError(db.Close())

	return out
}

func GetIpAddressInfo(id_in string) string {

	out := "initial"

	ipAddressInfo := new(IpAddress)

	user := "root"
	pass := ""
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"

	db := mysql.New(proto, "", addr, user, pass, dbname)

	err := db.Connect()
	if err != nil {
		panic(err)
	}

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

	checkError(db.Close())

	return out
}

func GetIpListInfo(id_in string) string {

	out := "initial"

	ipListInfo := new(IpList)

	user := "root"
	pass := ""
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"

	db := mysql.New(proto, "", addr, user, pass, dbname)

	err := db.Connect()
	if err != nil {
		panic(err)
	}

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

	checkError(db.Close())

	return out
}

func GetKeyLogsInfo(id_in string) string {

	out := "initial"

	keyLogsInfo := new(KeyLogs)

	user := "root"
	pass := ""
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"

	db := mysql.New(proto, "", addr, user, pass, dbname)

	err := db.Connect()
	if err != nil {
		panic(err)
	}

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

	checkError(db.Close())

	return out
}

func GetLaptopDeviceInfo(id_in string) string {

	out := "initial"

	laptopDeviceInfo := new(LaptopDevice)

	user := "root"
	pass := ""
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"

	db := mysql.New(proto, "", addr, user, pass, dbname)

	err := db.Connect()
	if err != nil {
		panic(err)
	}

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

	checkError(db.Close())

	return out
}

func disconnect() {

}

//func main() {

//fmt.Println(GetAccountInfo("12"))

//fmt.Println(GetCoordinatesInfo("1"))

//fmt.Println(GetCustomerInfo("15"))

//fmt.Println(GetGpsDeviceInfo("15"))

//fmt.Println(GetIpAddressInfo("1"))

//fmt.Println(GetIpListInfo("1"))

//fmt.Println(GetKeyLogsInfo("1"))

//fmt.Println(GetLaptopDeviceInfo("1"))

//fmt.Println(GetUserDevices("15"))

//SignUp("steven", "whaley", "steven@facebook.gov", "911", "usersteve", "password1")

//fmt.Println(VerifyAccountInfo("sadfk", "eieiei"))

//fmt.Println(VerifyAccountInfo("wrongusernameexample", "369d841cdf0dd150a680931769e868d9e487452f"))

//fmt.Println(VerifyAccountInfo("leo@auburn.edu", "wrongpasswordexample"))
//}
