all: glossary.exe

glossary.exe:
	go build . -o glossary.exe

clean: 
	rm -f glossary.exe
migrate-dev:
  npx prisma migrate dev --schema ./data/schema.prisma