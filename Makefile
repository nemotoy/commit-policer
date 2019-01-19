APPNAME=commit-policer
LOCALPORT=8080

ping: 
	curl -is https://commit-policer.herokuapp.com/ping

## heroku cmd memo
up: 
	heroku ps:scale web=1 --app $(APPNAME)

down: 
	heroku ps:scale web=0 --app $(APPNAME)

hinfo:
	heroku apps:info

hlog:
	heroku logs

hpush:
	git push heroku master
