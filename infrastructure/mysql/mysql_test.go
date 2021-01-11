package mysql

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func TestMysql(t *testing.T) {
	db, err := xorm.NewEngine("mysql", "root:123@tcp(127.0.0.1:3306)/test")
	if err != nil {
		t.Fatal(err)
	}
	_ = db
	if rows, err := db.Query("show tables"); err != nil {
		t.Error(err)
	} else {
		t.Log(string(rows[0]["Tables_in_test"]))
	}
}
