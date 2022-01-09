Deploy docker image to GCloud

```bash
# After commit
docker build -t cowculator99/realtime-stream-server . && docker push cowculator99/realtime-stream-server

#  Then on GCP
docker pull cowculator99/realtime-stream-server
docker tag cowculator99/realtime-stream-server gcr.io/learning-gcp-325508/cowculator99/realtime-stream-server
docker push gcr.io/learning-gcp-325508/cowculator99/realtime-stream-server
```