FROM golang:1.9 AS build-env

WORKDIR /go/src/github.com/larsp/co2monitor

COPY . ./

RUN go-wrapper download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -v -ldflags "-w -s"  -a -installsuffix cgo -o /out/co2monitor

FROM scratch
COPY --from=build-env /out/co2monitor /
ENTRYPOINT [ "./co2monitor" ]