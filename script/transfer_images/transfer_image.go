package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // MySQLドライバ
	"github.com/jmoiron/sqlx"
)

type Icon struct {
	ID     int    `db:"id"`
	UserID string `db:"user_id"`
	Image  []byte `db:"image"`
}

// ISUCON_DB_HOST=${ISUCON13_MYSQL_DIALCONFIG_ADDRESS:-127.0.0.1}
// ISUCON_DB_PORT=${ISUCON13_MYSQL_DIALCONFIG_PORT:-3306}
// ISUCON_DB_USER=${ISUCON13_MYSQL_DIALCONFIG_USER:-isucon}
// ISUCON_DB_PASSWORD=${ISUCON13_MYSQL_DIALCONFIG_PASSWORD:-isucon}
// ISUCON_DB_NAME=${ISUCON13_MYSQL_DIALCONFIG_DATABASE:-isupipe}
func main() {
	// db接続
	// dsn := "username:password@tcp(localhost:3306)/dbname?parseTime=true"

	username := os.Getenv("ISUCON13_MYSQL_DIALCONFIG_USER")
	password := os.Getenv("ISUCON13_MYSQL_DIALCONFIG_PASSWORD")
	host := os.Getenv("ISUCON13_MYSQL_DIALCONFIG_ADDRESS")
	port := os.Getenv("ISUCON13_MYSQL_DIALCONFIG_PORT")
	dbname := os.Getenv("ISUCON13_MYSQL_DIALCONFIG_DATABASE")

	// DSNを組み立て
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, dbname)
	fmt.Println(dsn)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}
	defer db.Close()

	// クエリを実行して複数行のデータを取得
	var images []Icon
	query := "SELECT id, user_id, image FROM icons"
	err = db.Select(&images, query)
	if err != nil {
		log.Fatalln("Failed to execute query:", err)
	}

	for _, v := range images {
		// バイナリデータをReaderに変換
		reader := bytes.NewReader(v.Image)

		// 画像データをデコード
		img, _, err := image.Decode(reader)
		if err != nil {
			log.Fatal(err)
		}

		// 出力ファイルを作成
		dir := "../../webapp/img/"
		f, err := os.Create(dir + v.UserID + ".jpg")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		// JPEGエンコードオプション
		opt := jpeg.Options{
			Quality: 90,
		}

		// JPEGとしてファイルに書き出す
		if err := jpeg.Encode(f, img, &opt); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("画像が正常に保存されました。")
}
