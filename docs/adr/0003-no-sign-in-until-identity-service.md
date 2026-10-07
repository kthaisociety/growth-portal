# 0003. No sign-in until the identity service exists

Status: Accepted

## Context
This is a proof of concept. Sign-in for KTHAIS apps is planned as the shared identity service
(`infrastructure/docs/identity-service.md`), which isn't built yet.

## Decision
No sign-in for now. To make that safe, the API serves aggregates only, never individual members
(PCR 0005).

## Consequences
Anyone who can reach the app can see the aggregates. When the identity service exists, the app becomes
one of its OIDC clients, and this ADR is superseded.
