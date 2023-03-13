# Sonnenbatterie trend notification

![build workflow](https://github.com/dabump/sonnenbatterie/actions/workflows/build.yaml/badge.svg)
## Background
If you are the owner of a Sonnen battery in your house, there is a good chance that the battery
will be connected to your home network, with a local IP, that can be used to obtain status info.
The battery has a unauthenticated REST endpoint that can be used to determine the
Usable state of charge (USOC).

## Problem statement
This is a personal project, with a bespoke use case. 
I own a Tesla electric vehicle, and want to capatalise on charging 
my car with mostly free solar energy. This utililty will monitor the battery state of charge 
on my behalf, based on simplistic upward or downward trends, and notify me when a high or low
threshold was crossed that will trigger a manual, human action, to start or stop the charge
process of my Tesla to avoid running the risk using grid power (ie, paying for electricity).

## How to get running
1. Clone this repo to a local directory
2. Copy the template config `cp config.cfg.template config.cfg`
3. Change the config file, and add your shoutrrr url string (https://containrrr.dev/shoutrrr/)
4. Run the daemon by either using a locally built docker container, or as GO service
4.1 Docker: `make docker`
4.2 GO: `go run cmd/daemon/main.go`

## HTTP API endpoint
Two HTTP API endpoints are available for default configured port `8881` (port value can be changed
via the confguration file). Default rate limit set to 1 per second request to sonnen batterie
 - `http://host:8881/` <- Returns the json payload directly from the sonnen batterie
 - `http://host:8881/status` <- returns `200 OK` when service up and accepting traffic
## Shoutrrr?
Instead of building custom / bespoke notification logic, using Shoutrrr allows more range of
notification services to be used (example, slack, telegram, email, ms teams, etc)

## Future changes / known limitations
- Investigate Tesla API's to authomatically start / stop charge
- Abstract Sonnenbatterie specifics to allow extension / addition of alternative battery systems
- Improve trend analysis
- Add support for multiple low and upper threshold notifications
- Add notification black out periods (ie, swallow notifications between 10pm-6am)
- Look at extracting the daemon so other GO developers can use a package approach
