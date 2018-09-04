package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"strings"
	"log"
)

func ListAndUpload(config *Config) {

	ListObjectsInBatches(
		&Config{
			Bucket:   config.Bucket,
			Prefix:   config.Prefix,
			S3:       config.S3,
			Uploader: config.Uploader,
		},
		func(objects []*s3.Object) {

			var objectsToUpload []s3manager.BatchUploadObject

			for _, object := range objects {

				if strings.HasSuffix(*object.Key, "/") {
					continue
				}

				objectData, _ := config.S3.GetObject(&s3.GetObjectInput{
					Bucket: &config.Bucket,
					Key:    object.Key,
				})
				key := strings.Replace(*object.Key, config.Prefix, "", -1)
				if !strings.HasPrefix(key, "/") {
					key = "/" + key
				}

				objectsToUpload = append(objectsToUpload, s3manager.BatchUploadObject{
					Object: &s3manager.UploadInput{
						Bucket: &config.BackupBucket,
						Key:    aws.String(config.BackupPrefix + "/" + config.ArchiveFolderName + key),
						Body:   objectData.Body,
					},
				})

				log.Println("Added file to queue: " + *object.Key)
			}

			log.Printf("Queued %d objects for upload. Starting upload.", len(objectsToUpload))

			iterator := &s3manager.UploadObjectsIterator{Objects: objectsToUpload}
			if err := config.Uploader.UploadWithIterator(aws.BackgroundContext(), iterator); err != nil {
				panic(err)
			} else {
				log.Printf("Upload completed.")
			}

		},
	)

}
