# minecraft-gobot (WIP)

minecraft-gobot is a set of tool (mcs + bot) to run a discord-bot that help manage your minecraft server on AWS.

- mcs: grpc server that wraps the minecraft server (server.jar) and pipes stdin/stdout for executing your regular minecraft server console commands.

- bot: discord bot, that is meant to run 24/7 and performs rpc calls to the mcs wrapper to control the minecraft server (server.jar).

## How to setup

- First create a config.yaml file in the following format:
```yaml
mcs:
  port: PORT_TO_RUN_MCS
  server_jar: "server.jar" # specify the minecraft launch server.

bot:
  mcs_addr: ELASTIC_IP_RUNNING_MCS
  mcs_port: PORT_MCS_IS_RUNNING
  whitelisted_channel_ids:
    - CHANNEL_ID_1
    - CHANNEL_ID_2

# Required only for S3 backup uploads.
aws:
  region: AWS_REGION
  s3_bucket_name: AWS_S3_BUCKET_NAME
```

## Basic commands
```bash
# Start the discord-bot
go run gobot.go bot -t [DISCORD_BOT_TOKEN]

# Start the minecraft server
go run gobot.go mcs
```

## Setting up your AWS instances
TODO: complete doc

## Setting up your credentials

- discord app
TODO

- AWS (iam)
TODO
