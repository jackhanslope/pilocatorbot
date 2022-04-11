# pilocatorbot

Telegram bot to parse the RSS feed at <https://rpilocator.com/feed.rss> and send a message when there is an update

## Setup
1. Create a bot on telegram by messaging the botfather and set the environment variable `PILOCBOT_TOKEN` to be the provided api key
2. Send the `\start` message to your new bot from your personal account
3. Set the environment variable `PILOCBOAT_CHAT_ID` to the output of 
```
curl "https://api.telegram.org/bot$PILOCBOT_TOKEN/getUpdates" | jq .result[0].message.from.id
```
