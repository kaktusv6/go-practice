package db

//go:generate bash -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i route256/libs/db.QueryEngine -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i route256/libs/db.QueryEngineProvider -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i route256/libs/db.Transactor -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i route256/libs/db.DB -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i route256/libs/db.Client -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i route256/libs/db.Manager -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i github.com/jackc/pgx/v4.Tx -o ./mocks/tx_minimock.go
