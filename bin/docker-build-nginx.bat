cd "%~dp0/../nginx"

docker build -t ericflo/siacdn-nginx:latest .
docker push ericflo/siacdn-nginx:latest