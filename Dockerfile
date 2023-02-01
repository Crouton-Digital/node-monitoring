##
## Build
##

FROM golang:1.19.0-buster AS build

COPY . /opt/

WORKDIR /opt/backend

RUN git config --global user.email "f6b9g8wdbqehaq@estchange.io"
RUN git config --global user.name "Google Cloud Build"
RUN git config --global url."git@bitbucket.org:am-bitbucket".insteadOf "https://bitbucket.org/am-bitbucket"
RUN go env -w GOPRIVATE=bitbucket.org/am-bitbucket

#COPY ./ssh /root/.ssh

#RUN go mod init node-balancer
RUN go mod tidy
RUN go build -o /main ./app/main.go

##
## Deploy
##

FROM golang:1.19.0 AS deploy

RUN mkdir /app
WORKDIR /app

COPY --from=build /main /app/main
#COPY --from=build /opt/backend/config /app/config

#EXPOSE 8080
#EXPOSE 9090

ENTRYPOINT [ "/app/main" ]