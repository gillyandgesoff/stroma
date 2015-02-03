export GOOS=linux
go install
cp ~/Development/Code/Go/bin/linux_amd64/stroma ./stroma
docker rmi bengesoff/stroma:latest
docker build -t bengesoff/stroma .
