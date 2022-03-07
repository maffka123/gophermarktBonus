# Gophermarkt Bonus App

It is a aservice that allows to recieve Bonuses for orders

## API endpoints

POST /api/user/register — create user

POST /api/user/login — authenticate user

POST /api/user/orders — register new order

GET /api/user/orders — return all orders

GET /api/user/balance — current balance and withdrawals

POST /api/user/balance/withdraw — available bonus points

GET /api/user/balance/withdrawals — all withdrawals

## Testing with cURL

 `curl -d '{"login":"test1","password":"mypass"}' -H "Content-Type: application/json" -X POST http://localhost:8081/api/user/register`

  `curl -d '{"login":"test1","password":"mypass"}' -H "Content-Type: application/json" -X POST http://localhost:8081/api/user/login`

   `curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.uwbhqVZMHjeX9nvVpbw-AHXZ2YAfNToBR1IGjITmxo4" -H "Content-Type: text/plain" -d 18 -X POST http://localhost:8081/api/user/orders`

   `curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.uwbhqVZMHjeX9nvVpbw-AHXZ2YAfNToBR1IGjITmxo4"  -X GET http://localhost:8081/api/user/orders`

`curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.uwbhqVZMHjeX9nvVpbw-AHXZ2YAfNToBR1IGjITmxo4"  -X GET http://localhost:8081/api/user/balance`

`curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.uwbhqVZMHjeX9nvVpbw-AHXZ2YAfNToBR1IGjITmxo4" -H "Content-Type: application/json" -d '{"order":"1234566","sum":0.5}' -X POST http://localhost:8081/api/user/balance/withdraw`

`curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.uwbhqVZMHjeX9nvVpbw-AHXZ2YAfNToBR1IGjITmxo4" -H "Content-Type: application/json" -d '{"order":"1234566","sum":10000}' -X POST http://localhost:8081/api/user/balance/withdraw`

`curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.uwbhqVZMHjeX9nvVpbw-AHXZ2YAfNToBR1IGjITmxo4" -X GET http://localhost:8081/api/user/balance/withdrawals`

 `curl -d '{"login":"test2","password":"mypass"}' -H "Content-Type: application/json" -X POST http://localhost:8081/api/user/register`

   `curl -d '{"login":"test2","password":"mypass"}' -H "Content-Type: application/json" -X POST http://localhost:8081/api/user/login`

  `curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyfQ.XUjieZQLFHd61t9ZjifbQ6c1BGB6ANYD1Xo-aog249U" -H "Content-Type: text/plain" -d 18 -X POST http://localhost:8081/api/user/orders`

   `curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyfQ.XUjieZQLFHd61t9ZjifbQ6c1BGB6ANYD1Xo-aog249U"  -X GET http://localhost:8081/api/user/orders`

   `curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyfQ.XUjieZQLFHd61t9ZjifbQ6c1BGB6ANYD1Xo-aog249U" -X GET http://localhost:8081/api/user/balance/withdrawals`

 `curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.uwbhqVZMHjeX9nvVpbw-AHXZ2YAfNToBR1IGjITmxo4" -H "Content-Type: text/plain" -d 123456758 -X POST http://localhost:8081/api/user/orders`

  `curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.uwbhqVZMHjeX9nvVpbw-AHXZ2YAfNToBR1IGjITmxo4" -H "Content-Type: text/plain" -d 123456766 -X POST http://localhost:8081/api/user/orders`

  `curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.uwbhqVZMHjeX9nvVpbw-AHXZ2YAfNToBR1IGjITmxo4" -H "Content-Type: text/plain" -d 123456774 -X POST http://localhost:8081/api/user/orders`

  `curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.uwbhqVZMHjeX9nvVpbw-AHXZ2YAfNToBR1IGjITmxo4" -H "Content-Type: text/plain" -d 162124202183 -X POST http://localhost:8081/api/user/orders`
