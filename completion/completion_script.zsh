#compdef WARG_COMPLETION_APPNAME

date >> ~/_WARG_COMPLETION_APPNAME_completion.log

local -a comp_values=()
local -a comp_descriptions=()
local -a comp_type

local state="expecting_type"
local line
# TODO: switch on comp_type to switch on the rest of the lines
while IFS= read -r line && [[ -n "$line" ]]; do
    if [[ "$state" == "expecting_type" ]]; then
        comp_type="$line"
        state="expecting_value"
    elif [[ "$state" == "expecting_value" ]]; then
        comp_values+=("$line")
        state="expecting_description"
    else
        comp_descriptions+=("$line")
        state="expecting_value"
    fi
done < <(${words[1]} --completion-zsh "${(@)words[2,$CURRENT]}")

echo "$comp_values" >> ~/_WARG_COMPLETION_APPNAME.log

compadd -d comp_descriptions -a comp_values
