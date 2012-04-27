## Installation

	go build
	install gonoip /usr/local/sbin

	gonoip -a > /etc/gonoip.conf
	chmod 0600 /etc/gonoip.conf

Put the following in /etc/cron.hourly/gonoip

	#!/bin/sh
	/usr/local/bin/gonoip

And the following in /etc/cron.weekly/gonoip

	#!/bin/sh
	/usr/local/bin/gonoip -f

