#!/bin/sh -e

init_web_hook_local() {
  API_KEY=$2

  stop_web_hook_local
  ngrok http 8080 > /dev/null &
  touch ngrok_output.txt
  sleep 1
  curl -s http://localhost:4040/api/tunnels/command_line > ngrok_output.txt
  NGROK_URL=$(grep -o 'https.*\.io' ngrok_output.txt)
  curl -F "url=${NGROK_URL}" https://api.telegram.org/bot"${API_KEY}"/setWebhook > ngrok_output.txt
  cat ngrok_output.txt
  rm ngrok_output.txt
}

stop_web_hook_local() {
    killall ngrok
}

# Добавьте сюда список команд
using(){
  echo "Commands list:"
  echo "  init_web_hook_local - inits web hook for local development"
  echo "  stop_web_hook_local - stops web hook for local development"
}

command="$1"
if [ -z "$command" ]
then
 using
 exit 0;
else
 $command $@
fi
