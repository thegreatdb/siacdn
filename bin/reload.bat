kubectl delete pod -l "appinstance=siacdn-uploader-0"
SLEEP 60
kubectl delete pod -l "appinstance=siacdn-uploader-1"
SLEEP 60
kubectl delete pod -l "appinstance=siacdn-uploader-2"
SLEEP 60
kubectl delete pod -l "appinstance=siacdn-uploader-3"
SLEEP 60
kubectl delete pod -l "appinstance=siacdn-uploader-4"
SLEEP 60
kubectl delete pod -l "appinstance=siacdn-viewer-0"
SLEEP 60
kubectl delete pod -l "appinstance=siacdn-viewer-1"
SLEEP 60
kubectl delete pod -l "appinstance=siacdn-viewer-2"
SLEEP 60
kubectl delete pod -l "appinstance=siacdn-viewer-3"
SLEEP 60
kubectl delete pod -l "appinstance=siacdn-viewer-4"