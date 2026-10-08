# fossil

A tool which parses exported Spotify user data, stores it in an SQLite database and can query/run statistics on the data.

To download the required data, carry out the following steps:
1. Go to the [Spotify](https://spotify.com) website.
2. Click on *Account privacy* under the *Security and privacy* section.
3. Select the types of data you wish to download. Currently, `fossil` only has support for the *Extended streaming history* section.
4. Click *Request data*.
5. You will receive an email within 30 days containing a link allowing you to download the data.

The relevant files that can be parsed by `fossil` have the filename format `Streaming_History_Audio_{year}.json` where `{year}` is the
year that the data pertains to. If there is a large amount of data, some years will overflow into multiple files,
which will be titled such as `Streaming_History_Audio_{year}_1.json`.

# Build

Ensure you have `go` `1.26.6` or later installed, and run `go build . -o fossil` in the project root. This will generate a `fossil` executable
file in the project root..

# Usage

To import an *audio streaming history* file, from the project root, run `./fossil import {filepath}` where `{filepath}` is the path of the file.
`fossil` will then parse the file and store its data in the database, skipping data that already exists, so you can re-import a file that has already
been partially imported (maybe you imported the `2026` file earlier in the year, and then re-requested the data later in the year and want to import the
latest data) without duplicating data.

Several statistics-related commands exist, which will be cleaned up soon into a more usable CLI but were just implemented quickly.
Currently, there are (replace `{executable}` with the path to the executable file, for example `./fossil`):
1. `{executable} year-leaderboard {year}` where `{year}` is the year to query data for, which will generate a leaderboard of tracks listened in that year,
including the amount of time spent listening to each track.
2. `{executable} total-listening-time` sums the time spent listening to every track you've ever listened to (only including data that has been imported, of course),
and gives you a total amount of time that you have spent listening to Spotify.
1. `{executable} lifetime-leaderboard` where `{year}` which will generate a leaderboard of tracks listened in your Spotify lifetime, as well as a leaderboard
of artists across your Spotify lifetime.

# TODO

- [x] Consider renaming `Song` to something more generic because it can cover podcasts etc too
- [ ] Allow importing multiple files at a time
- [ ] Allow parsing of 'video' history files too
- [ ] Investigate feasibility of YouTube Music support
