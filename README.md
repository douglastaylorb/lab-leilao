## Como Rodar o Projeto

1. **Clone o repositório**:
```bash
git clone <URL_REPOSITÓRIO>
cd lab-leilao
```

2. **Execute o projeto com Docker Compose**:
```bash
docker-compose up --build
```

3. **Teste fechamento automático**:

O arquivo para testes está localizado em internal -> infra -> database -> auction -> create_auction.go

### Endpoints

- `GET /auction` - Listar todos os leilões
- `GET /auction/:auctionId` - Buscar leilão por ID
- `POST /auction` - Criar novo leilão
- `GET /auction/winner/:auctionId` - Buscar lance vencedor do leilão
- `POST /bid` - Criar novo lance
- `GET /bid/:auctionId` - Listar lances de um leilão
- `GET /user/:userId` - Buscar usuário por ID
