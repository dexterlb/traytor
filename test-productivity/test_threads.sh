#!/bin/zsh

#1 is the range (for example 3-5 for numbers of threads to test for)
#2 is log file 
#3 is how many times to run the script 
#4 is scene 5 is otuput image file
#6+ all the remaining arguments for traytor

if [[ -z ${TRAYTOR_BIN} ]]; then
    TRAYTOR_BIN=traytor
fi

if [[ "${1}" =~ '([0-9]+)-([0-9]+)' ]]; then
    first="${match[1]}"
    last="${match[2]}"
    echo "will test from $first to $last threads"
else
    echo "please enter correct number of threads"
fi

echo "num_threads start_time end_time" >> ${2}
for (( i = first ; i <= last ; i++))
do
		colour="\033[33;32m"
		echo -e "\n${colour}===> Testing for ${i}"
		echo -n -e "\e[0m"
		for (( j = 1 ; j <= ${3} ; j++ ))
		do
			arguments=(-j ${i})

			before=`date +%s%N`
			"${TRAYTOR_BIN}" render "${@:6}" "${arguments[@]}" "${4}" "${5}"
			after=`date +%s%N`
			sleep 1
			echo "$i $before $after" >> "${2}"
		done
		echo -e "$colour****** done ******\e[0m"
done
