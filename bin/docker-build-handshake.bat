cd "%~dp0/../handshake"

docker build -t ericflo/siacdn-handshake:latest .
docker push ericflo/siacdn-handshake:latest