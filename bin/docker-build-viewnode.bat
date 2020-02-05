cd "%~dp0/../viewnode"

docker build -t ericflo/siacdn-viewnode:latest .
docker push ericflo/siacdn-viewnode:latest