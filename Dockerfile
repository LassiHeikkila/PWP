##
## BUILD STAGE
##
FROM golang:1.17 AS build

WORKDIR /src

# copy dependency info
COPY go.mod ./
COPY go.sum ./
# download dependencies, only re-done when dependencies change
RUN go mod download

# copy everything else
COPY . .

# build binaries
RUN echo "Building taskey" && go build -v -o taskey ./cmd/taskey/
RUN echo "Building taskeyd" && go build -v -o taskeyd ./cmd/taskeyd/
RUN echo "Building taskey-cli" && go build -v -o taskey-cli ./cmd/taskey-cli/

##
## DEPLOY STAGE
##
FROM gcr.io/distroless/base-debian11 AS deploy

COPY --from=build /src/taskey /taskey
COPY --from=build /src/taskeyd /taskeyd
COPY --from=build /src/taskey-cli /taskey-cli

EXPOSE 80

ENTRYPOINT ["/taskey"]
