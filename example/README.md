
```
AWS_PRROFILE=yourprofile
S3_BUCKET=yours3bucket

make deps build
sam package --template-file template.yaml --s3-bucket ${S3_BUCKET} --output-template-file packaged.yaml
aws cloudformation deploy --template-file ./packaged.yaml --stack-name ExampleLogone --capabilities CAPABILITY_IAM
aws lambda invoke --function-name example-logone-lambda-go --payload '{"event": "example"}' --log-type Tail /dev/null
```
