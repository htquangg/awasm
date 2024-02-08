GET {{base_api_url}}/api/v1/healthz/

HTTP 200
[Asserts]
jsonpath "$.code" == 200
jsonpath "$.reason" == "success"
