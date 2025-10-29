# Where should we write the GIF?
Output demo.gif

Set Shell "bash"
Set FontSize 14
Set LineHeight 1.2
Set Width 1200
Set Height 800
Set Framerate 24
Set Padding 5
Set FontFamily "FiraCode Nerd Font"
Set WindowBar Colorful

# Type a command in the terminal.
Env GRPC_SERVER_ADDRESS "localhost:8080"
Env LOG_OUTPUT "bin/client.log"
Env SPEW_DUMP_OUTPUT "bin/spew.log"
Env LOG_TRUNCATE "true"
Type "./bin/gophkeeper-client"

Enter
# Pause for dramatic effect...
Sleep 1s

Type "demo"
Sleep 300ms
Tab
Type "demo"

Sleep 500ms
# Run the command by pressing enter.
Enter

Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Down

# Admire the output for a bit.
Sleep 2s