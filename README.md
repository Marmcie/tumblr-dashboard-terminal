# Tumblr dashboard on terminal
Your Tumblr dashboard on terminal.
![image](https://github.com/RayMatsuo/tumblr-dashboard-terminal/blob/main/preview.png?raw=true)
## Set up
1. Register a new application at [Registration page](https://www.tumblr.com/oauth/apps).
2. Copy the consumer key and secret key.
3. Create "tumblr-dt.json" in "~/.config" folder.
4. Fill the necessary fields.
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
- [ ] Fix post count on initial render.
## Continuing
- [ ] Render various NPF posts.
- [ ] Update README.md
### Low priority
- [ ] Ability to switch feed? (e.g. user dashboard, tag search, likes etc.)
- [ ] Possibly implement image viewer feature.
    - Decide on whether to use ASCII image or kitty protocol.
- [ ] Possibly implement other Tumblr features as well.
    - Make sure it doesn't go against the Tumblr API policies.
