# go-auth-service

jwt authentication microservice

## Description

production ready auth service with gin server and gorm orm

## Endpoints
- POST `/auth/user/signup`
- POST `/auth/user/signin`
- POST `/auth/admin/sigin`
- POST `/auth/refresh`
- POST `/auth/email-verify/request`
- POST `/auth/email-verify/confirm`
- GET, POST `/users`
- GET, PATCH, DELETE `/users/:id`

## Authentication
- [x] local jwt
- [ ] google oauth2.0
- [ ] facebook oauth2.0

### Credits

- [go-blueprint](https://github.com/Melkeydev/go-blueprint)
- [simplebank](https://github.dev/techschool/simplebank)
- [seeder] (https://github.dev/nayonacademy/gin-seed)
