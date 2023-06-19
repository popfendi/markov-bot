# Markov chain tg bot

A telegram bot that uses markov chains to learn to speak like your groups users.

## How it works

The bot trains each message into a markov chain, and can be configured to send messages at a time interval. It's like your little mascot / friend who tries to imitate you.

## Usage

1. build the binary
2. you can edit the templateRunScript.sh to run the bot with your args.

### Project Structure

```
├── handlers.go -- most of the bot logic is in here
├── init2.txt -- contains some starting messages to train the bot on, can edit this or remove it all together for a completely dumb start.
├── main.go
├── markov
│   └── markov.go -- markov chain training and message generation
├── README.md
└── templateRunScript.sh

```
