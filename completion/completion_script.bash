_WARG_COMPLETION_APPNAME_bash_autocomplete() {
	local cur="${COMP_WORDS[COMP_CWORD]}"
	local output
	output=$(${COMP_WORDS[0]} --completion-bash "${COMP_WORDS[@]:1:$COMP_CWORD}" 2>/dev/null)

	COMPREPLY=()
	[[ -z "$output" ]] && return 0

	local comp_type
	comp_type=$(printf '%s' "$output" | head -n 1)

	case "$comp_type" in
	COMPLETION_TYPE_DIRECTORIES)
		compopt -o dirnames 2>/dev/null
		;;
	COMPLETION_TYPE_DIRECTORIES_FILES)
		compopt -o default 2>/dev/null
		;;
	COMPLETION_TYPE_NONE) ;;
	COMPLETION_TYPE_VALUES | COMPLETION_TYPE_VALUES_DESCRIPTIONS)
		local values
		values=$(printf '%s' "$output" | tail -n +2)
		COMPREPLY=($(compgen -W "$values" -- "$cur"))
		;;
	esac
	return 0
}
complete -F _WARG_COMPLETION_APPNAME_bash_autocomplete WARG_COMPLETION_APPNAME
