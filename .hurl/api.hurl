GET {{base_api_url}}/api/v1/healthz/

HTTP 200
[Asserts]
jsonpath "$.code" == 200
jsonpath "$.reason" == "success"

# Add endpoint
POST {{base_api_url}}/api/v1/endpoints/
{
  "name": "test001",
  "runtime": "go"
}

HTTP 200
[Asserts]
jsonpath "$.code" == 200
jsonpath "$.reason" == "success"
jsonpath "$.data.id" != ""
