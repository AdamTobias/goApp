# FROM golang:1.6

# ADD . /go/src/github.com/adamtobias/goApp/auth
# RUN go get golang.org/x/crypto/bcrypt
# RUN go get github.com/dgrijalva/jwt-go 
# RUN go install github.com/adamtobias/goApp/auth

# ENTRYPOINT env=docker auth

FROM golang:1.6-onbuild

EXPOSE 8080