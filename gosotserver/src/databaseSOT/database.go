// personal note: export GOPATH=/Users/stevenwhaley/go/

package databaseSOT

import (
<<<<<<< HEAD
        "fmt"
        "os"
        "strconv"
        "github.com/ziutek/mymysql/mysql"
        //_ "github.com/ziutek/mymysql/native" // Native engine
         _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
    ) 
=======
	"fmt"
	"os"
	"strconv"
	//"bytes"
	//"encoding/json"
>>>>>>> FETCH_HEAD

	"github.com/ziutek/mymysql/mysql"
	//_ "github.com/ziutek/mymysql/native" // Native engine
	_ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
)

type Account struct {
<<<<<<< HEAD
  CustomerId int64
  Id int64
    UserName string
    Password string
    Admin bool

=======
	CustomerId int64
	Id         int64
	UserName   string
	Password   string
	Admin      bool
>>>>>>> FETCH_HEAD
}
type Coordinates struct {
	DeviceId  int64
	Latitude  float64
	Longitude float64
	Timestamp string
	Id        int64
}
type Customer struct {
<<<<<<< HEAD
    Id int64
  PhoneNumber string
    Address string
    Email string
    FirstName string
    LastName string
=======
	Id          int64
	PhoneNumber string
	Address     string
	Email       string
	FirstName   string
	LastName    string
>>>>>>> FETCH_HEAD
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

<<<<<<< HEAD
=======
func printOK() {
	fmt.Println("OK")
}
>>>>>>> FETCH_HEAD

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

<<<<<<< HEAD
func connect() (connection mysql.Conn){
    user := "root"
    pass := "toor"
    dbname := "trackerdb"
    proto := "tcp"
    addr := "127.0.0.1:3306"
=======
func SignUp(firstname string, lastname string, email string, phoneNumber string, password string) {

	user := "root"
	pass := "toor"
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"
>>>>>>> FETCH_HEAD

	db := mysql.New(proto, "", addr, user, pass, dbname)

<<<<<<< HEAD
    return db
}

func disconnect(connection mysql.Conn) {

  checkError(connection.Close())

}

func SignUp(firstname string, lastname string, email string, phoneNumber string, password string){

    db := connect()

    db.Query("INSERT INTO customer (firstName, lastName, email, phoneNumber) VALUES ('" + firstname + "', '" + lastname+ "', '" + email + "', '" + phoneNumber + "')") 
    
    db.Query("INSERT INTO account (userName, password, customerId) SELECT '" + email + "', '" + password + "', id FROM customer WHERE email='" + email + "'")
  
    disconnect(db)
    
    return
=======
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	db.Query("INSERT INTO customer (firstName, lastName, email, phoneNumber) VALUES ('" + firstname + "', '" + lastname + "', '" + email + "', '" + phoneNumber + "')")

	db.Query("INSERT INTO account (userName, password, customerId) SELECT '" + email + "', '" + password + "', id FROM customer WHERE email='" + email + "'")

	checkError(db.Close())

	return
>>>>>>> FETCH_HEAD
}

func VerifyAccountInfo(username string, password string) (bool, bool) {

<<<<<<< HEAD
    bool1 := false
    bool2 := false
    
    accountInfo := new (Account)

    db := connect()
=======
	bool1 := false
	bool2 := false

	accountInfo := new(Account)

	user := "root"
	pass := "toor"
	dbname := "trackerdb"
	//proto := "unix"
	//addr := "/var/run/mysqld/mysqld.sock"
	proto := "tcp"
	addr := "127.0.0.1:3306"
>>>>>>> FETCH_HEAD

	db := mysql.New(proto, "", addr, user, pass, dbname)

	err := db.Connect()
	if err != nil {
		panic(err)
	}

	rows, res, err := db.Query("select * from account")
	if err != nil {
		panic(err)
	}

<<<<<<< HEAD
        val4 := row[2].([]byte)
        val5 := row[3].([]byte)
        
        accountInfo.UserName = string(val4[:])
        accountInfo.Password = string(val5[:])
      

        if accountInfo.UserName == username{
            bool1 = true
        }
        if accountInfo.Password == password{
            bool2 = true
        }
    }

    disconnect(db)
    
    return bool1, bool2
}
=======
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
>>>>>>> FETCH_HEAD

	}

	checkError(db.Close())
	//fmt.Print("Connection Closed... ")
	//printOK()

	return bool1, bool2
}

func GetUserDevices(email string) ([]string)  {

<<<<<<< HEAD
    customerId := "initial"
    
    var list []string
    
    db := connect()
=======
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
>>>>>>> FETCH_HEAD

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
<<<<<<< HEAD
        customerId = string(val1[:])
     
    }

    //adding laptopDevices to the the devices list
    rows2, res2, err2 := db.Query("select deviceName from laptopDevice where customerId = '" + customerId + "'")
    if err2 != nil {
        panic(err2)
    }
=======
        val2 := row[1].([]byte)
        val3 := row[2].([]byte)
        val4 := row[3].([]byte)
        val5 := row[4].([]byte)

        out = string(val1[:])
    }

    checkError(db.Close())

>>>>>>> FETCH_HEAD

    res2 = res2
    

<<<<<<< HEAD
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

func GetAccountInfo(id_in string) (string)  {

    out := "initial"
    
    accountInfo := new (Account)

    db := connect()

    rows, res, err := db.Query("select * from account where id = " + id_in)
    if err != nil {
        panic(err)
    }
=======
func GetAccountInfo(id_in string) string {

	out := "initial"

	accountInfo := new(Account)

	user := "root"
	pass := ""
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"

	db := mysql.New(proto, "", addr, user, pass, dbname)
>>>>>>> FETCH_HEAD

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

<<<<<<< HEAD
    err2 = err2
    err3 = err3
=======
		err2 = err2
		err3 = err3
>>>>>>> FETCH_HEAD

	}

<<<<<<< HEAD
      out = fmt.Sprint(accountInfo.CustomerId, accountInfo.Id, accountInfo.UserName, accountInfo.Password)

    disconnect(db)
=======
	out = fmt.Sprint(accountInfo.CustomerId, accountInfo.Id, accountInfo.UserName, accountInfo.Password)

	checkError(db.Close())
>>>>>>> FETCH_HEAD

	return out
}

func GetCoordinatesInfo(id_in string) string {

	out := "initial"

<<<<<<< HEAD
    db := connect()
=======
	coordinatesInfo := new(Coordinates)

	user := "root"
	pass := ""
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"
>>>>>>> FETCH_HEAD

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

<<<<<<< HEAD
    err2 = err2
    err3 = err3
    err4 = err4
    err5 = err5
=======
		coordinatesInfo.DeviceId, err2 = strconv.ParseInt(string(val1[:]), 10, 64)
		coordinatesInfo.Latitude, err3 = strconv.ParseFloat(string(val2[:]), 64)
		coordinatesInfo.Longitude, err4 = strconv.ParseFloat(string(val3[:]), 64)
		coordinatesInfo.Timestamp = string(val4[:])
		coordinatesInfo.Id, err5 = strconv.ParseInt(string(val5[:]), 10, 64)
>>>>>>> FETCH_HEAD

		err2 = err2
		err3 = err3
		err4 = err4
		err5 = err5

	}

<<<<<<< HEAD
    disconnect(db)
=======
	out = fmt.Sprint(coordinatesInfo.DeviceId, coordinatesInfo.Latitude, coordinatesInfo.Longitude, coordinatesInfo.Timestamp, coordinatesInfo.Id)

	checkError(db.Close())
>>>>>>> FETCH_HEAD

	return out
}

func GetCustomerInfo(id_in string) string {

	out := "initial"

<<<<<<< HEAD
    db := connect()
=======
	customerInfo := new(Customer)

	user := "root"
	pass := ""
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"
>>>>>>> FETCH_HEAD

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

<<<<<<< HEAD
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
=======
	for _, row := range rows {
		for _, col := range row {
			if col == nil {
				// col has NULL value
			} else {
				// Do something with text in col (type []byte)
			}
		}
>>>>>>> FETCH_HEAD

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

<<<<<<< HEAD
    disconnect(db)
=======
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
>>>>>>> FETCH_HEAD

	checkError(db.Close())

	return out
}

func GetGpsDeviceInfo(id_in string) string {

	out := "initial"

<<<<<<< HEAD
    db := connect()
=======
	gpsDeviceInfo := new(GpsDevice)

	user := "root"
	pass := ""
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"
>>>>>>> FETCH_HEAD

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

<<<<<<< HEAD
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
=======
	for _, row := range rows {
		for _, col := range row {
			if col == nil {
				// col has NULL value
			} else {
				// Do something with text in col (type []byte)
			}
		}
>>>>>>> FETCH_HEAD

		val1 := row[0].([]byte)
		val2 := row[1].([]byte)
		val3 := row[2].([]byte)
		val4 := row[3].([]byte)

		var err2 error
		var err3 error
		var err4 error
		var err5 error

<<<<<<< HEAD
    disconnect(db)
=======
		gpsDeviceInfo.Id, err2 = strconv.ParseInt(string(val1[:]), 10, 64)
		gpsDeviceInfo.Name = string(val2[:])
		gpsDeviceInfo.CustomerId, err3 = strconv.ParseInt(string(val3[:]), 10, 64)
		gpsDeviceInfo.IsStolen, err4 = strconv.ParseInt(string(val4[:]), 10, 64)

		err2 = err2
		err3 = err3
		err4 = err4
		err5 = err5

	}
>>>>>>> FETCH_HEAD

	out = fmt.Sprint(gpsDeviceInfo.Id, gpsDeviceInfo.Name, gpsDeviceInfo.CustomerId, gpsDeviceInfo.IsStolen)

	checkError(db.Close())

	return out
}

func GetIpAddressInfo(id_in string) string {

	out := "initial"

<<<<<<< HEAD
    db := connect()
=======
	ipAddressInfo := new(IpAddress)

	user := "root"
	pass := ""
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"
>>>>>>> FETCH_HEAD

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

<<<<<<< HEAD
      ipAddressInfo.ListId, err2 = strconv.ParseInt(string(val1[:]), 10, 64)
      ipAddressInfo.IpAddress = string(val2[:])
      ipAddressInfo.Id, err3 = strconv.ParseInt(string(val3[:]), 10, 64)
  
    err2 = err2
    err3 = err3
    err4 = err4
    
    }
=======
		var err2 error
		var err3 error
		var err4 error
>>>>>>> FETCH_HEAD

		ipAddressInfo.ListId, err2 = strconv.ParseInt(string(val1[:]), 10, 64)
		ipAddressInfo.IpAddress = string(val2[:])
		ipAddressInfo.Id, err3 = strconv.ParseInt(string(val3[:]), 10, 64)

<<<<<<< HEAD
    disconnect(db)
=======
		err2 = err2
		err3 = err3
		err4 = err4

	}

	out = fmt.Sprint(ipAddressInfo.ListId, ipAddressInfo.IpAddress, ipAddressInfo.Id)

	checkError(db.Close())
>>>>>>> FETCH_HEAD

	return out
}

func GetIpListInfo(id_in string) string {

	out := "initial"

<<<<<<< HEAD
    db := connect()
=======
	ipListInfo := new(IpList)

	user := "root"
	pass := ""
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"
>>>>>>> FETCH_HEAD

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

<<<<<<< HEAD
      ipListInfo.Id, err2 = strconv.ParseInt(string(val1[:]), 10, 64)
      ipListInfo.DeviceId, err3 = strconv.ParseInt(string(val2[:]), 10, 64)
      ipListInfo.Timestamp = string(val3[:])
      
    err2 = err2
    err3 = err3
    err4 = err4
    }
=======
		var err2 error
		var err3 error
		var err4 error
>>>>>>> FETCH_HEAD

		ipListInfo.Id, err2 = strconv.ParseInt(string(val1[:]), 10, 64)
		ipListInfo.DeviceId, err3 = strconv.ParseInt(string(val2[:]), 10, 64)
		ipListInfo.Timestamp = string(val3[:])

<<<<<<< HEAD
    disconnect(db)
=======
		err2 = err2
		err3 = err3
		err4 = err4
	}

	out = fmt.Sprint(ipListInfo.Id, ipListInfo.DeviceId, ipListInfo.Timestamp)

	checkError(db.Close())
>>>>>>> FETCH_HEAD

	return out
}

func GetKeyLogsInfo(id_in string) string {

	out := "initial"

<<<<<<< HEAD
    db := connect()
=======
	keyLogsInfo := new(KeyLogs)

	user := "root"
	pass := ""
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"
>>>>>>> FETCH_HEAD

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

<<<<<<< HEAD
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
=======
	for _, row := range rows {
		for _, col := range row {
			if col == nil {
				// col has NULL value
			} else {
				// Do something with text in col (type []byte)
			}
		}
>>>>>>> FETCH_HEAD

		val1 := row[0].([]byte)
		val2 := row[1].([]byte)
		val3 := row[2].([]byte)

<<<<<<< HEAD
    disconnect(db)
=======
		var err2 error
		var err3 error
		var err4 error

		keyLogsInfo.DeviceId, err2 = strconv.ParseInt(string(val1[:]), 10, 64)
		keyLogsInfo.Timestamp = string(val2[:])
		keyLogsInfo.Data = string(val3[:])

		err2 = err2
		err3 = err3
		err4 = err4
>>>>>>> FETCH_HEAD

	}

	out = fmt.Sprint(keyLogsInfo.DeviceId, keyLogsInfo.Timestamp, keyLogsInfo.Data)

	checkError(db.Close())

	return out
}

func GetLaptopDeviceInfo(id_in string) string {

	out := "initial"

<<<<<<< HEAD
    db := connect()
=======
	laptopDeviceInfo := new(LaptopDevice)

	user := "root"
	pass := ""
	dbname := "trackerdb"
	proto := "tcp"
	addr := "127.0.0.1:3306"
>>>>>>> FETCH_HEAD

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

<<<<<<< HEAD
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
=======
	for _, row := range rows {
		for _, col := range row {
			if col == nil {
				// col has NULL value
			} else {
				// Do something with text in col (type []byte)
			}
		}
>>>>>>> FETCH_HEAD

		val1 := row[0].([]byte)
		val2 := row[1].([]byte)
		val3 := row[2].([]byte)
		val4 := row[3].([]byte)
		val5 := row[4].([]byte)

<<<<<<< HEAD
    disconnect(db)
=======
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
>>>>>>> FETCH_HEAD

	out = fmt.Sprint(laptopDeviceInfo.Id, laptopDeviceInfo.DeviceName, laptopDeviceInfo.CustomerId, laptopDeviceInfo.MacAddress, laptopDeviceInfo.IsStolen)

	checkError(db.Close())

	return out
}

<<<<<<< HEAD
func main() {

  //fmt.Println(GetAccountInfo("12"))

  //fmt.Println(GetCoordinatesInfo("1"))

  //fmt.Println(GetCustomerInfo("15"))

  //fmt.Println(GetGpsDeviceInfo("15"))

  //fmt.Println(GetIpAddressInfo("1"))

  //fmt.Println(GetIpListInfo("1"))

  //fmt.Println(GetKeyLogsInfo("1"))

  //fmt.Println(GetLaptopDeviceInfo("1"))

  //fmt.Println(GetUserDevices("15"))

  //fmt.Println(GetUserDevices("sadfk"))

  //SignUp("steven", "whaley", "steven@facebook.gov", "911", "password1")
=======
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
>>>>>>> FETCH_HEAD

//fmt.Println(VerifyAccountInfo("wrongusernameexample", "369d841cdf0dd150a680931769e868d9e487452f"))

//fmt.Println(VerifyAccountInfo("leo@auburn.edu", "wrongpasswordexample"))
//}
