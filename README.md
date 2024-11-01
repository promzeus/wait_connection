# wait_connection

Outputs the number of active connections on a specified port to stdout. By default, it monitors port `8080`. When the number of connections decreases to `0` or `1`, the application exits with `exit 0`.

To change the port, use the `--port` flag, for example: `--port 80`.

![Demo](https://github.com/promzeus/wait_connection/blob/main/connection_wait.gif)

### Deployment.yaml:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: {{ .Release.Name }}
  name: {{ .Release.Name }}
spec:
  .....
        {{- if .Values.probes.enabled }}
        lifecycle:
          postStart:
            exec:
              command: [ "/bin/bash","-c", "touch /tmp/health" ]
          preStop:
            exec:
              command: [ "/bin/bash","-c", "rm -rf /tmp/health && /usr/local/bin/connection_wait --port=10000 && kill -SIGTERM 1"]
        readinessProbe:
          initialDelaySeconds: {{ .Values.probes.readinessProbe.initialDelaySeconds }}
          periodSeconds: {{ .Values.probes.readinessProbe.periodSeconds }}
          timeoutSeconds: {{ .Values.probes.readinessProbe.timeoutSeconds }}
          successThreshold: {{ .Values.probes.readinessProbe.successThreshold }}
          failureThreshold: {{ .Values.probes.readinessProbe.failureThreshold }}
          exec:
            command: ["test","-e","/tmp/health"]
        {{- end }}
    ....