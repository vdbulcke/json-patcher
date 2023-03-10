## json-patcher completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(json-patcher completion bash)

To load completions for every new session, execute once:

#### Linux:

	json-patcher completion bash > /etc/bash_completion.d/json-patcher

#### macOS:

	json-patcher completion bash > $(brew --prefix)/etc/bash_completion.d/json-patcher

You will need to start a new shell for this setup to take effect.


```
json-patcher completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -d, --debug   debug mode enabled
```

### SEE ALSO

* [json-patcher completion](json-patcher_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 11-Mar-2023
