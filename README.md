rain-check
===

## Overview
I consistently leave the house/office without checking the weather or looking out the window first, and I often waste time having to go back and get my umbrella.  Rain-check is a cli utility meant to run as a background process that sends push notifications to your phone when it starts/stops raining/snowing/hailing, so you never go outside unprepared.

## Usage
Rain-check requires a [Pushover](https://pushover.net/) account with an app key and user key, along with an [OpenWeatherMap](http://openweathermap.org/api) API token.

To build:
`go install`

Rain check runs and a command-line utility.  Type `rain-check help` for option information.
