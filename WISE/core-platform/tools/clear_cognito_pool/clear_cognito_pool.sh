#!/bin/bash

usage() {
    echo "Usage: $0 [cognito_user_pool_id]" >&2
    exit 1
}

# Make sure we have one param for user pool id
if [ ! $# = 1 ]; then
	usage
fi

export USER_POOL_ID=$1
export AWS_DEFAULT_REGION=us-west-2
export AWS_DEFAULT_PROFILE=wiseus
 
RUN=1
until [ $RUN -eq 0 ] ; do
	echo "Listing users"
	USERS=`aws cognito-idp list-users  --user-pool-id ${USER_POOL_ID} | grep Username | awk -F: '{print $2}' | sed -e 's/\"//g' | sed -e 's/,//g'`

	if [ ! "x$USERS" = "x" ]; then
		for user in $USERS; do
			echo "Deleting user $user"
			aws cognito-idp admin-delete-user --user-pool-id ${USER_POOL_ID} --username ${user}
			echo "Result code: $?"
			echo "Done"
		done
	else
		echo "Done, no more users"
		RUN=0
	fi
done
