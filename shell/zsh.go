package shell

var ZSH = `
[[ ! $(type -p kit) ]] && return
REAL_KIT=$(which kit)
CODE_ROOT_FOLDER="$HOME/src/github.com"
kit() {
    case $1 in
        cd)
            if [[ ! -z $2 && -d "$CODE_ROOT_FOLDER/$2" ]]; then
                cd "$CODE_ROOT_FOLDER/$2"
            else
                dir=$(find $CODE_ROOT_FOLDER -maxdepth 2 -mindepth 2 | cut -c$((${#CODE_ROOT_FOLDER}+2))- | sort | fzf)
                [[ ! -z $dir ]] && cd "$CODE_ROOT_FOLDER/$dir"
            fi
        ;;
        edit)
            if [[ ! -z $2 && -d "$CODE_ROOT_FOLDER/$2" ]]; then
                eval "$EDITOR $CODE_ROOT_FOLDER/$2"
            else
                dir=$(find $CODE_ROOT_FOLDER -maxdepth 2 -mindepth 2 | cut -c$((${#CODE_ROOT_FOLDER}+2))- | sort | fzf)
                [[ ! -z $dir ]] && eval "$EDITOR $CODE_ROOT_FOLDER/$dir"
            fi
        ;;
        clone)
            if [[ ! -z $2 ]]; then
                if [[ ! -d "$CODE_ROOT_FOLDER/$2" ]]; then
                    git clone "https://github.com/$2.git" "$CODE_ROOT_FOLDER/$2"
                fi
                if [[ -d "$CODE_ROOT_FOLDER/$2" ]]; then
                    cd "$CODE_ROOT_FOLDER/$2"
                fi
            else
                eval "$REAL_KIT clone --help"
            fi
        ;;
        *)
            eval "$REAL_KIT $@"
        ;;
    esac
}
`
