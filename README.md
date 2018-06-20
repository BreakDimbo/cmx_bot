# Feature
This is a bot colletion for mastodon. It will be extended as needed.
Now there are 2 kinds of bot:
* intelbot
* firebot

## intelbot
Intelbot is a bot which collects and analyzes the toots.
It will show the following items daily and weekly
* keywords of top five
* total toots in local timeline
* total count of users who tooted in local timeline
* the top 3 users who like toot most in local timeline
* the top 1 user who like toot most in home

## firebot
Firebot is a bot for reblogging the toot with review in local timeline.
Reply the toot with review and @firebot.

# Install
1. [install golang](https://golang.org/doc/install)
1. git clone git@github.com:BreakDimbo/cmx_bot.git
2. change to working directory
3. make get
4. make deps
5. make bot
6. make fbot
7. [install elastic search 6.3.0](https://www.elastic.co/downloads/elasticsearch) (maybe you need install java runtime env)
8. sudo touch /var/log/mastodon_bot
9. sudo chown ${your_current_user} /var/log/mastodon_bot
10. set the config file in config/development.toml according to config/development.demo.toml

# Usage
1. run elasticsearch
2. run intelbot
```bash
#{working_directory}/bin/bot
```
3. run firebot
```bash
#{working_directory}/bin/fbot
```

there are xxx.service files for systemd usage in dir config

# TODO

- [x] set config.yml
- [x] follow & unfflow automatically
- [x] post weekly and daily
- [x] collect local toots of followers
- [x] refactor the structure
- [x] log
- [x] remove config.toml file
- [x] readme
- [ ] show the favourited plus reblogged most toot daily
- [ ] generate graph about toots count everyday
- [ ] cache info in nsq when es connect down
- [ ] fix the way of reading dict file