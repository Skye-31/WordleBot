# WordleBot

Public instance: [Invite](https://discord.com/oauth2/authorize?client_id=948289152514924584&scope=applications.commands)

## Self-Hosting
#### Setup
```shell
git clone https://github.com/Skye-31/WordleBot
cd WordleBot
cp example.config.json config.json
# Edit your config.json file as appropriate
# You should use a Postgresql 12+ database
# Install Docker & Docker Compose
docker-compose build
```
#### Running
The first time you run the bot, you should use the flags --sync-commands and --sync-db to automatically set up your bot commands & database tables.
You can do this by commenting out the line in your `docker-compose.yml`.

After that, you run the following command.
```shell
docker-compose up -d
```
