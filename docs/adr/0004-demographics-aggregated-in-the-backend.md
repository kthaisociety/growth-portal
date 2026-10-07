# 0004. Demographics are aggregated in the backend

Status: Accepted

## Context
The dashboards show how the Community tier's members answered their registration questions (graduation
year, areas of interest, …).

## Decision
`internal/demographics` counts answers per question, from members fetched through luma-service, and
`GET /api/community/demographics` returns only those counts. Only choice questions (`select`) are counted;
free-text answers are left out.

## Consequences
Personal data never reaches the browser. New breakdowns (e.g. members joined over time) are added as
functions in `internal/demographics`, each with a plain-value test. Members are fetched on every request
for now; cache when that becomes slow.
