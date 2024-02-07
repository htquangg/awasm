GET {{base_api_url}}/api/v1/healthz/

HTTP 200
[Asserts]
jsonpath "$.code" == 1 # 1: success, 2: bad request, 3: not found, 4: unauthorized, 5: forbidden, 6: internal server error
