export GIN_MODE=debug
export SESSION_SECRET=secret
export GAMMA_CLIENT_SECRET=secret
export GAMMA_CLIENT_ID=id
export GAMMA_URL=http://localhost:8081
export GAMMA_REDIRECT_URL=http://localhost:3001/api/authenticate
export COOKIE_DOMAIN=localhost

go run goldapps/web/main.go
