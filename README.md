# nuki2mqtt

Go program that receives Nuki Bridge webhooks and publishes contents to MQTT.

# gokrazy

If you’re running this program on gokrazy, see
https://gokrazy.org/userguide/package-config/ for how to set command-line flags
to influence the listening address, MQTT broker and MQTT topic.

# Setup

Configure your Nuki Bridge to send webhooks to this process:
https://developer.nuki.io/page/nuki-bridge-http-api-1-12/4#heading--callback

E.g.:

```
curl -v http://Nuki_Bridge_XXYYZZ02:8080/callback/add\?url\=http://10.0.0.54:8319/nuki\&token\=SECRET | jq '.'
```
