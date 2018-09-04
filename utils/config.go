package utils

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Config struct {
	Region              string
	Bucket              string
	Prefix              string
	BackupBucket        string
	BackupPrefix        string
	BackupRetentionDays int
	ArchiveFolderName   string
	S3                  *s3.S3
	Uploader            *s3manager.Uploader
}

