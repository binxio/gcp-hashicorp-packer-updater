FROM golang:1.14

WORKDIR /app
ADD	. /app
RUN go mod download
RUN	go build -o gcp-hashicorp-packer-updater -ldflags '-extldflags "-static"' .

FROM gcr.io/cloud-builders/gcloud:latest

COPY --from=0 /app/gcp-hashicorp-packer-updater /usr/local/bin/

ENTRYPOINT 	["/usr/local/bin/gcp-hashicorp-packer-updater"]
