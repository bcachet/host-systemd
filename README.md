```shell
docker build -f Dockerfile --tag vault:custom .
docker run -v /run/dbus/:/run/dbus/ vault:custom systemctl status
```
