# 0004. Luma data comes only from luma-service

Status: Accepted

## Convention
The app never calls the Luma API and holds no Luma key. All Luma data comes from luma-service, through
`internal/lumaservice`, authenticated with `X-Service-Secret`.

## Why
luma-service is the one place that handles Luma's API key, pagination and rate limit, for every KTHAIS
service.

## How to apply
Need Luma data luma-service doesn't serve yet? Add the endpoint to luma-service first, then a method to
`internal/lumaservice`.
