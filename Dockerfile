FROM golang:1.9.3-stretch AS build-env

RUN go-wrapper download github.com/golang/dep/cmd/dep && go-wrapper install github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/larsp/co2monitor

# only copy Gopkg.toml and Gopkg.lock for now so dep ensure run keeps cached till those files are changed
COPY Gopkg.* ./
RUN dep ensure -v --vendor-only

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -v -ldflags "-w -s"  -a -installsuffix cgo -o /out/co2monitor

FROM scratch
COPY --from=build-env /out/co2monitor /
ENTRYPOINT [ "./co2monitor" ]