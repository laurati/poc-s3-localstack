package download

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// import (
// 	"context"
// 	"io"
// 	"log"
// 	"os"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/service/s3"
// )

// // BucketBasics encapsulates the Amazon Simple Storage Service (Amazon S3) actions
// // used in the examples.
// // It contains S3Client, an Amazon S3 service client that is used to perform bucket
// // and object actions.
// type BucketBasics struct {
// 	S3Client *s3.Client
// }

// // DownloadFile gets an object from a bucket and stores it in a local file.
// func (basics BucketBasics) DownloadFile(bucketName string, objectKey string, fileName string) error {

// 	result, err := basics.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
// 		Bucket: aws.String(bucketName),
// 		Key:    aws.String(objectKey),
// 	})
// 	if err != nil {
// 		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
// 		return err
// 	}
// 	defer result.Body.Close()
// 	file, err := os.Create(fileName)
// 	if err != nil {
// 		log.Printf("Couldn't create file %v. Here's why: %v\n", fileName, err)
// 		return err
// 	}
// 	defer file.Close()
// 	body, err := io.ReadAll(result.Body)
// 	if err != nil {
// 		log.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
// 	}
// 	_, err = file.Write(body)
// 	return err
// }

func downloadManager(bucketName string, objectKey string, filename string) error {

	// The session the S3 Downloader will use
	sess := session.Must(session.NewSession())

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	// Create a file to write the S3 Object contents to.
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %q, %v", filename, err)
	}

	// Write the contents of S3 Object to the file
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to download file, %v", err)
	}
	fmt.Printf("file downloaded, %d bytes\n", n)

	return nil
}
