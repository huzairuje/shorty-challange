####################################
# STEP 1 build executable binary
####################################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
RUN apk add ca-certificates
WORKDIR $GOPATH/src/test_amartha_muhammad_huzair

#copy all the content to container
COPY . .

##Fix go mod cant download without using proxy
ENV GOPROXY="https://goproxy.cn,direct"

# Build the binary
RUN export CGO_ENABLED=0 && go build -o /go/bin/tiny_url

#change the permission on binary
RUN chmod +x /go/bin/tiny_url

##############################################
# STEP 2 build a small image using alpine:3.14
##############################################
FROM scratch

# Copy our static executable.
COPY --from=builder /go/bin/tiny_url ./tiny_url

# Run the entrypoints.
ENTRYPOINT [ "./tiny_url" ]