cd "%~dp0/../handshake-api"

docker build -t ericflo/siacdn-handshake-api:latest .
docker push ericflo/siacdn-handshake-api:latest