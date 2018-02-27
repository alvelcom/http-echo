FROM golang:1.10 as build
ADD main.go .
RUN CGO_ENABLED=0 go build -o /http-echo main.go

FROM alpine:3.7
RUN apk --no-cache add curl bash bind-tools wget openssh-client sysstat
CMD "/http-echo"
COPY --from=build /http-echo /http-echo

