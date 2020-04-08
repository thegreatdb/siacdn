@echo off
kubectl exec -it siacdn-uploader-0 -c sia -- bash -c "echo '------------' && echo 'UPLOADER-0:' && ls -lhtra /tmp/*.zip ; ls -lhtra /etc/sia/**/*.db ; cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null ; echo '------------' && echo ''"
kubectl logs --tail=5 siacdn-uploader-0 -c sia
kubectl exec -it siacdn-uploader-1 -c sia -- bash -c "echo && echo '------------' && echo 'UPLOADER-1:' && ls -lhtra /tmp/*.zip ; ls -lhtra /etc/sia/**/*.db ; cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null ; echo '------------' && echo ''"
kubectl logs --tail=5 siacdn-uploader-1 -c sia
kubectl exec -it siacdn-uploader-2 -c sia -- bash -c "echo && echo '------------' && echo 'UPLOADER-2:' && ls -lhtra /tmp/*.zip ; ls -lhtra /etc/sia/**/*.db ; cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null ; echo '------------' && echo ''"
kubectl logs --tail=5 siacdn-uploader-2 -c sia
kubectl exec -it siacdn-uploader-3 -c sia -- bash -c "echo && echo '------------' && echo 'UPLOADER-3:' && ls -lhtra /tmp/*.zip ; ls -lhtra /etc/sia/**/*.db ; cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null ; echo '------------' && echo ''"
kubectl logs --tail=5 siacdn-uploader-3 -c sia
kubectl exec -it siacdn-uploader-4 -c sia -- bash -c "echo && echo '------------' && echo 'UPLOADER-4:' && ls -lhtra /tmp/*.zip ; ls -lhtra /etc/sia/**/*.db ; cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null ; echo '------------' && echo ''"
kubectl logs --tail=5 siacdn-uploader-4 -c sia
kubectl exec -it siacdn-viewer-0 -c sia -- bash -c "echo && echo '------------' && echo 'VIEWER-0:' && ls -lhtra /tmp/*.zip ; ls -lhtra /etc/sia/**/*.db ; cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null ; echo '------------' && echo ''"
kubectl logs --tail=5 siacdn-viewer-0 -c sia
kubectl exec -it siacdn-viewer-1 -c sia -- bash -c "echo && echo '------------' && echo 'VIEWER-1:' && ls -lhtra /tmp/*.zip ; ls -lhtra /etc/sia/**/*.db ; cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null ; echo '------------' && echo ''"
kubectl logs --tail=5 siacdn-viewer-1 -c sia
kubectl exec -it siacdn-viewer-2 -c sia -- bash -c "echo && echo '------------' && echo 'VIEWER-2:' && ls -lhtra /tmp/*.zip ; ls -lhtra /etc/sia/**/*.db ; cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null ; echo '------------' && echo ''"
kubectl logs --tail=5 siacdn-viewer-2 -c sia