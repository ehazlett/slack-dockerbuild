# Dockerbuild
Notifies Slack when a Docker Trusted Build is complete.

# Build
`go build`

# Usage
Listens on port 8080.

Add an "Incoming Webhook Integration" in Slack and then add a webhook for this service (`/notify`) in the Docker Trusted Build.

```
./slack-dockerbuild -url <slack-incoming-webhook-url> -channel <slack-channel>
```

