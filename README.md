# S3 updater

Used to update contents within S3 bucket.

## Usage

1. Configure AWS SSO profiles in ~/.aws/config using `aws configure sso` or simply by updating `~/.aws/config` like so:
```
[profile your-profile-name]
sso_start_url = https://yourssologinname.awsapps.com/start#/
sso_region = eu-west-1
sso_account_id = your-account-id
sso_role_name = AdministratorAccess
region = eu-west-1
output = json
```
2. Run `go build .` to create a binary in ./s3-updater.
3. Run `sudo cp ./s3-updater /usr/local/bin/s3-updater` to be able to call the binary from anywhere.
4. Call `s3-updater` when you want to update the contents of the S3 bucket.
5. Select profile -> Authenticate (if needed - sso session by default lasts 8 hours.) -> Select bucket -> Modify files within /tmp/s3-uploader -> Confirm to reupload changes to S3 bucket.
6. Done.
