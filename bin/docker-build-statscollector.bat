cd "%~dp0/../statscollector"

docker build -t ericflo/siacdn-statscollector:latest .
docker push ericflo/siacdn-statscollector:latest