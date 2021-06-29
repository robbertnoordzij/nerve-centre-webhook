# nerve-centre-webhook
A simple applicatation to push the schedule defined in Nerve Centre to a Slack webhook.

## Usage

```bash
docker build -t nerve-centre-webhook:latest .
docker run nerve-centre-webhook:latest --namespace "<<nerve-centre-namespace>>" --username "<<nerve-centre-username>>" --password "<<nerve-centre-password>>" --webhook "<<slack-webhook-url>>" --channel "<<slack-channel-override>>"
```

## Docker hub

Also available on Docker hub: https://hub.docker.com/r/robbert0001/nerve-centre-webhook
