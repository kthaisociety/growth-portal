# 0005. No personal data leaves the API

Status: Accepted

## Convention
`/api` responses contain aggregates only: no names, emails, Luma IDs or free-text answers of individual
members. `internal/lumaservice` doesn't even decode names and emails.

## Why
The app has no sign-in yet (ADR 0003), so anything the API returns is effectively public.

## How to apply
A new endpoint returns counts or other aggregates. Revisit this convention, with a new PCR, once sign-in
exists.
