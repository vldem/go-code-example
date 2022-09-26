<h3>Supervisor's events notifier</h3>

## Features

This service listens events from [supervisor](https://supervisor.readthedocs.io/en/latest/) and sends corresponding notifications to email and/or telegram.

## Configuration

You should set needed parameters in the Yaml config file.
spv_notif_cfg.yml example:

    email:
      from: YOUR_FROM_EMAIL_HERE
      to: YOUR_TO_EMAIL_HERE
      subject: "[supervisor notification]"
      mail_prog:
        cmd: "/usr/sbin/sendmail"
        args: - "-t" - "-oi"

    telegram:
      botkey: YOUR_TELEGRAM_BOT_APIKEY_HERE
      chatid: YOUR_TELEGRAM_CHAT_ID_WITH_BOT_HERE

You should place the config file to directory `../config` or `../../config` from binary file's directory.

# Supervisor configuration

There are two points in supervisor configuration that should be done in order to receive notification about needed supervisor's events.

You need to configure eventlistener in supervisord.conf file. Below is an example of such configuration:

    [eventlistener:spv_notif]
    command=/YOUR_PATH/spv_notif
    events=PROCESS_STATE_EXITED,PROCESS_LOG_STRERR
    process_name=%(program_name)s_%(process_num)s
    numprocs=1
    autorestart=true

Also it's necessary to add events line to program configuration like that:

    [program:worker]
    command=/YOUR_PATH/myworker
    numprocs=1
    startsecs=0
    autostart=true
    autorestart=true
    process_name=%(program_name)s*%(process_num)02d
    events=PROCESS_STATE_EXITED,PROCESS_LOG_STRERR

If your worker exits due to some reason the event listener will catch the event PROCESS_STATE_EXITED and
will send a notification to email or the telegram chat that you specified in configuration.
