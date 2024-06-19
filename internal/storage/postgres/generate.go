package postgres

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../../bin/minimock -i UserStorage -o ./mocks/ -s "_minimock.go"
//go:generate ../../../bin/minimock -i LogStorage -o ./mocks/ -s "_minimock.go"
