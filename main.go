package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var wasabi *Wasabi
var mongodb *MongoDB
var archive *Archive
var sqlite *SQLite

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env")
		return
	}

	mongodb = &MongoDB{}
	archive = &Archive{}
	wasabi = &Wasabi{
		Region:          os.Getenv("WASABI_REGION"),
		AccessKeyID:     os.Getenv("WASABI_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("WASABI_SECRET_ACCESS_KEY"),
		Endpoint:        os.Getenv("WASABI_URL"),
		Bucket:          os.Getenv("WASABI_BUCKET"),
	}

	sqlite = &SQLite{}
	err = sqlite.connect()
	if err != nil {
		log.Println("Could not connect to SQLite")
		return
	}

	takeBackup()
}

func takeBackup() {
	backupFileName := fmt.Sprintf("%s%s%s", "backup-", time.Now().String(), ".zip")
	path, err := mongodb.dump()
	if err != nil {
		log.Println("Failed to dump database")
		return
	}
	defer cleanUp(path)

	err = archive.compress(path, backupFileName)
	if err != nil {
		log.Println("Failed to compress")
		return
	}
	defer cleanUp(backupFileName)

	_, err = wasabi.upload(backupFileName)
	if err != nil {
		log.Println("Could not upload file")
		return
	}

	err = sqlite.record(backupFileName)
	if err != nil {
		log.Println("Could not store backup record")
		return
	}
}

func cleanUp(items ...string) error {
	log.Printf("Cleaning up files")
	var err error
	for _, item := range items {
		err = os.RemoveAll(item)
	}

	return err
}
