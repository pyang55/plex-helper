# Plex-helper
A helper tool to remove and filter movies based on Rotton Tomatoes ratings. This is built primarly for my own use but I figured I would share it. I was trying to clear some space in my Plex and I realized how many old movies I downloaded that were high quality 10GB files of absolute shit (Jupiter Ascending in 4k). This helps clear out some of the cruft. 

## Setup
- Download binary from release
- chmod +x the binary to make it executable ```chmod +x plex-helper-*``` 
- rename the binary to be plex-helper if it makes it easier on you ```mv plex-helper-* plex-helper```
- Move binary to /usr/local/bin or /usr/bin/ ```mv plex-helper /usr/local/bin/```
- If you are on mac you may have to allow permissions for a non apple store app (Apple menu > System Preferences, click Security & Privacy , click General)
- Know the ip (public or private) of your Plex server
- Retrieve plex auth token https://support.plex.tv/articles/204059436-finding-an-authentication-token-x-plex-token/

## Use Cases
- Allows users to list/delete movies/shows based on rating
- Allows users to choose movies(or shows) to keep (because we all have those 40% RT cult classics that we love...Eurotrip)

## Common Commands
```plex-helper list <movies|shows> --ip-addr <plex-ipaddr> --token <token> --rating 5 # list movies or shows with RT score of 5 or less```

```plex-helper list <movies|shows> --ip-addr <plex-ipaddr> --token <token> --rating 5 --minimal # non pretty output of shows/movies```

```plex-helper delete <movies|shows> --ip-addr <plex-ipaddr> --token <token> --rating 5 # delete shows with RT score of 5 or less```

```plex-helper list <movies|shows> --ip-addr <plex-ipaddr> --token <token> --rating 10 --minimal # just get a list of all your movies or shows```

## Attention
- Be careful when using this!! This can wipe out alot of shows/movies by accident if you are not careful. There are prompts to help guide you but it can't save you from yourself

## Notes
- I have only really used this on Mac or some kind of linux distribution. I have no idea how this will work on windows (even though i will build a binary for it)
