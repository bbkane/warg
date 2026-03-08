function __warg_completion_WARG_COMPLETION_APPNAME
    set -l tokens (commandline -xpc)
    set -l current (commandline -ct)

    set -l output (WARG_COMPLETION_APPNAME --completion-fish $tokens[2..] "$current" 2>/dev/null)

    # Check if we got any output
    test (count $output) -eq 0; and return 1

    # First line is the type
    set -l comp_type $output[1]

    switch $comp_type
        case COMPLETION_TYPE_DIRECTORIES
            __fish_complete_directories "$current" ""
            return
        case COMPLETION_TYPE_DIRECTORIES_FILES
            __fish_complete_path "$current"
            return
        case COMPLETION_TYPE_NONE
            return
        case COMPLETION_TYPE_VALUES COMPLETION_TYPE_VALUES_DESCRIPTIONS
            for i in (seq 2 (count $output))
                echo $output[$i]
            end
    end
end

complete -c WARG_COMPLETION_APPNAME -f -a '(__warg_completion_WARG_COMPLETION_APPNAME)'
