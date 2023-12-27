package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	os.Setenv("AWS_ACCESS_KEY_ID", "foo")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "bar")

	file := "UmV0ZWtCcmFuZCxSZXRla0VBTixSZXRla0RFUFQsUmV0ZWtDTEFTUyxSZXRla1NVQkNMQVNTLFByb2R1Y3RNb2RlbERlc2MKMjY2MTYsMjEwMDAwMTMxMDcxNCw2ODAyLDE1LDIsQ2FycmVnYW1lbnRvIFByZVBhZ28gTk9TCjI2NjE2LDIxMDAwMDEzMTA4NzUsNjgwMiwxNSwyLENhcnJlZ2FtZW50byBQcmVQYWdvIE5PUwoyNjYxNiwyMTAwMDAxMzEwOTA1LDY4MDIsMTUsMixDYXJyZWdhbWVudG8gUHJlUGFnbyBOT1MKMjY2MTYsMjEwMDAwMTMxMDkyOSw2ODAyLDE1LDIsQ2FycmVnYW1lbnRvIFByZVBhZ28gTk9T"

	bucketName := "my-bucket-name"

	sessao, err := session.NewSession(
		aws.NewConfig().WithEndpoint("http://127.0.0.1:4566").WithRegion("eu-west-3").WithS3ForcePathStyle(true),
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	// req, resp := s3.New(sessao).ListBucketsRequest(&s3.ListBucketsInput{})

	// err = req.Send()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// if resp.Buckets == nil {
	// 	_, err = s3.New(sessao).CreateBucket(&s3.CreateBucketInput{
	// 		Bucket: aws.String(bucketName),
	// 	})
	// 	if err != nil {
	// 		log.Fatal(err)
	// 		return
	// 	}

	// 	req, resp = s3.New(sessao).ListBucketsRequest(&s3.ListBucketsInput{})

	// 	err = req.Send()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 		return
	// 	}

	// 	for _, v := range resp.Buckets {
	// 		println(*v.Name)
	// 	}
	// } else {
	// 	for _, v := range resp.Buckets {
	// 		println(*v.Name)
	// 	}
	// }

	// nome do arquivo dentro do bucket
	// se for o mesmo nome substitui, se nao cria outro arquivo dentro do bucket
	key := "my_key3"
	// mensagem := []byte("Mensagem")

	fileDecoded, err := base64.StdEncoding.DecodeString(file)
	if err != nil {
		log.Fatal(err)
		return
	}

	contentLength := aws.Int64(int64(len(fileDecoded)))
	contentType := aws.String(http.DetectContentType(fileDecoded))

	_, err = s3.New(sessao).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucketName),
		Key:                  aws.String(key),
		ACL:                  aws.String("public-read-write"),
		Body:                 bytes.NewReader(fileDecoded),
		ContentLength:        contentLength,
		ContentType:          contentType,
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("length:", *contentLength, " nome do arquivo:", key)

	// download file
	// results, err := s3.New(sessao).GetObject(&s3.GetObjectInput{
	// 	Bucket: aws.String(bucketName),
	// 	Key:    aws.String(key),
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// defer results.Body.Close()

	// bytes, err := io.ReadAll(results.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// println(string(bytes))

	err = DownloadFile(sessao, bucketName, key, "lauraFile.csv")
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("sucesso")

}

func DownloadFile(sessao *session.Session, bucketName string, objectKey string, fileName string) error {

	result, err := s3.New(sessao).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}
	_, err = file.Write(body)
	return err
}
