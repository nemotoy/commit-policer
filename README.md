# commit-policer

自分のGithubのコミットを監視する的なやつをつくる

## reference

* [Github Developers](https://developer.github.com/v3/)
    * [PushEvent](https://developer.github.com/v3/activity/events/types/#pushevent)
    * [rate limit](https://developer.github.com/v3/#rate-limiting)

* [go-github](https://github.com/google/go-github)


## API

* 指定リポジトリのコミット全件取得  
$ curl -is -H "Authorization: token {token}" https://api.github.com/repos/:user/:repo/commits

* 指定ユーザーのイベント取得    
$ curl -is https://api.github.com/users/:user/events
