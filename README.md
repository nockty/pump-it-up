# pump-it-up

**Do not expect to earn money with this bot**

pump-it-up is a simple bot to trade cryptocurrencies during pumps. Big Pump Signal is a Discord server and Telegram channel that regularly schedules cryptocurrency 'pumps'. The concept: everyone buys the coin at the same time, which makes the price go up. It sends a notification through Telegram and Discord at the last moment with the symbol of the coin, so that people can buy it at that time.

This POC polls the Big Pump Signal Telegram channel, waiting for the notification, then buys the revealed coin, waits 1 second and sells the coin. Ideally, the bot is supposed to buy low and sell high.

## build

Dependencies:
- [github.com/Arman92/go-tdlib](https://github.com/Arman92/go-tdlib)
- [github.com/adshao/go-binance/v2](https://github.com/adshao/go-binance/)

`go build -o pump-it-up cmd/pump-it-up/main.go`

## run

`./pump-it-up "BTCCount"` where `BTCCount` is the number of BTC you want to use to buy the coin.

_e.g._ `./pump-it-up "0.00187682"`

## limitations

It doesn't earn money.

- Nothing guarantees that the bot receives the Telegram notification at the same time as other clients. So, lots of people may be faster to buy the coin before the bot. The bot should check the price of the coin before buying (so that it doesn't buy the coin at a high price). Or, at least, somehow, not buy if the notification was received late.
- The bot should also poll Discord.
- The strategy of 'waiting 1s' can be improved.
- Does someone really think the Big Pump Signal race is 100% legit? Chances are some privileged members have the name of the coin before the notification is even sent to anyone.
