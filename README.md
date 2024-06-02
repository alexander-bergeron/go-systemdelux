# go-systemdelux

A lightweight golang program to start, stop, and monitor processes. This is similar to systemd but can be run at the user level.

## Architecture

- cron/systemd setup for start on reboot

- Daemon process that runs constantly, checking and then sleeping

- Utility function to work on the config files. Migrate to API down the line

## Features

1. List monitored services

2. Start a service

3. Stop a service

4. Restart a service

5. Add/Delete/Pause/Update a monitored service

6. Inspect a service

