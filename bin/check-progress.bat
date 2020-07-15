@echo off
kubectl exec -it deployment/siacdn-uploader-0 -c sia -- bash -c "echo '------------' && echo 'UPLOADER-0:' && ls -lhtra /tmp/*.zip ; ls -lhtra /etc/sia/**/*.db ; cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null ; echo '------------' && echo ''"
kubectl logs --tail=5 deployment/siacdn-uploader-0 -c sia
kubectl exec -it deployment/siacdn-uploader-1 -c sia -- bash -c "echo && echo '------------' && echo 'UPLOADER-1:' && ls -lhtra /tmp/*.zip ; ls -lhtra /etc/sia/**/*.db ; cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null ; echo '------------' && echo ''"
kubectl logs --tail=5 deployment/siacdn-uploader-1 -c sia
kubectl exec -it deployment/siacdn-uploader-2 -c sia -- bash -c "echo && echo '------------' && echo 'UPLOADER-2:' && ls -lhtra /tmp/*.zip ; ls -lhtra /etc/sia/**/*.db ; cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null ; echo '------------' && echo ''"
kubectl logs --tail=5 deployment/siacdn-uploader-2 -c sia
kubectl exec -it deployment/siacdn-uploader-3 -c sia -- bash -c "echo && echo '------------' && echo 'UPLOADER-3:' && ls -lhtra /tmp/*.zip ; ls -lhtra /etc/sia/**/*.db ; cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null ; echo '------------' && echo ''"
kubectl logs --tail=5 deployment/siacdn-uploader-3 -c sia
kubectl exec -it deployment/siacdn-viewer-0 -c sia -- bash -c "echo && echo '------------' && echo 'VIEWER-0:' && ls -lhtra /tmp/*.zip ; ls -lhtra /etc/sia/**/*.db ; cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null ; echo '------------' && echo ''"
kubectl logs --tail=5 deployment/siacdn-viewer-0 -c sia
kubectl exec -it deployment/siacdn-viewer-1 -c sia -- bash -c "echo && echo '------------' && echo 'VIEWER-1:' && ls -lhtra /tmp/*.zip ; ls -lhtra /etc/sia/**/*.db ; cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null ; echo '------------' && echo ''"
kubectl logs --tail=5 deployment/siacdn-viewer-1 -c sia
kubectl exec -it deployment/siacdn-viewer-2 -c sia -- bash -c "echo && echo '------------' && echo 'VIEWER-2:' && ls -lhtra /tmp/*.zip ; ls -lhtra /etc/sia/**/*.db ; cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null ; echo '------------' && echo ''"
kubectl logs --tail=5 deployment/siacdn-viewer-2 -c sia
kubectl exec -it deployment/siacdn-viewer-3 -c sia -- bash -c "echo && echo '------------' && echo 'VIEWER-3:' && ls -lhtra /tmp/*.zip ; ls -lhtra /etc/sia/**/*.db ; cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null ; echo '------------' && echo ''"
kubectl logs --tail=5 deployment/siacdn-viewer-3 -c sia