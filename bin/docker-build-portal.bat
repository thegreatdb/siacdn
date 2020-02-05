cd "%~dp0/../portal"

docker build -t ericflo/siacdn-portal:latest .
docker push ericflo/siacdn-portal:latest