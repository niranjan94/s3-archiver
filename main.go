package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/uniplaces/carbon"
	"os"
	"strconv"
	"fmt"
	"github.com/niranjan94/s3-archiver/utils"
)

var config utils.Config

func init() {

	awsSession := session.Must(session.NewSession())

	config = utils.Config{
		Bucket:              os.Getenv("BUCKET"),
		Prefix:              os.Getenv("PREFIX"),
		BackupPrefix:        os.Getenv("BACKUP_PREFIX"),
		BackupBucket:        os.Getenv("BACKUP_BUCKET"),
		BackupRetentionDays: 7,
		ArchiveFolderName:   carbon.Now().Format("2006-01-02"),
		S3:                  s3.New(awsSession),
		Uploader:            s3manager.NewUploader(awsSession),
	}

	backupRetentionDaysString := os.Getenv("BACKUP_RETENTION_DAYS")

	if len(backupRetentionDaysString) != 0 {
		backupRetentionDays, err := strconv.Atoi(backupRetentionDaysString)
		if err == nil {
			config.BackupRetentionDays = backupRetentionDays
		} else {
			fmt.Printf("Recived invalid number %s for backup retention days\n", backupRetentionDaysString)
		}
	}
}

func main() {
	utils.ListAndUpload(&config)
	utils.PerformCleanup(&config)
}
