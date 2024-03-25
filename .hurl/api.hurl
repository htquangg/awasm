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
POST {{base_api_url}}/api/v1/users/auth/sign-up
{
  "email": "htquangg@gmail.com",
  "ott": "123456"
}

HTTP 200
[Asserts]
jsonpath "$.code" == 200
jsonpath "$.reason" == "base.success"
jsonpath "$.data.token" != ""
