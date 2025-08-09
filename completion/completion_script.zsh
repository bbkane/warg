#compdef WARG_COMPLETION_APPNAME

# date >> ~/_WARG_COMPLETION_APPNAME_completion.log

local -a comp_values=()
local -a comp_descriptions=()
local -a comp_type
local line

local -a output
output=("${(@f)$(${words[1]} --completion-zsh "${(@)words[2,$CURRENT]}")}")

# Check if we got any output
[[ ${#output} -eq 0 ]] && return 1

# First line is the type
comp_type="${output[1]}"

# Log type
# echo "TYPE: $comp_type" >> ~/_WARG_COMPLETION_APPNAME_completion.log

# Process based on type
case "$comp_type" in
    COMPLETION_TYPE_DIRECTORIES)
        _files -/
        ;;

    COMPLETION_TYPE_DIRECTORIES_FILES)
        _files
        ;;

    COMPLETION_TYPE_NONE)
        return 0
        ;;

    COMPLETION_TYPE_VALUES)
        local i=2
        while (( i <= ${#output} )); do
            comp_values+=("${output[i]}")
            (( i++ ))
        done
        compadd -a comp_values
        ;;

    COMPLETION_TYPE_VALUES_DESCRIPTIONS)
        local i=2
        while (( i <= ${#output} )); do
            comp_values+=("${output[i]}")
            (( i++ ))
            if (( i <= ${#output} )); then
                comp_descriptions+=("${output[i]}")
                (( i++ ))
            fi
        done
        compadd -d comp_descriptions -a comp_values
        ;;

    *)
        # echo "Unknown completion type: $comp_type" >> ~/_WARG_COMPLETION_APPNAME_completion.log
        return 1
        ;;
esac
