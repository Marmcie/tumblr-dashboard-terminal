# Tumblr dashboard on terminal
Your Tumblr dashboard on terminal!
## Features
With this project you can:
- Look at posts from your dashboard.
- Read posts from any blog.
- Search for posts with specific tag.
## Screenshot
![image](https://github.com/Marmcie/tumblr-dashboard-terminal/blob/main/doc/preview_01.png?raw=true)
## Set up
1. Register a new application at [Registration page](https://www.tumblr.com/oauth/apps).
2. Copy the consumer key and secret key.
3. Create "tumblr-dt.json" in "~/.config" folder.
> Alternatively you can use -config flag to pass the path to the config file directly on launch.
```
#Example
tumblr-dt -config="Path-to-your-config-json-file"
```
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
- Install [Go](https://go.dev/doc/install)
### Steps
Run the following command to build the program.
```sh
git clone git@github.com:Marmcie/tumblr-dashboard-terminal.git

cd tumblr-dashboard-terminal

go build .
```
