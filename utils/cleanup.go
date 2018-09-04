package utils

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"regexp"
	"strings"
	"github.com/uniplaces/carbon"
	"log"
	"github.com/aws/aws-sdk-go/aws"
)

var slashPrefixRegex = regexp.MustCompile(`^/+`)

func PerformCleanup(config *Config)  {

	batch := 1

	ListObjectsInBatches(
		&Config{
			Bucket:   config.BackupBucket,
			Prefix:   config.BackupPrefix,
			S3:       config.S3,
			Uploader: config.Uploader,
		},
		func(objects []*s3.Object) {

			var objectsToDelete []*s3.ObjectIdentifier

			for _, object := range objects {
				folderName := strings.Replace(*object.Key, config.BackupPrefix, "", -1)
				folderName = slashPrefixRegex.ReplaceAllString(folderName, "")
				folderName = strings.Split(folderName, "/")[0]
				date, _ := carbon.Parse("2006-01-02", folderName, "UTC")
				if date.Before(carbon.Now().SubDays(7).StartOfDay().Time) {
					objectsToDelete = append(objectsToDelete, &s3.ObjectIdentifier{
						Key: object.Key,
					})
					log.Printf("Marked %s for deletion.\n", *object.Key)
				}
			}

			if length := len(objectsToDelete); length > 0 {
				output, err := config.S3.DeleteObjectsWithContext(aws.BackgroundContext(), &s3.DeleteObjectsInput{
					Bucket: &config.BackupBucket,
					Delete: &s3.Delete{
						Objects: objectsToDelete,
						Quiet: aws.Bool(false),
					},
				})

				log.Printf("%v\n", output)

				if err != nil {
					panic(err)
				} else {
					log.Printf("%d objects deleted.\n", length)
				}
			} else {
				log.Printf("No objects to cleanup in batch %d.\n", batch)
			}

			batch = batch + 1
		},
	)
}
