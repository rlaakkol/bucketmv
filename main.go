package main

import (
    "os"
    "path"
    "fmt"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
)

func main() {
    bucket_name := os.Getenv("MV_BUCKET")
    dl_path, _ := os.Getwd()
    if env_path, ok := os.LookupEnv("MV_PATH"); ok {
        dl_path = env_path
    }
    sess := session.New()
    svc := s3.New(sess)
    downloader := s3manager.NewDownloader(sess)
    input := &s3.ListObjectsInput{
        Bucket:  aws.String(bucket_name),
    }

    result, err := svc.ListObjects(input)
    if err != nil {
        exitErrorf("Unable to list objects %q, %v", err)
    }

    for _, item := range result.Contents {
        key := *item.Key
        file, err := os.Create(path.Join(dl_path, key))
        if err != nil {
            exitErrorf("Unable to open file %q, %v", err)
        }
        get_input := &s3.GetObjectInput{
            Bucket: aws.String(bucket_name),
            Key: aws.String(key),
        }
        fmt.Println("Copying file", key, "into", dl_path)
        numBytes, err := downloader.Download(file, get_input)
        if err != nil {
            exitErrorf("Unable to download item %q, %v", item, err)
        }
        fmt.Println("Success! Bytes downloaded: ", numBytes)
        delete_input := &s3.DeleteObjectInput{
            Bucket: aws.String(bucket_name),
            Key: aws.String(key),
        }
        fmt.Println("Deleting object from bucket...")
        _, err = svc.DeleteObject(delete_input)
        if err != nil {
            exitErrorf("Unable to delete object %q from bucket %q, %v", key, bucket_name, err)
        }
        check_input := &s3.HeadObjectInput{
            Bucket: aws.String(bucket_name),
            Key: aws.String(key),
        }
        err = svc.WaitUntilObjectNotExists(check_input)
        if err != nil {
            exitErrorf("Error occurred while waiting for object %q to be deleted, %v", key)
        }
        fmt.Println("Deleting succeeded!")
    }
}

func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}
