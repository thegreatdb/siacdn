@echo off
kubectl exec -it siacdn-uploader-0 -c sia -- bash -c "echo 'UPLOADER-0:' && cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null && echo ''"
kubectl exec -it siacdn-uploader-1 -c sia -- bash -c "echo 'UPLOADER-1:' && cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null && echo ''"
kubectl exec -it siacdn-uploader-2 -c sia -- bash -c "echo 'UPLOADER-2:' && cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null && echo ''"
kubectl exec -it siacdn-uploader-3 -c sia -- bash -c "echo 'UPLOADER-3:' && cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null && echo ''"
kubectl exec -it siacdn-uploader-4 -c sia -- bash -c "echo 'UPLOADER-4:' && cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null && echo ''"
kubectl exec -it siacdn-viewer-0 -c sia -- bash -c "echo 'VIEWER-0:' && cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null && echo ''"
kubectl exec -it siacdn-viewer-1 -c sia -- bash -c "echo 'VIEWER-1:' && cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null && echo ''"
kubectl exec -it siacdn-viewer-2 -c sia -- bash -c "echo 'VIEWER-2:' && cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null && echo ''"
