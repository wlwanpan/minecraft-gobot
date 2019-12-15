# minecraft-gobot

minecraft-gobot is a set of tools (launcher + bot) to run a discord-bot that help manage your minecraft server on AWS.

- launcher: grpc server that launches the minecraft server (server.jar) and interacts with its stdin/stdout for executing your regular minecraft server console commands.

- bot: discord bot, that is meant to run 24/7 and performs rpc calls to the launcher to control the minecraft server (server.jar).

## Basic commands
```bash
# Start the discord-bot
go run main.go bot

# Start the minecraft-launcher
go run main.go launcher
```

## Setting up your AWS instances
TODO: complete doc

## Setting up your credentials

- discord app
TODO

- AWS (iam)
TODO
