@echo off
kubectl exec -it deployment/siacdn-uploader-0 -c sia -- bash -c "echo 'UPLOADER-0:' && grep -H blockheight /etc/sia/**/*.json && echo && echo"
kubectl exec -it deployment/siacdn-uploader-1 -c sia -- bash -c "echo 'UPLOADER-1:' && grep -H blockheight /etc/sia/**/*.json && echo && echo"
kubectl exec -it deployment/siacdn-uploader-2 -c sia -- bash -c "echo 'UPLOADER-2:' && grep -H blockheight /etc/sia/**/*.json && echo && echo"
kubectl exec -it deployment/siacdn-uploader-3 -c sia -- bash -c "echo 'UPLOADER-3:' && grep -H blockheight /etc/sia/**/*.json && echo && echo"
kubectl exec -it deployment/siacdn-viewer-0 -c sia -- bash -c "echo 'VIEWER-0:' && grep -H blockheight /etc/sia/**/*.json && echo && echo"
kubectl exec -it deployment/siacdn-viewer-1 -c sia -- bash -c "echo 'VIEWER-1:' && grep -H blockheight /etc/sia/**/*.json && echo && echo"
kubectl exec -it deployment/siacdn-viewer-2 -c sia -- bash -c "echo 'VIEWER-2:' && grep -H blockheight /etc/sia/**/*.json && echo && echo"
kubectl exec -it deployment/siacdn-viewer-3 -c sia -- bash -c "echo 'VIEWER-3:' && grep -H blockheight /etc/sia/**/*.json && echo && echo"