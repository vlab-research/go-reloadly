docker build \
	-t go-reloadly \
	. &> /dev/null

docker run \
	-v $(pwd):/app \
	-t go-reloadly \
	sh -c "goimports -l -w ."
