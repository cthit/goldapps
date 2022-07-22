export GIN_MODE=debug
export SESSION_SECRET=secret
export AUTH_SECRET=secret
export CLIENT_ID=id
export GAMMA_URL=http://localhost:8081
export REDIRECT_GAMMA_URL=http://localhost:8081
export CALLBACK_URL=http://localhost:3001/callback
export COOKIE_DOMAIN=localhost

go run goldapps/web/main.go
