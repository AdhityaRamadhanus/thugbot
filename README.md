# thugbot
[![Go Report Card](https://goreportcard.com/badge/github.com/AdhityaRamadhanus/thugbot)](https://goreportcard.com/report/github.com/AdhityaRamadhanus/thugbot)

<p>
  <a href="#installation">Installation |</a>
  <a href="#usage">Usage</a> |
  <a href="#results">Results</a> |
  <a href="#docker">Docker (soon)</a> |
  <a href="#licenses">License</a>
  <br><br>
  <blockquote>
	Slackbot to thuglify your image
  </blockquote>
</p>

Installation
------------
* Install openCV version 2.4.x [opencv](http://docs.opencv.org/2.4/doc/tutorials/introduction/table_of_content_introduction/table_of_content_introduction.html)
* git clone
* go get -v
* go build

Usage
------------
* You need a slack token for this, please read this if you haven't got one 
[slack token](https://api.slack.com/bot-users)
* You need to define your slack token for bot users in an environment variable called SLACK_TOKEN
* I prefer to use .env file for this (example below)
```
SLACK_TOKEN=your-token-here
```

Results
------------
![thugbot-1](https://cloud.githubusercontent.com/assets/5761975/24082152/4f2dc7fa-0cf3-11e7-93bd-2ba626192b2a.png)
![thugbot2](https://cloud.githubusercontent.com/assets/5761975/24082151/4f26ae0c-0cf3-11e7-9424-f871f1c13e13.png)

License
----

MIT © [Adhitya Ramadhanus]# thugbot
Slack bot to thuglify your photos
