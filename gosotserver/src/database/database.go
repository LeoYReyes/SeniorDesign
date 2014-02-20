// testproj project testproj.go

package main

import (
		"fmt"
		"os"
        "encoding/json"
		"github.com/ziutek/mymysql/mysql"
		//_ "github.com/ziutek/mymysql/native" // Native engine
		 _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
	)

type Account struct {
    field1 int64
    field2 int64
    field3 string
    field4 string
    field5 int64
}
type Coordinates struct {
    field1 int64
    field2 float
    field3 float
    field4 string
    field5 int64   
}
type Customer struct {
    field1 int64
    field2 string
    field3 string
    field4 string
    field5 string
    field6 string
}
type GpsDevice struct {
    field1 int64
    field2 string
    field3 int64
}
type IpAddress struct {
    field1 int64
    field2 string
    field3 int64
}
type IpList struct {
    field1 int64
    field2 int64
    field3 string
}
type KeyLogs struct {
    field1 int64
    field2 string
    field3 string
}
type LaptopDevice struct {
    field1 int64
    field2 string
    field3 int64
    field4 string
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

func getAccountInfo() {

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

        val1 := row[0].([]byte)
        val2 := row[1].([]byte)
        val3 := row[2].([]byte)
        val4 := row[3].([]byte)
        //val5 := row[4].([]byte)

        os.Stdout.Write(val1)
        os.Stdout.Write(val2)
        os.Stdout.Write(val3)
        os.Stdout.Write(val4)
       // os.Stdout.Write(val5)


        a := Account{val1, val2, 1294706395881547000}
        b, err := json.Marshal(m)
    }
}

func getCoordinatesInfo() {

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

    rows, res, err := db.Query("select * from coordinates")
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

        os.Stdout.Write(val1)
        os.Stdout.Write(val2)
        os.Stdout.Write(val3)
        os.Stdout.Write(val4)
        os.Stdout.Write(val5)
    }
}

func getCustomerInfo() {

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
    
    rows, res, err := db.Query("SELECT * FROM customer")
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

        os.Stdout.Write(val1)
        os.Stdout.Write(val2)
        os.Stdout.Write(val3)
        os.Stdout.Write(val4)
        os.Stdout.Write(val5)
        os.Stdout.Write(val6)

    }
}

func getGpsDeviceInfo() {

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

    rows, res, err := db.Query("select * from gpsDevice")
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
        
        os.Stdout.Write(val1)
        os.Stdout.Write(val2)
        os.Stdout.Write(val3)
    }
}

func getIpAddressInfo() {

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

    rows, res, err := db.Query("select * from ipAddress")
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

        os.Stdout.Write(val1)
        os.Stdout.Write(val2)
        os.Stdout.Write(val3)
    }
}

func getIpListInfo() {

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

    rows, res, err := db.Query("select * from ipList")
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

        os.Stdout.Write(val1)
        os.Stdout.Write(val2)
        os.Stdout.Write(val3)
    }
}

func getKeyLogsInfo() {

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

    rows, res, err := db.Query("select * from keyLogs")
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

        os.Stdout.Write(val1)
        os.Stdout.Write(val2)
        os.Stdout.Write(val3)
    }
}

func getLaptopDeviceInfo() {

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

    rows, res, err := db.Query("select * from laptopDevice")
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

        os.Stdout.Write(val1)
        os.Stdout.Write(val2)
        os.Stdout.Write(val3)
        os.Stdout.Write(val4)
    }
}

func disconnect() {

}

func main() {

   // getAccountInfo()
    //getCoordinatesInfo()
    getCustomerInfo()
    //getGpsDeviceInfo()
    //getIpAddressInfo()
    //getIpListInfo()
    //getKeyLogsInfo()
    //getLaptopDeviceInfo()
}
