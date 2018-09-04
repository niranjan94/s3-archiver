package utils

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type callback func(objects []*s3.Object)


func ListObjectsInBatches(config *Config, sendResult callback, continuationTokenArgs ...*string )  {

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(config.Bucket),
		Prefix: aws.String(config.Prefix),
		MaxKeys: aws.Int64(100),
	}

	if len(continuationTokenArgs) > 0 {
		input.ContinuationToken = continuationTokenArgs[0]
	}

	if result, err := config.S3.ListObjectsV2(input); err != nil {
		panic(err)
	} else {
		sendResult(result.Contents)
		if result.NextContinuationToken != nil && len(*result.NextContinuationToken) > 0 {
			ListObjectsInBatches(config, sendResult, result.NextContinuationToken)
		}
	}

}
