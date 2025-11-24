# pb (personal blog)

*pb* is simple pet-project to demonstrate REST API knowledge.
It uses *postgresql*, *pgx*, *standard library* and *goose* for migrations.

## Run with docker compose

```shell
git clone https://github.com/ochamekan/pb
cd pb
cp .env.example .env 
docker compose up -d
```

## API

API running on `http://localhost:8080`

## Swagger

SwaggerUI running on `http://localhost:8080/swagger/index.html`

## Clean up

Stop docker and remove containers:
`docker compose down`

Delete images:
`docker rmi pb-backend pb-migrate postgres`

Delete volumes:
`docker volume rm pb_articles`

 
 







