#!/bin/zsh

#1 is data file for addresses
#2 is log file 
#3 is how many times to run the script 
#4 is scene 5 is otuput image file
#6+ all the remaining attributes for traytor

#generate_random n k
#generates n random values from 1 to k
function generate_random {
		shuf -n "$1" "$2" 
}

declare -a addresses
declare -a indices
addresses=( `cat "$1"` )
len=${#addresses[@]}
for (( i = 1 ; i <= len ; i++))
do
		colour="\033[33;3$(( i % 8 ))m"
		echo -e "${colour}===> Testing with" ${i}
		echo -n -e "\e[0m"
		for (( j = 1 ; j <= $3 ; j++ ))
		do
			arguments=()
			used=""
			generate_random $i $1 | while read index
			do
				arguments+=(-w "$index")
				used="$used$index,"
			done

			before=`date +%s%N`
			echo traytor client "${@:6}" "${arguments[@]}" "${4}" "${5}"
			traytor client "${@:6}" "${arguments[@]}" "${4}" "${5}"
			after=`date +%s%N`
			sleep 1
			echo "$i $before $after ${used%,}" >> $2
		done
		echo -e "$colour****** done ******"
done