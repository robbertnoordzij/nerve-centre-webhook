# nerve-centre-webhook
Simple webhook to Slack to push the schedule defined in Nerve Centre

## Usage

```bash
docker build -t nerve-centre-webhook:latest .
docker run nerve-centre-webhook:latest -username "<<nerve-centre-username>>" -password "<<nerve-centre-password>>" --webhook "<<slack-webhook-url>>"
```
