# Weather Forecast Telegram Bot

# Overview
This Telegram bot provides users with daily weather forecast notifications at a chosen time in the UTC timezone. Users can subscribe and unsubscribe to receive these notifications based on their preferences.

## Features
**Subscription**: Users can subscribe to receive daily weather forecast notifications.\
**Unsubscription**: Users can unsubscribe at any time to stop receiving weather updates.

## Installation
Clone this repository:
```
git clone https://github.com/c1kzy/Telegram-Subscription-Bot.git
cd task-2.5-subscription-bot

```
Click the link "https://t.me/SubscriptionWeatherBot" and run the bot

Deploy a server
You can deploy your server any way you want, but I find it really quick and easy to use ngrok. Ngrok allows you to expose applications running on your local machine to the public internet.

How to install it?
You can download it from the website directly

```
https://ngrok.com/download
```
or install ngrok via Chocolatey
```
$ choco install ngrok
```

Running a server using ngrok
Once you install ngrok, you can run this command on another terminal on your system:
```
ngrok http <port>
```
Example:
```
ngrok http port 8080
```
Here, https://ed40-178-150-143-55.ngrok.io is the public temporary URL for the server running on port 8080 of my machine.
Now, all we need to do is let telegram know that our bot has to talk to this url whenever it receives any message. We do this through the telegram API. Enter this in your terminal :
```
curl -F "url=https://ed40-178-150-143-55.ngrok.io"  https://api.telegram.org/bot<your_api_token>/setWebhook
```

