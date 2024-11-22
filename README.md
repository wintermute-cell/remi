# Remi

> Remi is a straightforward command-line reminder for programmers.

<br>
<img alt="remi mascot" src="logo.webp" width="240" height="240">

It is designed to be simple and quick to use, to not get in your way.
The recommended way to use Remi is to add `remi` to your shell rc, so the
relevant reminders are printed every time you open a terminal emulator.

## How to Use Remi

Calling Remi without any arguments prints all relevant reminders.

Beyond that, this is how you use Remi:

- **Add a Reminder:** You can add a reminder by specifying a timestamp and a
  message. Timestamps can optionally contain a specific timestamp. If you want
  the reminder to appear a certain time before the actual deadline, provide the
  duration parameter.

    ```bash
    remi add <timestamp> <message>
    remi add <timestamp> <duration> <message>
    ```

    - Example:
      ```bash
      remi add 25.12.21@16:30 "Christmas party"
      remi add 25.12.21 "3d" "Reminder before Christmas"
      ```

    - Aliases: `"add", "a", "+"`

- **List Reminders:** You can view all existing reminders.

    ```bash
    remi list
    ```

    - Aliases: `"list", "l", "ls"`

- **Remove a Reminder:** You can delete a specific reminder by providing the
  number next to it **from the latest Remi output** as an ID. (These may change
  from run to run, as remi always tries to index incrementally from 0 to X)

    ```bash
    remi remove <reminderId>
    ```

    - Aliases: `"remove", "r", "rm", "delete", "d", "-"`

## Installation

```bash
git clone https://github.com/wintermute-cell/remi
go build
./remi
```

You can then place the remi binary somewhere in your `$PATH` or make a shell
alias pointing to it.

### Looking for help packaging!

If you'd like to help by packaging for your platform, I'd gladly accept :)

## Usage in Shell Configuration

To have reminders appear when opening a new shell session, you can call the
Remi program from your shell configuration file (e.g., `~/.bashrc`, `~/.zshrc`).

Add the following line to your shell configuration file:

```bash
go run /path/to/remi/main.go
```

By executing this line in the shell configuration, Remi will run automatically
when a new shell session is started. It will show relevant reminders based on
the current time.
