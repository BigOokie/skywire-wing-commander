#!/bin/bash
#Will Automatically get data from the Telegram API to get Chat ID and Telegram username. 

echo After starting a chat with the Botfather press start on your Chat with your bot. 
echo Please enter your APIKey from the BotFather
read APIKEYS

TELEUSERNAME=$(curl -s "https://api.telegram.org/bot$APIKEYS/getUpdates" | grep -Po '"username":\K.*?(?=,)'| head -1)
echo $TELEUSERNAME

ChatID=$(curl -s "https://api.telegram.org/bot$APIKEYS/getUpdates" | grep -Po '"id":\K.*?(?=,)'| head -1 )
echo $ChatID

sed -i -e "s/apikey = .*/apikey = "$APIKEYS"/" config.toml
sed -i -e "s/admin = .*/admin = "$TELEUSERNAME"/" config.toml
sed -i -e "s/chatid = .*/chatid = $ChatID/" config.toml


