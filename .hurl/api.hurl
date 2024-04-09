GET {{base_api_url}}/api/v1/healthz/

HTTP 200
[Asserts]
jsonpath "$.code" == 200
jsonpath "$.reason" == "base.success"

# Add endpoint
POST {{base_api_url}}/api/v1/endpoints/
{
  "name": "test00122",
  "runtime": "go"
}

HTTP 200
[Asserts]
jsonpath "$.code" == 200
jsonpath "$.reason" == "base.success"
jsonpath "$.data.id" != ""

# Sign Up

## Step 1: Begin email signup process
POST {{base_api_url}}/api/v1/users/auth/email/signup
{
  "email": "htquangg@gmail.com"
}

HTTP 200
[Asserts]
jsonpath "$.code" == 200
jsonpath "$.reason" == "base.success"

## Step 2: Verify OTP
POST {{base_api_url}}/api/v1/users/auth/email/verify
{
  "email": "htquangg@gmail.com",
  "otp": "123456"
}

HTTP 200
[Asserts]
jsonpath "$.code" == 200
jsonpath "$.reason" == "base.success"
[Captures]
access_token: jsonpath "$.data.accessToken"

## Step 3: Setup SRP account
POST {{base_api_url}}/api/v1/users/auth/email/setup-srp
Authorization: Bearer {{access_token}}
{
  "srpUserId": "B4E3FC0E-AC6E-4C3F-A875-7ADCFAE4F30E",
  "srpSalt": "E1D0FE2E-3194-46EA-983C-DB90EC30B423",
  "srpVerifier": "A45411C1-098D-42A1-8434-B27F865CFA95",
  "srpA": "0E2EFFD9-3573-463A-A202-6447C6B5E236"
}

HTTP 200
