# Tumblr dashboard on terminal
Your Tumblr dashboard on terminal.
![image](https://github.com/RayMatsuo/tumblr-dashboard-terminal/blob/main/preview.png?raw=true)
## Set up
1. Register a new application at [Registration page](https://www.tumblr.com/oauth/apps).
2. Copy the consumer key and secret key.
3. Create "tumblr-dt.json" in "~/.config" folder.
4. Fill out the necessary fields.
```json
{
  "consumer_key": "consumer key",
  "secret_key": "secret key"
}
```
5. Launch the program.
6. Authenticate the client using the URL provided in the terminal.
## Build form source
### Prerequisites
- Install Go
### Steps
Run the following command to build the program.
```sh
git clone git@github.com:RayMatsuo/tumblr-dashboard-terminal.git

cd tumblr-dashboard-terminal

go build .
```
## TODO
### High priority
- [x] Figure out how to display asks.
- [x] Switch to bubletea.
- [ ] Fix UI breaking when rendering certain emojis.
### Continuing
- [ ] Render various NPF posts.
- [ ] Update README.md
### Low priority
- [ ] Ability to switch feed? (e.g. user dashboard, tag search, likes etc.)
- [ ] Possibly implement image viewer feature.
    - Decide on whether to use ASCII image or kitty protocol.
- [ ] Possibly implement other Tumblr features as well.
    - Make sure it doesn't go against the Tumblr API policies.
## Known issues
- Certain full width emojis breaks the UI, especially joined emojis.
- When loading new posts, sometimes newly loaded posts contains duplicate posts from the list of already loaded posts.
This happens because dashboard API only allows for offsetting from the latest post on dashboard. 
So if any new post appeared on dashboard since the last time you loaded posts, newly loaded posts contains already loaded posts.
