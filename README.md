# Golang Auth Starter
JWT Based Password less Authentication
## Features
- JWT with private key signature 
- JTI based session with Redis
- Password less otp authentication
- Rate limit
- Emailer with AWS SES client
- Database migration with Atlas
- Health endpoints
## Directory Structure
```bash
.
├── cmd
│   ├── root.go
│   ├── schema.go
│   └── srv_start.go
├── config
│   └── srv_config.go
├── internal
│   ├── api
│   │   ├── auth.go
│   │   ├── health.go
│   │   ├── router.go
│   │   └── user.go
│   ├── comm
│   │   ├── aws_ses.go
│   │   └── email.go
│   ├── core
│   │   ├── app.go
│   │   ├── aws_session.go
│   │   ├── identity.go
│   │   ├── rate_limit.go
│   │   └── validate.go
│   ├── enum
│   │   ├── gender.go
│   │   └── identity_provider.go
│   ├── errors
│   │   └── http_error.go
│   ├── health
│   │   └── health.go
│   ├── middleware
│   │   └── auth_interceptor.go
│   ├── misc
│   │   ├── crypto.go
│   │   ├── jwt_helper.go
│   │   └── otp.go
│   ├── model
│   │   └── user.go
│   ├── oidc
│   │   ├── apple.go
│   │   ├── google.go
│   │   └── oidc.go
│   ├── repo
│   │   ├── access_token.go
│   │   ├── auth.go
│   │   ├── repo.go
│   │   └── user.go
│   ├── storage
│   │   ├── cache_store.go
│   │   └── db_store.go
│   ├── templates
│   │   ├── email_update_otp_template.html
│   │   └── sign_in_otp_template.html
│   └── uid
│       ├── id.go
│       ├── id_generator.go
│       └── id_serializer.go
├── atlas.hcl
├── go.mod
├── go.sum
├── main.go
├── migrations
│   ├── 20240225050014.sql
│   └── atlas.sum
```