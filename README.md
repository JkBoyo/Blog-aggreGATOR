# Blog AggreGATOR
Thanks for taking a look at GATOR! This is my implementation of the [Boot.dev](https://www.boot.dev) [blog aggregator](https://www.boot.dev/courses/build-blog-aggregator-golang).

## Installation Requirements
To install GATOR you have to meet two requirements.
1. Have the [Go](https://go.dev/) toolchain installed on your computer. 
    (This can be checked with `go version`.  If you don't have go installed already I recommend installing this with [webi](webinstall.dev/golang/). )

2. Have an instance of a Postgres Database that you can hook into your program.

## Installation
Once you have verified that the required tooling is insstalled it's time to actually install the thing.

1. Download Source Code
    - I personnaly prefer to do `git clone "https://github.com/jkboyo/Blog-aggregator" gator` inside of the directory you would like to place it.
    - Otherwise you can download it manually from this repo but the above will let you easily download and then update with future changes that might come later.

2. Setup DB
    - DB migrations are handled using [goose](https://github.com/pressly/goose). The installation instructions are really simple on their git page. 
    - Once you have your postgres instance up and running for your distro, run the migrations with `goose postgres "db-connection-string-here/gator" Up`.  An env var will keep you from having to type in the connection string every time you want to run migrations.
    - Due to the nature of postgres db's you should be able to setup this db anywhere you want either locally or on the cloud as long as you correctly configure the db string.
    - It should have the format of "postgres://user:password@hostname:port/gator"

3. Install Program
    - The command `go install` in the source dir should compile the code and place it in your go bin making it available anywhere in your system. The program will be named whatever your root dir of the project is named.
        (I chose GATOR because I love holding shift :upside_down:)
    - Now you have the GATOR program installed! From anywhere you should be able to use the available commands.
    - The very first thing you should do is properly set the db string using `GATOR set-db "db-string"`.
    - From there create your user using the command `GATOR register "username"` which will add your username to the db and will set it in the config file.

## Use the Program
- From here you will have access to the below commands
    ```
    usage GATOR [command] <args>
        login <username>            Logs in user using username
        users                       lists registered users
        agg <formattedtime>         taggregates posts from feeds every time amount passed in. ex. 1s, 1m, 1h, 1d
        feeds                       lists feeds and user that added it
        addfeed <feedname> <url>    Adds feed at specified url wiht specified name
        browse <limit>              list feeds by most recently published. Defaults to limit of 2
        register <username>         Registers a new user
        reset                       resets the database to an empty slate. Only for testing purposes
        set-db <db-conn-str>        update db connection string in config
        follow <url>                Follows feed at the given url
        following                   Lists feeds the logged in user is following
        unfollow <feed-url>         unfollows feed at provided url for logged in user
    ```
- You can list these commands using the -h or --help flags
- To list the description for an individual command you can append the same flags after the command name


