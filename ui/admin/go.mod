module admin

go 1.20

require github.com/golang-jwt/jwt v3.2.2+incompatible

require github.com/gorilla/mux v1.8.0

require (
	github.com/go-sql-driver/mysql v1.5.0
	gopkg.in/yaml.v2 v2.4.0
)

require golang.org/x/crypto v0.11.0 // indirect

require (
	github.com/google/uuid v1.3.0
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/gorilla/sessions v1.2.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/sys v0.10.0 // indirect
)
