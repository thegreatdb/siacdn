@echo off
kubectl exec -it siacdn-uploader-0 -c nginx -- bash -c "echo -n 'uploader-0> ' && cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null"
kubectl exec -it siacdn-uploader-1 -c nginx -- bash -c "echo -n 'uploader-1> ' && cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null"
kubectl exec -it siacdn-uploader-2 -c nginx -- bash -c "echo -n 'uploader-2> ' && cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null"
kubectl exec -it siacdn-uploader-3 -c nginx -- bash -c "echo -n 'uploader-3> ' && cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null"
kubectl exec -it siacdn-uploader-4 -c nginx -- bash -c "echo -n 'uploader-4> ' && cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null"
kubectl exec -it siacdn-viewer-0 -c nginx -- bash -c "echo -n 'viewer-1> ' && cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null"
kubectl exec -it siacdn-viewer-1 -c nginx -- bash -c "echo -n 'viewer-2> ' && cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null"
kubectl exec -it siacdn-viewer-2 -c nginx -- bash -c "echo -n 'viewer-3> ' && cat /etc/sia/host/host.json 2>/dev/null | grep blockheight 2>/dev/null"