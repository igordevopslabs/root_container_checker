# Use a imagem oficial do Go como base
FROM golang:1.20 AS builder

# Defina o diretório de trabalho
WORKDIR /app

# Copie os arquivos necessários para o diretório de trabalho
COPY . .

# Faça o build do programa Go
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

# Use a imagem base scratch para um contêiner minimalista
FROM scratch

# Crie um usuário não-root
# No caso do scratch, você precisa especificar o UID e GID manualmente
USER 1000:1000

# Copie o binário compilado do builder para a imagem final
COPY --from=builder --chown=1000:1000 /app/main /main

# Execute o binário compilado como o usuário não-root
CMD ["/main"]
