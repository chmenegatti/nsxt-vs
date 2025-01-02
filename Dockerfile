# Etapa 1: Construção
FROM golang:1.23-alpine AS build

# Definir o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copiar os arquivos de código para dentro do contêiner
COPY . .

# Instalar as dependências e compilar a aplicação
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s"

# Etapa 2: Imagem final
FROM alpine:latest

# Definir o diretório de trabalho dentro do contêiner
WORKDIR /root/

# Copiar o binário compilado da etapa de construção
COPY --from=build /app/nsxt-vs .
COPY --from=build /app/config.yaml .

# Expor a porta 4040
EXPOSE 4040

# Comando para rodar a aplicação
CMD ["./nsxt-vs"]
