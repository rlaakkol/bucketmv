bucketmv
========

Simple service for syncing files from an S3 bucket with deletion.

## Usage
1. Set up your AWS credentials and region in a standard way
2. Set the MV_BUCKET env variable with the bucket name
3. OPTIONAL: Set MV_PATH to point to the directory where files are to be downloaded (default is cwd)
4. Run `bucketmv`

## Building for windows
`GOOS=windows GOARCH=386 go build -o bucketmv.exe`
