#!/bin/zsh

#1 is data file for addresses
#2 is the range (for example 3-5 for numbers of hosts to test for)
#3 is log file 
#4 is how many times to run the script 
#5 is scene 6 is otuput image file
#6+ all the remaining attributes for traytor

#generate_random n k
#generates n random values from 1 to k
function generate_random {
		shuf -n "$1" "$2" 
}

declare -a addresses
declare -a indices
addresses=( `cat "$1"` )

if [[ "${2}" =~ '([0-9]+)-([0-9]+)' ]]; then
    first="${match[1]}"
    last="${match[2]}"
    echo "will test from $first to $last hosts"
else
    echo "please enter correct number of hosts"
fi

echo "num_hosts start_time end_time hosts" >> "${3}"
for (( i = first ; i <= last ; i++))
do
		colour="\033[33;32m"
		echo -e "${colour}===> Testing for ${i}\n"
		echo -n -e "\e[0m"
		for (( j = 1 ; j <= $4 ; j++ ))
		do
			arguments=()
			used=""
			generate_random $i $1 | while read index
			do
				arguments+=(-w "$index")
				used="$used$index,"
			done
			echo -e "${colour}===> with [${used%,}]"
			echo -n -e "\e[0m"


			before=`date +%s%N`
			traytor client "${@:7}" "${arguments[@]}" "${5}" "${6}"
			after=`date +%s%N`
			sleep 1
			echo "$i $before $after ${used%,}" >> "${3}"
		done
		echo -e "$colour****** done ******"
done
